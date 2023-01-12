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

package cli

import (
	"github.com/fabelx/go-solc-select/pkg/config"
	ver "github.com/fabelx/go-solc-select/pkg/versions"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

var currentCmd = &cobra.Command{
	Use:   "current",
	Short: "Current solc version",
	Long: `gsolc-select

Prints out current solc versions and exit.
`,
	Example: `  gsolc-select versions current`,
	Args:    cobra.NoArgs,
	RunE:    getCurrentVersions,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// Checks if there is a file to store the current version of the compiler
		// - `global-version` file. File:<$HomeDir/.gsolc-select/global-version>
		if _, err := os.Stat(config.CurrentVersionFilePath); os.IsNotExist(err) {
			// Creates file if it doesn't exist
			err = os.WriteFile(config.CurrentVersionFilePath, []byte(""), 0755)
			if err != nil {
				return err
			}
		}

		return nil
	},
}

func getCurrentVersions(cmd *cobra.Command, args []string) error {
	var currentVersion, err = ver.GetCurrent()
	if err != nil {
		return err
	}

	log.Warn(currentVersion)
	return nil
}

func init() {
	RegisterCmd(versionsCmd, currentCmd)
}
