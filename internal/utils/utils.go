/*
	Copyright © 2022 Vladyslav Novotnyi <daprostovseeto@gmail.com>.

	fabelx/go-solc-select is licensed under the
	GNU Affero General Public License v3.0

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU Affero General Public License as
    published by the Free Software Foundation, either version 3 of the
    License, or (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU Affero General Public License for more details.
    You should have received a copy of the GNU Affero General Public License
    along with this program.  If not, see <https://www.gnu.org/licenses/>.

go-solc-select is a tool written in Golang for managing and switching between versions of the Solidity compiler.
*/

package utils

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"fmt"
	"github.com/Masterminds/semver"
	"github.com/fabelx/go-solc-select/internal/errors"
	"github.com/fabelx/go-solc-select/pkg/config"
	"golang.org/x/crypto/sha3"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
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

// IsOldLinuxVersion Determines if the compiler version for Linux is old
//
// Older versions are available at a different address
func IsOldLinuxVersion(version string) bool {
	ver, _ := semver.NewVersion(version)
	v, _ := semver.NewVersion("0.4.10")
	return ver.LessThan(v)
}

// IsOldWindowsVersion Determines if the compiler version for Windows is old
//
// Older versions have a different file structure
func IsOldWindowsVersion(version string) bool {
	ver, _ := semver.NewVersion(version)
	v, _ := semver.NewVersion("0.7.2")
	return ver.LessThan(v)
}

// Unzip Decompresses a file to a specific folder, returns an error on failure during decompression
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

// VerifyChecksum Checks the checksum of the received file, returns an error if the sum is incorrect
func VerifyChecksum(k256 string, s256 string, data []byte) error {
	if fmt.Sprintf("0x%x", sha256.Sum256(data)) != s256 {
		return &errors.ChecksumMismatchError{HashFunc: "Sha256", Platform: runtime.GOOS}
	}

	k := sha3.NewLegacyKeccak256()
	k.Write(data)
	if fmt.Sprintf("0x%x", k.Sum(nil)) != k256 {
		return &errors.ChecksumMismatchError{HashFunc: "Keccak256", Platform: runtime.GOOS}
	}

	return nil
}

// Clean Removes passed versions
func Clean(versions []string) {
	for _, version := range versions {
		folder := filepath.Join(config.SolcArtifacts, fmt.Sprintf("solc-%s", version))
		os.RemoveAll(folder)
	}
}
