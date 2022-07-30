package versions

import (
	"fmt"
	"github.com/fabelx/go-solc-select/internal/errors"
	"github.com/fabelx/go-solc-select/internal/utils"
	"github.com/fabelx/go-solc-select/pkg/config"
	"github.com/stretchr/testify/assert"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"testing"
)

var testCurrentVersion = "0.4.7"
var testVersions = map[string]string{
	testCurrentVersion: testCurrentVersion,
	"0.5.1":            "0.5.1",
	"0.6.7":            "0.6.7",
	"0.8.1":            "0.8.1",
}
var notInstalledVersion = "0.4.1"

func TestMain(m *testing.M) {
	err := setup()
	if err != nil {
		log.Fatalf("Failed to run tests during setup. Error: %v", err)
	}

	code := m.Run()
	shutdown()
	os.Exit(code)
}

// setup Setups test environment
func setup() error {
	// creates dirs for testing
	config.SolcDir = filepath.Join(config.HomeDir, ".test-gsolc-select")
	config.SolcArtifacts = filepath.Join(config.SolcDir, "artifacts")
	err := os.MkdirAll(config.SolcArtifacts, 0755)
	if err != nil {
		return err
	}

	// initializes current global version as EMPTY value
	config.CurrentVersionFilePath = filepath.Join(config.SolcDir, "global-version")
	err = os.WriteFile(config.CurrentVersionFilePath, []byte(""), 0755)
	if err != nil {
		return err
	}

	// adds fake solc compilers
	for _, v := range testVersions {
		name := fmt.Sprintf("solc-%s", v)
		dirPath := filepath.Join(config.SolcArtifacts, name)
		err = os.Mkdir(dirPath, 0755)
		if err != nil {
			return err
		}

		fakeFilePath := filepath.Join(dirPath, name)
		err = os.WriteFile(fakeFilePath, []byte(""), 0755)
		if err != nil {
			return err
		}
	}

	return nil
}

// shutdown Removes all test files and dirs
func shutdown() {
	os.RemoveAll(config.SolcDir) //nolint:errcheck
}

func TestGetInstalled(t *testing.T) {
	result := GetInstalled()
	assert.Equal(t, testVersions, result)
}

func TestGetCurrent(t *testing.T) {
	testCases := []struct {
		name     string
		expected string
		err      error
		before   func()
		after    func()
	}{
		{
			name:     "test EMPTY version",
			expected: "",
			err:      &errors.NoCompilerSelected{},
			before:   func() {},
			after:    func() {},
		},
		{
			name:     "test failed to read file with version",
			expected: "",
			err:      &fs.PathError{},
			before: func() {
				config.CurrentVersionFilePath = "pathDoesNotExist"
			},
			after: func() {
				config.CurrentVersionFilePath = filepath.Join(config.SolcDir, "global-version")
			},
		},
		{
			name:     "test not installed version, but specified as current version",
			expected: "",
			err: &errors.NotInstalledError{
				Version: notInstalledVersion,
			},
			before: func() { //nolint:errcheck
				os.WriteFile(config.CurrentVersionFilePath, []byte(notInstalledVersion), 0755) //nolint:errcheck
			},
			after: func() {
				os.WriteFile(config.CurrentVersionFilePath, []byte(""), 0755) //nolint:errcheck
			},
		},
		{
			name:     "test installed version",
			expected: testCurrentVersion,
			err:      nil,
			before: func() {
				os.WriteFile(config.CurrentVersionFilePath, []byte(testCurrentVersion), 0755) //nolint:errcheck
			},
			after: func() {},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.before()

			result, err := GetCurrent()
			assert.IsType(t, testCase.err, err)
			assert.Equal(t, testCase.expected, result)

			testCase.after()
		})
	}
}

func TestGetAvailable(t *testing.T) {
	expectedType := map[string]string{}
	result, err := GetAvailable()

	assert.NoError(t, err)
	assert.IsType(t, expectedType, result)
	for key := range testVersions {
		assert.Contains(t, result, key)
	}
}

func TestLinuxGetAvailableVersions(t *testing.T) {
	expected := []string{"0.4.10", "0.4.0", "0.8.13"}
	platform := &LinuxPlatform{Name: config.LinuxAmd64}
	result, _ := platform.GetAvailableVersions()
	for _, key := range expected {
		assert.Contains(t, result, key)
	}
}

func TestMacGetAvailableVersions(t *testing.T) {
	expected := []string{"0.3.6", "0.5.17", "0.8.13"}
	platform := &MacPlatform{Name: config.MacosxAmd64}
	result, _ := platform.GetAvailableVersions()
	for _, key := range expected {
		assert.Contains(t, result, key)
	}
}

