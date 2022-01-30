package installer

import (
	"reflect"
	"testing"
)

func TestInstallSolc(t *testing.T) {
	testTable := []struct {
		input    string
		expected error
	}{
		{
			input:    "0.7.4",
			expected: nil,
		},
	}

	for _, testCase := range testTable {
		result := InstallSolc(testCase.input)
		if !(reflect.TypeOf(result) == reflect.TypeOf(testCase.expected)) {
			t.Errorf("Unexpected result! Ecpected %v - got %v", testCase.expected, result)
		}
	}
}

func TestInstallSolcs(t *testing.T) {
	testTable := []struct {
		input        []string
		installed    []string
		notInstalled []string
	}{
		{
			input:        []string{"0.4.11", "0.7.4", "0.7.5"},
			installed:    []string{"0.4.11", "0.7.4", "0.7.5"},
			notInstalled: nil,
		},
	}

	for _, testCase := range testTable {
		installed, _ := InstallSolcs(testCase.input)
		if len(installed) != len(testCase.installed) {
			t.Errorf("Unexpected result! Ecpected %v - got %v", testCase.installed, installed)
		}
	}
}
