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

package versions

import (
	"encoding/json"
	"fmt"
	"github.com/fabelx/go-solc-select/internal/errors"
	"github.com/fabelx/go-solc-select/internal/utils"
	"github.com/fabelx/go-solc-select/pkg/config"
	"os"
	"path/filepath"
	"strings"
)

type Platform interface {
	GetAvailableVersions() (map[string]string, error)
	GetBuilds() ([]*utils.BuildData, error)
	GenerateBuildUrl(build *utils.BuildData) string
}

type LinuxPlatform struct {
	Name string `json:"name"`
}

type MacPlatform struct {
	Name string `json:"name"`
}

type WindowsPlatform struct {
	Name string `json:"name"`
}

func get(url string) (*utils.ResponseData, error) {
	data, err := utils.Get(url)
	if err != nil {
		return nil, err
	}

	respData := utils.ResponseData{}
	err = json.Unmarshal(data, &respData)
	if err != nil {
		return nil, err
	}

	return &respData, nil
}

// GetAvailableVersions Returns an array of compiler versions for linux
func (r *LinuxPlatform) GetAvailableVersions() (map[string]string, error) {
	versions, err := getVersions(fmt.Sprintf("%s/%s/list.json", config.SoliditylangUrl, r.Name))
	if err != nil {
		return nil, err
	}

	oldVersions, err := getVersions(config.OldSolcListUrl)
	if err != nil {
		return nil, err
	}

	for key, value := range oldVersions {
		versions[key] = value
	}

	return versions, nil
}

// GetAvailableVersions Returns an array of compiler versions for mac
func (r *MacPlatform) GetAvailableVersions() (map[string]string, error) {
	return getVersions(fmt.Sprintf("%s/%s/list.json", config.SoliditylangUrl, r.Name))
}

// GetAvailableVersions Returns an array of compiler versions for windows
func (r *WindowsPlatform) GetAvailableVersions() (map[string]string, error) {
	return getVersions(fmt.Sprintf("%s/%s/list.json", config.SoliditylangUrl, r.Name))
}

// GetBuilds Returns an array of compiler versions
func getVersions(url string) (map[string]string, error) {
	respData, err := get(url)
	if err != nil {
		return nil, err
	}

	return respData.Releases, nil
}

// GetBuilds Returns an array of meta information about compilers for linux
func (r *LinuxPlatform) GetBuilds() ([]*utils.BuildData, error) {
	builds, err := getBuilds(fmt.Sprintf("%s/%s/list.json", config.SoliditylangUrl, r.Name))
	if err != nil {
		return nil, err
	}

	oldBuilds, err := getBuilds(config.OldSolcListUrl)
	if err != nil {
		return nil, err
	}

	builds = append(builds, oldBuilds...)

	return builds, nil
}

// GetBuilds Returns an array of meta information about compilers for mac
func (r *MacPlatform) GetBuilds() ([]*utils.BuildData, error) {
	return getBuilds(fmt.Sprintf("%s/%s/list.json", config.SoliditylangUrl, r.Name))
}

// GetBuilds Returns an array of meta information about compilers for windows
func (r *WindowsPlatform) GetBuilds() ([]*utils.BuildData, error) {
	return getBuilds(fmt.Sprintf("%s/%s/list.json", config.SoliditylangUrl, r.Name))
}

// getBuilds Returns an array of meta information about compilers
func getBuilds(url string) ([]*utils.BuildData, error) {
	respData, err := get(url)
	if err != nil {
		return nil, err
	}

	return respData.Builds, nil
}

// GenerateBuildUrl Returns the url of solc compiler file(s) for linux
func (r *LinuxPlatform) GenerateBuildUrl(build *utils.BuildData) string {
	if utils.IsOldLinuxVersion(build.Version) {
		return fmt.Sprintf("%s/%s", config.OldSolcUrl, build.Name)
	}
	return fmt.Sprintf("%s/%s/%s", config.SoliditylangUrl, r.Name, build.Path)
}

// GenerateBuildUrl Returns the url of solc compiler file(s) for mac
func (r *MacPlatform) GenerateBuildUrl(build *utils.BuildData) string {
	return fmt.Sprintf("%s/%s/%s", config.SoliditylangUrl, r.Name, build.Path)
}

// GenerateBuildUrl Returns the url of solc compiler file(s) for windows
func (r *WindowsPlatform) GenerateBuildUrl(build *utils.BuildData) string {
	return fmt.Sprintf("%s/%s/%s", config.SoliditylangUrl, r.Name, build.Path)
}

// GetInstalled Returns installed versions on system
func GetInstalled() (map[string]string, error) {
	matches, err := filepath.Glob(filepath.Join(config.SolcArtifacts, "solc-*"))
	if err != nil {
		return nil, err
	}

	versions := make(map[string]string)
	for _, path := range matches {
		version := strings.Replace(filepath.Base(path), "solc-", "", 1)
		versions[version] = version
	}

	return versions, nil
}

// GetAvailable Returns all installable versions of the solc compiler for the current operating system platform
func GetAvailable() (map[string]string, error) {
	platform, err := GetPlatform()
	if err != nil {
		return nil, err
	}

	return platform.GetAvailableVersions()
}

// GetCurrent Returns current version on system
func GetCurrent() (string, error) {
	// Getting the compiler version from a file where the version is specified
	data, err := os.ReadFile(config.CurrentVersionFilePath)
	if err != nil {
		return "", err
	}

	version := string(data)

	if version != "" {
		installedVersions, err := GetInstalled()
		if err != nil {
			return "", err
		}

		// Checking if the compiler is installed on the host
		if installedVersions[version] == "" {
			return "", &errors.NotInstalledError{Version: version}
		}

		return version, nil
	} else {
		return "", &errors.NoCompilerSelected{}
	}

}