func TestWindowsGetAvailableVersions(t *testing.T) {
	expected := []string{"0.4.1", "0.5.17", "0.8.13"}
	platform := &WindowsPlatform{Name: config.WindowsAmd64}
	result, _ := platform.GetAvailableVersions()
	for _, key := range expected {
		assert.Contains(t, result, key)
	}
}

func TestGetBuilds(t *testing.T) {
	testCases := []struct {
		name     string
		expected []*utils.BuildData
		platform Platform
	}{
		{
			name:     "test builds for linux",
			expected: []*utils.BuildData{},
			platform: &LinuxPlatform{Name: config.LinuxAmd64},
		},
		{
			name:     "test builds for windows",
			expected: []*utils.BuildData{},
			platform: &WindowsPlatform{Name: config.WindowsAmd64},
		},
		{
			name:     "test builds for mac",
			expected: []*utils.BuildData{},
			platform: &MacPlatform{Name: config.MacosxAmd64},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result, err := testCase.platform.GetBuilds()
			assert.NoError(t, err)
			assert.IsType(t, testCase.expected, result)
		})
	}
}

func TestGenerateBuildUrl(t *testing.T) {
	testCases := []struct {
		name     string
		input    *utils.BuildData
		expected string
		platform Platform
	}{
		{
			name: "test url generation for linux solc compiler",
			input: &utils.BuildData{
				Path:      "solc-linux-amd64-v0.4.10+commit.9e8cc01b",
				Version:   "0.4.10",
				Keccak256: "0x66bd5478b31c7ad1ec9b148618945dc8b7d0dc9ca4a2469992b9728daf672f9f",
				Sha256:    "0xf3638225df24f444a72123956033f5743079118f0e1195ce6969aa16a7ef2283",
			},
			expected: "https://binaries.soliditylang.org/linux-amd64/solc-linux-amd64-v0.4.10+commit.9e8cc01b",
			platform: &LinuxPlatform{Name: config.LinuxAmd64},
		},
		{
			name: "test url generation for OLD linux solc compiler",
			input: &utils.BuildData{
				Name:      "solc-v0.4.0",
				Version:   "0.4.0",
				Keccak256: "0xd13ec91fa7d5893619fbfa253ce56907ac4e55043cb3beffa5c06dc99aee3af5",
				Sha256:    "0xe26d188284763684f3cf6d4900b72f7e45a050dd2b2707320273529d033cfd47",
			},
			expected: "https://raw.githubusercontent.com/crytic/solc/master/linux/amd64/solc-v0.4.0",
			platform: &LinuxPlatform{Name: config.LinuxAmd64},
		},
		{
			name: "test url generation for mac solc compiler",
			input: &utils.BuildData{
				Path:      "solc-macosx-amd64-v0.3.6+commit.988fe5e5",
				Version:   "0.3.6",
				Keccak256: "0x662f94643bbc549d00fe7cfdbcdff504c78c697b10dcfca645d6082d464c1402",
				Sha256:    "0x02d581b5b373d8160b9e75d690dd2ee898c3e0f6a39bbda54cd0698224c09df4",
			},
			expected: "https://binaries.soliditylang.org/macosx-amd64/solc-macosx-amd64-v0.3.6+commit.988fe5e5",
			platform: &MacPlatform{Name: config.MacosxAmd64},
		},
		{
			name: "test url generation for windows solc compiler",
			input: &utils.BuildData{
				Path:      "solc-windows-amd64-v0.4.1+commit.4fc6fc2c.zip",
				Version:   "0.4.1",
				Keccak256: "0x8ad763849cff88a5e6446bc8d261d4983f993319fad8947538800316b22ed3e0",
				Sha256:    "0xe2815a517b24f6695b5f85002dd5b6ba095a327687708cf0d762db311600f6e9",
			},
			expected: "https://binaries.soliditylang.org/windows-amd64/solc-windows-amd64-v0.4.1+commit.4fc6fc2c.zip",
			platform: &WindowsPlatform{Name: config.WindowsAmd64},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result := testCase.platform.GenerateBuildUrl(testCase.input)
			assert.Equal(t, testCase.expected, result)
		})
	}
}

func TestGet(t *testing.T) {
	var expected *utils.ResponseData
	urls := []string{
		"https://binaries.soliditylang.org/macosx-amd64/list.json",
		"https://binaries.soliditylang.org/linux-amd64/list.json",
		"https://binaries.soliditylang.org/windows-amd64/list.json",
		"https://raw.githubusercontent.com/crytic/solc/new-list-json/linux/amd64/list.json",
	}

	for _, url := range urls {
		result, err := get(url)
		assert.NoError(t, err)
		assert.IsType(t, expected, result)
	}
}
