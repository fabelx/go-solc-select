package installer

import (
	"fmt"
	"github.com/fabelx/go-solc-select/internal/errors"
	"github.com/fabelx/go-solc-select/internal/utils"
	"github.com/fabelx/go-solc-select/pkg/config"
	"github.com/fabelx/go-solc-select/pkg/versions"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"path/filepath"
	"testing"
)

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

	return nil
}

// shutdown Removes all test files and dirs
func shutdown() {
	os.RemoveAll(config.SolcDir)
}

func TestInstallSolc(t *testing.T) {
	testCases := []struct {
		name     string
		input    *utils.BuildData
		expected *utils.BuildData
		err      error
	}{
		{
			name: "test success install",
			input: &utils.BuildData{
				Path:      "solc-windows-amd64-v0.4.1+commit.4fc6fc2c.zip",
				Version:   "0.4.1",
				Keccak256: "0x8ad763849cff88a5e6446bc8d261d4983f993319fad8947538800316b22ed3e0",
				Sha256:    "0xe2815a517b24f6695b5f85002dd5b6ba095a327687708cf0d762db311600f6e9",
			},
			expected: &utils.BuildData{
				Path:      "solc-windows-amd64-v0.4.1+commit.4fc6fc2c.zip",
				Version:   "0.4.1",
				Keccak256: "0x8ad763849cff88a5e6446bc8d261d4983f993319fad8947538800316b22ed3e0",
				Sha256:    "0xe2815a517b24f6695b5f85002dd5b6ba095a327687708cf0d762db311600f6e9",
			},
			err: nil,
		},
		{
			name: "test failed install - wrong version",
			input: &utils.BuildData{
				Path:      "solc-windows-amd64-v0.4.111+commit.4fc6fc2c.zip",
				Version:   "0.4.111",
				Keccak256: "0x8ad763849cff88a5e6446bc8d261d4983f993319fad8947538800316b22ed3e0",
				Sha256:    "0xe2815a517b24f6695b5f85002dd5b6ba095a327687708cf0d762db311600f6e9",
			},
			expected: nil,
			err: &errors.UnexpectedStatusCode{
				StatusCode: 404,
				Url:        "https://binaries.soliditylang.org/windows-amd64/solc-windows-amd64-v0.4.111+commit.4fc6fc2c.zip",
			},
		},
		{
			name: "test failed install - wrong Keccak256",
			input: &utils.BuildData{
				Path:      "solc-windows-amd64-v0.4.2+commit.af6afb04.zip",
				Version:   "0.4.2",
				Keccak256: "0x8ad763849cff88a7e6446bc8d261d4983f993319fad8947538800316b22ed3e0",
				Sha256:    "0x34e10611651cbe9c2d7b8b4d1cc94779fc80d52a6c6975e308384308fe117eb9",
			},
			expected: nil,
			err:      &errors.ChecksumMismatchError{HashFunc: "Keccak256", Platform: "windows"},
		},
		{
			name: "test failed install - wrong Sha256",
			input: &utils.BuildData{
				Path:      "solc-windows-amd64-v0.4.2+commit.af6afb04.zip",
				Version:   "0.4.2",
				Keccak256: "0xe45a3d296656d66cdf9e7c5eec47b37afe260b9eed81dcbf60717b5c7b388e08",
				Sha256:    "0x34e10611651cbe9c8d7b8b4d1cc94779fc80d52a6c6975e308384308fe117eb9",
			},
			expected: nil,
			err:      &errors.ChecksumMismatchError{HashFunc: "Sha256", Platform: "windows"},
		},
	}

	platform, err := versions.GetPlatform("windows")
	assert.NoError(t, err)
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result, err := InstallSolc(platform, testCase.input)
			if result != nil {
				name := fmt.Sprintf("solc-%s", testCase.expected.Version)
				assert.FileExists(t, filepath.Join(config.SolcArtifacts, name, name))
			}

			assert.Equal(t, testCase.expected, result)
			assert.Equal(t, testCase.err, err)
		})
	}
}

func TestInstallSolcs(t *testing.T) {
	testCases := []struct {
		input                []string
		expectedInstalled    []string
		expectedNotInstalled []string
	}{
		{
			input:                []string{"0.7.1", "0.8.3"},
			expectedInstalled:    []string{"0.7.1", "0.8.3"},
			expectedNotInstalled: []string(nil),
		},
		{
			input:                []string{"0.0.0", "0.5.7"},
			expectedInstalled:    []string{"0.5.7"},
			expectedNotInstalled: []string{"0.0.0"},
		},
	}

	for _, testCase := range testCases {
		resultInstalled, resultNotInstalled, err := InstallSolcs(testCase.input)
		assert.NoError(t, err)
		assert.Equal(t, testCase.expectedInstalled, resultInstalled)
		assert.Equal(t, testCase.expectedNotInstalled, resultNotInstalled)
	}
}