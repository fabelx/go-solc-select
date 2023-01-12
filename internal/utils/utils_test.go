package utils

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"runtime"
	"testing"
)

func TestIsOldLinuxVersion(t *testing.T) {
	// 0.4.7 - is the old solc version for linux platform should return TRUE
	// 0.8.17 - is not the old solc version for linux platform should return FALSE
	assert.Equal(t, IsOldLinuxVersion("0.4.7"), true)
	assert.Equal(t, IsOldLinuxVersion("0.8.17"), false)
}

func TestIsOldWindowsVersion(t *testing.T) {
	// 0.7.0 - is the old solc version for windows platform should return TRUE
	// 0.8.17 - is not the old solc version for windows platform should return FALSE
	assert.Equal(t, IsOldWindowsVersion("0.7.0"), true)
	assert.Equal(t, IsOldWindowsVersion("0.8.17"), false)
}

func TestVerifyChecksum(t *testing.T) {
	// sha256(supernatural) == 0xc1b59ba0c56e6b217dfdd0de59e40718568517c076a012c4ea34ba867f08c5e2
	// keccak256(supernatural) == 0x07b3b30205e8fa3a189c8db0a754fc8e0ac331f3a02bed7626688fb27aa05891

	// should pass
	err := VerifyChecksum(
		"0x07b3b30205e8fa3a189c8db0a754fc8e0ac331f3a02bed7626688fb27aa05891",
		"0xc1b59ba0c56e6b217dfdd0de59e40718568517c076a012c4ea34ba867f08c5e2",
		[]byte("supernatural"),
	)

	assert.NoError(t, err)

	// should fail with wrong keccak256
	err = VerifyChecksum(
		"0x0000030205e8fa3a189c8db0a754fc8e0ac331f3a02bed7626688fb27aa00000",
		"0xc1b59ba0c56e6b217dfdd0de59e40718568517c076a012c4ea34ba867f08c5e2",
		[]byte("supernatural"),
	)

	assert.NotNil(t, err)
	assert.EqualError(t, err, fmt.Sprintf("%s checksum mismatch of files for %s platform.", "Keccak256", runtime.GOOS))

	// should fail with wrong sha256
	err = VerifyChecksum(
		"0x07b3b30205e8fa3a189c8db0a754fc8e0ac331f3a02bed7626688fb27aa05891",
		"0x0000030205e8fa3a189c8db0a754fc8e0ac331f3a02bed7626688fb27aa00000",
		[]byte("supernatural"),
	)

	assert.NotNil(t, err)
	assert.EqualError(t, err, fmt.Sprintf("%s checksum mismatch of files for %s platform.", "Sha256", runtime.GOOS))
}

func TestGet(t *testing.T) {
	testCases := []struct {
		name     string
		server   *httptest.Server
		response []byte
	}{
		{
			name: "http request with StatusOK",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("okay"))
			})),
			response: []byte("okay"),
		},
		{
			name: "http request with StatusBadRequest",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
			})),
			response: nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			data, err := Get(testCase.server.URL)

			assert.Equal(t, testCase.response, data)
			if err != nil {
				assert.EqualError(t, err, fmt.Sprintf("Recieved unexpected status code: '%d' from '%s' request.", http.StatusBadRequest, testCase.server.URL))
			}
		})
	}
}

func TestUnzip(t *testing.T) {
	// todo: write zip tests
}
