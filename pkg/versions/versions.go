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

func (r *MacPlatform) GetAvailableVersions() (map[string]string, error) {
	return getVersions(fmt.Sprintf("%s/%s/list.json", config.SoliditylangUrl, r.Name))
}

func (r *WindowsPlatform) GetAvailableVersions() (map[string]string, error) {
	return getVersions(fmt.Sprintf("%s/%s/list.json", config.SoliditylangUrl, r.Name))
}

func getVersions(url string) (map[string]string, error) {
	data, err := utils.Get(url)
	if err != nil {
		return nil, err
	}

	respData := utils.ResponseData{}
	err = json.Unmarshal(data, &respData)
	if err != nil {
		return nil, err
	}

	return respData.Releases, nil
}

func (r LinuxPlatform) GetBuilds() ([]*utils.BuildData, error) {
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

func (r MacPlatform) GetBuilds() ([]*utils.BuildData, error) {
	return getBuilds(fmt.Sprintf("%s/%s/list.json", config.SoliditylangUrl, r.Name))
}

func (r WindowsPlatform) GetBuilds() ([]*utils.BuildData, error) {
	return getBuilds(fmt.Sprintf("%s/%s/list.json", config.SoliditylangUrl, r.Name))
}

func getBuilds(url string) ([]*utils.BuildData, error) {
	data, err := utils.Get(url)
	if err != nil {
		return nil, err
	}

	respData := utils.ResponseData{}
	err = json.Unmarshal(data, &respData)
	if err != nil {
		return nil, err
	}

	return respData.Builds, nil
}

func (r LinuxPlatform) GenerateBuildUrl(build *utils.BuildData) string {
	if utils.IsOldLinuxVersion(build.Version) {
		return fmt.Sprintf("%s/%s", config.OldSolcUrl, build.Name)
	}
	return fmt.Sprintf("%s/%s/%s", config.SoliditylangUrl, r.Name, build.Path)
}

func (r MacPlatform) GenerateBuildUrl(build *utils.BuildData) string {
	return fmt.Sprintf("%s/%s/%s", config.SoliditylangUrl, r.Name, build.Path)
}

func (r WindowsPlatform) GenerateBuildUrl(build *utils.BuildData) string {
	return fmt.Sprintf("%s/%s/%s", config.SoliditylangUrl, r.Name, build.Path)
}

// GetInstalled Returns installed versions on system
func GetInstalled() (map[string]string, error) {
	matches, err := filepath.Glob(filepath.Join(config.SolcArtifacts, "**", "solc-*")) // todo: Improve pattern
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
	version := os.Getenv(config.EnvVariable)
	if version != "" {
		installedVersions, err := GetInstalled()
		if err != nil {
			return "", err
		}

		if installedVersions[version] == "" {
			return "", &errors.NotInstalledError{Version: version}
		}

		return version, err
	}

	data, err := os.ReadFile(config.CurrentVersionFilePath)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
