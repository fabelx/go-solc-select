package uninstaller

import (
	"fmt"
	"github.com/fabelx/go-solc-select/pkg/config"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"path/filepath"
	"testing"
)

var testCurrentVersion = "0.4.7"
var testVersions = map[string]string{
	testCurrentVersion: testCurrentVersion,
	"0.5.1":            "0.5.1",
}

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
	err = os.WriteFile(config.CurrentVersionFilePath, []byte(testCurrentVersion), 0755)
	if err != nil {
		return err
	}

	// adds fake solc compilers
	for _, v := range testVersions {
		name := fmt.Sprintf("solc-%s", v)
		dirPath := filepath.Join(config.SolcArtifacts, name)
		os.Mkdir(dirPath, 0755)
		fakeFilePath := filepath.Join(dirPath, name)
		os.WriteFile(fakeFilePath, []byte(""), 0755)
	}
	return nil
}

// shutdown Removes all test files and dirs
func shutdown() {
	os.RemoveAll(config.SolcDir)
}

func TestUninstallSolc(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
		err      error
	}{
		{
			input:    testCurrentVersion,
			expected: testCurrentVersion,
			err:      nil,
		},
		{
			input:    "0.5.1",
			expected: "0.5.1",
			err:      nil,
		},
	}

	t.Run("test successful uninstalling", func(t *testing.T) {
		for _, testCase := range testCases {
			result, err := UninstallSolc(testCase.input)
			assert.NoError(t, err)
			assert.Equal(t, testCase.expected, result)
		}
	})

}

func TestUninstallSolcs(t *testing.T) {
	input := []string{"0.5.1", testCurrentVersion}
	expected := []string{"0.5.1", testCurrentVersion}

	t.Run("test successful uninstalling", func(t *testing.T) {
		result, err := UninstallSolcs(input)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})
}
