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
)

var (
	windows bool
	linux   bool
	mac     bool
)

var installableCmd = &cobra.Command{
	Use:   "installable",
	Short: "Installable solc versions",
	Long: `gsolc-select

Prints out installable solc versions and exit.
`,
	Example: `  gsolc-select versions installable -l`,
	Args:    cobra.NoArgs,
	RunE:    getInstallableVersions,
}

func getInstallableVersions(cmd *cobra.Command, args []string) error {

	installableVersions := make(map[string]string)
	var err error
	if windows {
		platform := ver.WindowsPlatform{Name: config.WindowsAmd64}
		installableVersions, err = platform.GetAvailableVersions()
	} else if linux {
		platform := ver.LinuxPlatform{Name: config.LinuxAmd64}
		installableVersions, err = platform.GetAvailableVersions()
	} else if mac {
		platform := ver.MacPlatform{Name: config.MacosxAmd64}
		installableVersions, err = platform.GetAvailableVersions()
	} else {
		installableVersions, err = ver.GetAvailable()
	}

	if err != nil {
		return err
	}

	versions := ver.SortVersions(installableVersions)
	log.Warn(versions)

	return nil
}

func init() {
	installableCmd.Flags().BoolVarP(&windows, "windows", "w", false, "indicate if you want to get installable solc versions for windows OS")
	installableCmd.Flags().BoolVarP(&linux, "linux", "l", false, "indicate if you want to get installable solc versions for linux OS")
	installableCmd.Flags().BoolVarP(&mac, "mac", "m", false, "indicate if you want to get installable solc versions for mac OS")
	installableCmd.MarkFlagsMutuallyExclusive("windows", "linux", "mac")
	RegisterCmd(versionsCmd, installableCmd)
}
