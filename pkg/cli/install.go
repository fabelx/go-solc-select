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
	"github.com/fabelx/go-solc-select/pkg/installer"
	ver "github.com/fabelx/go-solc-select/pkg/versions"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install available solc versions",
	Long: `gsolc-select

Installs specific versions of the solc compiler.
You can specify multiple versions separated by spaces or 'all', which will install all available versions of the compiler.
`,
	Example: `gsolc-select install 0.8.1
gsolc-select install 0.8.1 0.4.23
gsolc-select install all`,
	Args: cobra.MinimumNArgs(1),
	Run:  installCompilers,
}

func installCompilers(cmd *cobra.Command, args []string) {
	var availableVersions, _ = ver.GetAvailable()
	fmt.Printf("There are %d versions of the solc compiler available for installation.\n", len(availableVersions))

	var installedVersions = ver.GetInstalled()
	var versionsToInstall []string
	for _, version := range args {

		if version == "all" {
			for key, _ := range availableVersions {
				versionsToInstall = append(versionsToInstall, key)
			}

			break
		}

		match := config.ValidSemVer.MatchString(version)
		if !match {
			fmt.Printf("Invalid version '%s'.\n", version)
			continue
		}

		if availableVersions[version] == "" {
			fmt.Printf("'%s' is not avaliable. Run `gsolc-select versions installable`.\n", version)
			continue
		}

		if installedVersions[version] != "" {
			fmt.Printf("Version '%s' is already installed. Run `gsolc-select versions`.\n", version)
			continue
		}

		versionsToInstall = append(versionsToInstall, version)
	}

	if len(versionsToInstall) == 0 {
		return
	}

	fmt.Printf("Installing %s...\n", versionsToInstall)
	installed, notInstalled, err := installer.InstallSolcs(versionsToInstall)
	if err != nil {
		fmt.Println(err) // todo: Exit?
		return
	}

	for _, version := range notInstalled {
		fmt.Printf("Failed to install version %s.\n", version)
	}

	for _, version := range installed {
		fmt.Printf("Version %s installed.\n", version)
	}
}

func init() {
	RegisterCmd(rootCmd, installCmd)
}
