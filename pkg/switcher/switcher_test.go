package switcher

import (
	"fmt"
	"github.com/fabelx/go-solc-select/internal/errors"
	"github.com/fabelx/go-solc-select/pkg/config"
	"github.com/fabelx/go-solc-select/pkg/versions"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"path/filepath"
	"testing"
)

var testCurrentVersion = "0.4.9"
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

func TestSwitchSolc(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
		err      error
	}{
		{
			name:     "test failed switch version",
			input:    "0.0.0",
			expected: "",
			err:      &errors.NotInstalledError{Version: "0.0.0"},
		},
		{
			name:     "test success switch version",
			input:    testCurrentVersion,
			expected: testCurrentVersion,
			err:      nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			err := SwitchSolc(testCase.input)
			if err == nil {
				version, err2 := versions.GetCurrent()
				assert.NoError(t, err2)
				assert.Equal(t, testCase.expected, version)
			}

			assert.Equal(t, testCase.err, err)
		})
	}
}
