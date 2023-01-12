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
	"github.com/fabelx/go-solc-select/pkg/uninstaller"
	ver "github.com/fabelx/go-solc-select/pkg/versions"
	"github.com/spf13/cobra"
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Remove installed solc versions",
	Long: `gsolc-select

Removes certain versions of the installed solc compiler.
You can specify multiple versions separated by spaces or 'all', which will remove all installed versions of the compiler.
`,
	Example: `  gsolc-select uninstall 0.6.5
  gsolc-select uninstall 0.7.2 0.4.1
  gsolc-select install all
`,
	Args: cobra.MinimumNArgs(1),
	RunE: uninstallCompilers,
}

func uninstallCompilers(cmd *cobra.Command, args []string) error {
	var installedVersions = ver.GetInstalled()
	var versionsToUninstall []string
	for _, version := range args {

		if version == "all" {
			for key, _ := range installedVersions {
				versionsToUninstall = append(versionsToUninstall, key)
			}

			break
		}

		match := config.ValidSemVer.MatchString(version)
		if !match {
			fmt.Printf("Invalid version '%s'.\n", version)
			continue
		}

		if installedVersions[version] == "" {
			fmt.Printf("'%s' is not installed. Run `gsolc-select versions`.\n", version)
			continue
		}

		versionsToUninstall = append(versionsToUninstall, version)
	}

	if len(versionsToUninstall) == 0 {
		return nil
	}

	fmt.Printf("Uninstalling %s...\n", versionsToUninstall)
	uninstalled, notUninstalled, err := uninstaller.UninstallSolcs(versionsToUninstall)
	if err != nil {
		return err
	}

	for _, version := range notUninstalled {
		fmt.Printf("Failed to uninstall version: %s.\n", version)
	}

	for _, version := range uninstalled {
		fmt.Printf("Version %s uninstalled.\n", version)
	}

	return nil
}

func init() {
	RegisterCmd(rootCmd, uninstallCmd)
}
