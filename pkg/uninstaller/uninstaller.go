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

package uninstaller

import (
	"fmt"
	"github.com/fabelx/go-solc-select/pkg/config"
	ver "github.com/fabelx/go-solc-select/pkg/versions"
	"os"
	"path/filepath"
)

// UninstallSolc Returns nil if success
func UninstallSolc(version string) error {
	// reset the current version in the file if it gets deleted
	var currentVersion, _ = ver.GetCurrent()
	if currentVersion == version {
		os.WriteFile(config.CurrentVersionFilePath, []byte(""), 0755)
	}

	// remove a dir with solc compiler artifacts
	folderPath := filepath.Join(config.SolcArtifacts, fmt.Sprintf("solc-%s", version))
	err := os.RemoveAll(folderPath)
	if err != nil {
		return err
	}

	return nil
}

// UninstallSolcs performs sequentially uninstallation of compilers
// Returns slice of uninstalled compiler versions, slice of NOT uninstalled compiler versions and error
func UninstallSolcs(versions []string) ([]string, []string, error) {
	var uninstalled []string
	var notUninstalled []string
	for _, version := range versions {
		err := UninstallSolc(version)
		if err != nil {
			notUninstalled = append(notUninstalled, version)
			continue
		}

		uninstalled = append(uninstalled, version)
	}

	return uninstalled, notUninstalled, nil
}
