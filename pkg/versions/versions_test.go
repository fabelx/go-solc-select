package versions

import (
	"github.com/fabelx/go-solc-select/pkg/config"
	"reflect"
	"testing"
)

func TestGetAvailable(t *testing.T) {
	testTable := []struct {
		expected map[string]string
	}{
		{
			expected: map[string]string{},
		},
	}

	for _, testCase := range testTable {
		result, _ := GetAvailable()
		for key, value := range result {
			t.Log(key, value)
		}
		if !(reflect.TypeOf(result) == reflect.TypeOf(testCase.expected)) {
			t.Errorf("Unexpected result! Ecpected %v - got %v", testCase.expected, result)
		}
	}
}

func TestLinuxGetAvailableVersions(t *testing.T) {
	testTable := []struct {
		expected map[string]string
	}{
		{
			expected: map[string]string{},
		},
	}

	for _, testCase := range testTable {
		platform := &WindowsPlatform{Name: config.WindowsAmd64}
		result, _ := platform.GetAvailableVersions()
		for key, value := range result {
			t.Log(key, value)
		}
		if !(reflect.TypeOf(result) == reflect.TypeOf(testCase.expected)) {
			t.Errorf("Unexpected result! Ecpected %v - got %v", testCase.expected, result)
		}
	}
}

func TestMacGetAvailableVersions(t *testing.T) {
	testTable := []struct {
		expected map[string]string
	}{
		{
			expected: map[string]string{},
		},
	}

	for _, testCase := range testTable {
		platform := &WindowsPlatform{Name: config.WindowsAmd64}
		result, _ := platform.GetAvailableVersions()
		for key, value := range result {
			t.Log(key, value)
		}
		if !(reflect.TypeOf(result) == reflect.TypeOf(testCase.expected)) {
			t.Errorf("Unexpected result! Ecpected %v - got %v", testCase.expected, result)
		}
	}
}

func TestWindowsGetAvailableVersions(t *testing.T) {
	testTable := []struct {
		expected map[string]string
	}{
		{
			expected: map[string]string{},
		},
	}

	for _, testCase := range testTable {
		platform := &WindowsPlatform{Name: config.WindowsAmd64}
		result, _ := platform.GetAvailableVersions()
		for key, value := range result {
			t.Log(key, value)
		}
		if !(reflect.TypeOf(result) == reflect.TypeOf(testCase.expected)) {
			t.Errorf("Unexpected result! Ecpected %v - got %v", testCase.expected, result)
		}
	}
}
