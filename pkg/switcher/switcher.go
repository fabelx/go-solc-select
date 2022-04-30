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

package switcher

import (
	"github.com/fabelx/go-solc-select/internal/errors"
	"github.com/fabelx/go-solc-select/pkg/config"
	"github.com/fabelx/go-solc-select/pkg/versions"
	"os"
)

// SwitchSolc Returns an error if the version switch failed
func SwitchSolc(version string) error {
	installedVersions, err := versions.GetInstalled()
	if err != nil {
		return err
	}

	if installedVersions[version] == "" {
		return &errors.NotInstalledError{Version: version}
	}

	file, err := os.OpenFile(config.CurrentVersionFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(version)
	if err != nil {
		return err
	}

	os.Setenv("SOLC_VERSION", version)

	return nil
}
