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
	"github.com/fabelx/go-solc-select/pkg/config"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "The go-solc-select version",
	Long: `gsolc-select

Displays the version of this gsolc-select binary and exits.
`,
	Args: cobra.NoArgs,
	Run:  getVersion,
}

func getVersion(cmd *cobra.Command, args []string) {
	fmt.Println(fmt.Sprintf("Version: %s", config.GoSolcSelect))

}

func init() {
	RegisterCmd(rootCmd, versionCmd)
}
