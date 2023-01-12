package installer

import (
	"context"
	"fmt"
	"github.com/fabelx/go-solc-select/internal/errors"
	"github.com/fabelx/go-solc-select/internal/utils"
	"github.com/fabelx/go-solc-select/pkg/config"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestMain(m *testing.M) {

	if os.Getenv("CI") != "" {
		return
	}

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
		expected error
	}{
		{
			name: "test success install",
			input: &utils.BuildData{
				Path:      "solc-windows-amd64-v0.4.1+commit.4fc6fc2c.zip",
				Version:   "0.4.1",
				Keccak256: "0x8ad763849cff88a5e6446bc8d261d4983f993319fad8947538800316b22ed3e0",
				Sha256:    "0xe2815a517b24f6695b5f85002dd5b6ba095a327687708cf0d762db311600f6e9",
			},
			expected: nil,
		},
		{
			name: "test failed install - wrong version",
			input: &utils.BuildData{
				Path:      "solc-windows-amd64-v0.4.111+commit.4fc6fc2c.zip",
				Version:   "0.4.111",
				Keccak256: "0x8ad763849cff88a5e6446bc8d261d4983f993319fad8947538800316b22ed3e0",
				Sha256:    "0xe2815a517b24f6695b5f85002dd5b6ba095a327687708cf0d762db311600f6e9",
			},
			expected: &errors.UnknownVersionError{Version: "0.4.111"},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			err := InstallSolc(testCase.input.Version)
			if err == nil {
				name := fmt.Sprintf("solc-%s", testCase.input.Version)
				assert.FileExists(t, filepath.Join(config.SolcArtifacts, name, name))
			}

			assert.Equal(t, testCase.expected, err)
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
		resultInstalled, resultNotInstalled, err := InstallSolcs(context.Background(), testCase.input)
		assert.NoError(t, err)
		for _, installed := range resultInstalled {
			name := fmt.Sprintf("solc-%s", installed)
			assert.FileExists(t, filepath.Join(config.SolcArtifacts, name, name))
		}
		assert.ElementsMatch(t, testCase.expectedInstalled, resultInstalled)
		assert.ElementsMatch(t, testCase.expectedNotInstalled, resultNotInstalled)
	}
}

func TestAsyncInstallSolcs(t *testing.T) {
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
		resultInstalled, resultNotInstalled, err := AsyncInstallSolcs(context.Background(), testCase.input)
		assert.NoError(t, err)
		for _, installed := range resultInstalled {
			name := fmt.Sprintf("solc-%s", installed)
			assert.FileExists(t, filepath.Join(config.SolcArtifacts, name, name))
		}
		assert.ElementsMatch(t, testCase.expectedInstalled, resultInstalled)
		assert.ElementsMatch(t, testCase.expectedNotInstalled, resultNotInstalled)
	}
}
