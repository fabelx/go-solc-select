package utils

import (
	"archive/zip"
	"bytes"
	"github.com/Masterminds/semver"
	"github.com/fabelx/go-solc-select/internal/errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

type ResponseData struct {
	Releases map[string]string `json:"releases"`
	Builds   []*BuildData      `json:"builds"`
}

type BuildData struct {
	Path      string `json:"path"`
	Name      string `json:"name"`
	Version   string `json:"version"`
	Keccak256 string `json:"keccak256"`
	Sha256    string `json:"sha256"`
}

// Get Base implementation of request
func Get(url string) ([]byte, error) {
	r, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer r.Body.Close()
	if r.StatusCode == http.StatusOK {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}

		return data, nil
	}

	return nil, &errors.UnexpectedStatusCode{StatusCode: r.StatusCode, Url: url}
}

func IsOldLinuxVersion(version string) bool {
	ver, _ := semver.NewVersion(version)
	v, _ := semver.NewVersion("0.4.10")
	return ver.LessThan(v)
}

func IsOldWindowsVersion(version string) bool {
	ver, _ := semver.NewVersion(version)
	v, _ := semver.NewVersion("0.7.2")
	return ver.LessThan(v)
}

func Unzip(folder string, data []byte) error {
	zipReader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return err
	}

	for _, archiveReader := range zipReader.File {

		// Open the file in the archive
		archiveFile, err := archiveReader.Open()
		if err != nil {
			return err
		}
		defer archiveFile.Close()

		// Prepare to write the file
		finalPath := filepath.Join(folder, archiveReader.Name)

		// Check if the file to extract is just a directory
		if archiveReader.FileInfo().IsDir() {
			err = os.MkdirAll(finalPath, 0755)
			if err != nil {
				return err
			}
			// Continue to the next file in the archive
			continue
		}

		// Create all needed directories
		if os.MkdirAll(filepath.Dir(finalPath), 0755) != nil {
			return err
		}

		// Prepare to write the destination file
		destinationFile, err := os.OpenFile(finalPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, archiveReader.Mode())
		if err != nil {
			return err
		}
		defer destinationFile.Close()

		// Write the destination file
		if _, err = io.Copy(destinationFile, archiveFile); err != nil {
			return err
		}
	}

	return nil
}
