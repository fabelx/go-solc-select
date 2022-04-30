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
	"fmt"
	"github.com/fabelx/go-solc-select/pkg/versions"
	"github.com/spf13/cobra"
)

var currentCmd = &cobra.Command{
	Use:   "current",
	Short: "Current solc versions",
	Long: `gsolc-select

Prints out current solc versions and exit.
`,
	Args: cobra.NoArgs,
	Run:  getCurrentVersions,
}

func getCurrentVersions(cmd *cobra.Command, args []string) {
	var currentVersion, err = versions.GetCurrent()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%s\n", currentVersion)

}

func init() {
	RegisterCmd(versionCmd, currentCmd)
}
