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
	"errors"
	"fmt"
	"github.com/fabelx/go-solc-select/pkg/config"
	"github.com/fabelx/go-solc-select/pkg/uninstaller"
	ver "github.com/fabelx/go-solc-select/pkg/versions"
	log "github.com/sirupsen/logrus"
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
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 && all {
			return errors.New("using the --all flag and specifying explicit compiler versions are prohibited")
		}

		return nil
	},
	RunE: uninstallCompilers,
}

func uninstallCompilers(cmd *cobra.Command, args []string) error {
	var installedVersions = ver.GetInstalled()
	for _, version := range args {
		match := config.ValidSemVer.MatchString(version)
		if !match {
			return fmt.Errorf("invalid version '%s'", version)
		}

		if installedVersions[version] == "" {
			return fmt.Errorf("'%s' is not installed. Run `gsolc-select versions`", version)
		}
	}

	if all {
		for key, _ := range installedVersions {
			args = append(args, key)
		}
	}

	if len(args) == 0 {
		return nil
	}

	log.Warn("Uninstalling...")
	uninstalled, notUninstalled, err := uninstaller.UninstallSolcs(args)
	if err != nil {
		return err
	}

	for _, version := range notUninstalled {
		log.Infof("Failed to uninstall version: %s.", version)
	}

	for _, version := range uninstalled {
		log.Infof("Version %s uninstalled.", version)
	}

	return nil
}

func init() {
	uninstallCmd.Flags().BoolVarP(&all, "all", "a", false, "indicate if you want to uninstall all installed solc versions")
	RegisterCmd(rootCmd, uninstallCmd)
}
