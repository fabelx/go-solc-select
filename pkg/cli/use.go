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
	"github.com/fabelx/go-solc-select/pkg/switcher"
	"github.com/fabelx/go-solc-select/pkg/versions"
	"github.com/spf13/cobra"
)

var (
	install bool
)

var useCmd = &cobra.Command{
	Use:   "use",
	Short: "Change the version of global solc compiler",
	Long: `gsolc-select

Switch between installed versions of solc compiler. 
Using the -i / --installer flag automatically installer the required compiler version.
`,
	Example: `gsolc-select use 0.4.12
gsolc-select use -i 0.4.13`,
	Args: cobra.ExactArgs(1),
	Run:  useCompiler,
}

func useCompiler(cmd *cobra.Command, args []string) {
	version := args[0]
	match := config.ValidSemVer.MatchString(version)
	if !match {
		fmt.Printf("Invalid version '%s'.\n", version)
		return
	}
	var availableVersions, _ = versions.GetAvailable()
	if availableVersions[version] == "" {
		fmt.Printf("'%s' is not avaliable. Run `solc-select versions installable`.\n", version)
		return
	}

	if install {
		installCompilers(cmd, args)
	}

	err := switcher.SwitchSolc(version)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Switched global version to '%s'.\n", version)
}

func init() {
	useCmd.Flags().BoolVarP(&install, "install", "i", false, "indicate if you want to automatically installer versions that are not installed")
	RegisterCmd(rootCmd, useCmd)
}
