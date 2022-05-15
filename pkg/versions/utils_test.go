package versions

import (
	"fmt"
	"github.com/Masterminds/semver"
	"github.com/fabelx/go-solc-select/internal/errors"
	"github.com/fabelx/go-solc-select/pkg/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetPlatform(t *testing.T) {
	testCases := []struct {
		input    string
		expected Platform
		err      error
	}{
		{
			input:    "darwin",
			expected: &MacPlatform{Name: config.MacosxAmd64},
			err:      nil,
		},
		{
			input:    "linux",
			expected: &LinuxPlatform{Name: config.LinuxAmd64},
			err:      nil,
		},
		{
			input:    "windows",
			expected: &WindowsPlatform{Name: config.WindowsAmd64},
			err:      nil,
		},
	}

	t.Run("test platform selection", func(t *testing.T) {
		for _, testCase := range testCases {
			result, err := GetPlatform(testCase.input)
			assert.NoError(t, err)
			assert.Equal(t, testCase.expected, result)
		}
	})

	t.Run("test unsupported platform", func(t *testing.T) {
		input := "Unsupported"
		expectedErr := &errors.UnsupportedPlatformError{Platform: input}
		expectedErrMsg := fmt.Sprintf("'%s' platform is not currently supported.", input)
		_, err := GetPlatform(input)

		assert.IsType(t, expectedErr, err)
		assert.EqualError(t, err, expectedErrMsg) // todo: move to error tests
	})
}

func TestSortVersions(t *testing.T) {
	t.Run("test version sorting", func(t *testing.T) {
		input := map[string]string{
			"0.0.0":  "0.0.0",
			"0.0.9":  "0.0.0",
			"0.5.0":  "0.0.0",
			"0.3.1":  "0.0.0",
			"0.8.3":  "0.0.0",
			"0.3.10": "0.0.0",
			"0.9.9":  "0.0.0",
			"0.1.18": "0.0.0",
		}
		expected := []*semver.Version{
			semver.MustParse("0.0.0"),
			semver.MustParse("0.0.9"),
			semver.MustParse("0.1.18"),
			semver.MustParse("0.3.1"),
			semver.MustParse("0.3.10"),
			semver.MustParse("0.5.0"),
			semver.MustParse("0.8.3"),
			semver.MustParse("0.9.9"),
		}
		result := SortVersions(input)

		assert.Equal(t, expected, result)
	})
}
