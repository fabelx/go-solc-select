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

package installer

import (
	"fmt"
	"github.com/fabelx/go-solc-select/internal/utils"
	"github.com/fabelx/go-solc-select/pkg/config"
	ver "github.com/fabelx/go-solc-select/pkg/versions"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"sync"
)

// download Returns an error if downloading of the solc compiler fails
func download(platform ver.Platform, build *utils.BuildData) error {
	url := platform.GenerateBuildUrl(build)
	data, err := utils.Get(url)
	if err != nil {
		return err
	}

	// Verifying checksum of files
	err = utils.VerifyChecksum(build.Keccak256, build.Sha256, data)
	if err != nil {
		return err
	}

	name := fmt.Sprintf("solc-%s", build.Version)
	folder := filepath.Join(config.SolcArtifacts, name)

	// Old compiler versions (<0.7.2) for windows have a different file structure
	if reflect.TypeOf(platform) == reflect.TypeOf(&ver.WindowsPlatform{}) && utils.IsOldWindowsVersion(build.Version) {
		err = utils.Unzip(folder, data)
		if err != nil {
			return err
		}

		// Rename bin file solc.exe -> solc-[version] (solc-0.0.0)
		filePath := filepath.Join(folder, "solc.exe")
		newFilePath := filepath.Join(folder, name)
		err := os.Rename(filePath, newFilePath)
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

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil

}

// InstallSolc Returns compiler meta information if the installation completed successfully
func InstallSolc(platform ver.Platform, build *utils.BuildData) (*utils.BuildData, error) { // todo: what is the reason to return *BuildData?
	err := download(platform, build)
	if err != nil {
		return nil, err
	}

	return build, nil
}

// InstallSolcs Returns slice of installed compiler versions, slice of NOT installed compiler versions and error
func InstallSolcs(versions []string) ([]string, []string, error) {
	platform, err := ver.GetPlatform(runtime.GOOS)
	if err != nil {
		return nil, nil, err
	}

	builds, err := platform.GetBuilds()
	if err != nil {
		return nil, nil, err
	}

	var installed []string
	var notInstalled []string
	var buildsToInstall []*utils.BuildData

	for _, version := range versions {
		build, err := ver.GetBuild(builds, version)
		if err != nil {
			notInstalled = append(notInstalled, version)
			continue
		}

		buildsToInstall = append(buildsToInstall, build)
	}

	// Install solc compilers
	wg := sync.WaitGroup{}
	for _, build := range buildsToInstall {
		wg.Add(1)
		build := build
		go func() {
			defer wg.Done()
			_, err := InstallSolc(platform, build)
			if err != nil {
				notInstalled = append(notInstalled, build.Version)
				return
			}

			installed = append(installed, build.Version)
		}()
	}

	wg.Wait()

	return installed, notInstalled, nil
}
