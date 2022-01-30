package installer

import (
	"fmt"
	"github.com/Masterminds/semver"
	"github.com/fabelx/go-solc-select/internal/errors"
	"github.com/fabelx/go-solc-select/internal/utils"
	"github.com/fabelx/go-solc-select/pkg/config"
	"github.com/fabelx/go-solc-select/pkg/versions"
	"os"
	"path/filepath"
	"reflect"
)

func downloadCompiler(platform versions.Platform, build *utils.BuildData) error {
	url := platform.GenerateBuildUrl(build)
	//resp, err := http.Get(url)
	//defer resp.Body.Close()
	data, err := utils.Get(url) // todo: Manage with this
	if err != nil {
		return err
	}

	name := fmt.Sprintf("solc-%s", build.Version)
	folder := filepath.Join(config.SolcArtifacts, name)
	if reflect.TypeOf(platform) == reflect.TypeOf(&versions.WindowsPlatform{}) && utils.IsOldWindowsVersion(build.Version) {
		err = utils.Unzip(folder, data)
		if err != nil {
			return err
		}

		return nil
	}

	err = os.MkdirAll(folder, os.ModePerm)
	if err != nil {
		return err
	}
	file, err := os.OpenFile(filepath.Join(folder, name), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0775)
	if err != nil {
		return err
	}

	defer file.Close()
	//_, err = io.Copy(file, resp.Body)  // todo: Manage with this
	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil

}

func getBuild(builds []*utils.BuildData, version *semver.Version) (*utils.BuildData, error) {
	for _, build := range builds {
		if build.Version == version.String() {
			return build, nil
		}
	}

	return nil, &errors.UnknownVersionError{Version: version.String()}
}

func InstallSolc(version string) error {
	v, err := semver.NewVersion(version)
	if err != nil {
		return err
	}

	platform, err := versions.GetPlatform()
	if err != nil {
		return err
	}

	builds, err := platform.GetBuilds()
	if err != nil {
		return err
	}

	build, err := getBuild(builds, v)
	if err != nil {
		return err
	}

	err = downloadCompiler(platform, build)
	if err != nil {
		return err
	}

	return nil
}

func InstallSolcs(versions []string) ([]string, []string) {
	var installed []string
	var notInstalled []string
	for _, version := range versions {
		err := InstallSolc(version)
		if err != nil {
			notInstalled = append(notInstalled, version)
			continue
		}

		installed = append(installed, version)
	}

	return installed, notInstalled
}
