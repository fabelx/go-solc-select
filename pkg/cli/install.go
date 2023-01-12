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
	"context"
	"errors"
	"fmt"
	"github.com/fabelx/go-solc-select/pkg/config"
	"github.com/fabelx/go-solc-select/pkg/installer"
	ver "github.com/fabelx/go-solc-select/pkg/versions"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os/signal"
	"syscall"
)

var (
	async bool
	all   bool
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install available solc versions",
	Long: `gsolc-select

Installs specific versions of the solc compiler.
You can specify multiple versions separated by spaces or flag '--all/-a', which will install all available versions of the compiler.
`,
	Example: `  gsolc-select install 0.8.1
  gsolc-select install 0.8.1 0.4.23
  gsolc-select install --all
`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 && all {
			return errors.New("using the --all flag and specifying explicit compiler versions are prohibited")
		}

		return nil
	},
	RunE: installCompilers,
}

func installCompilers(cmd *cobra.Command, args []string) error {
	availableVersions, err := ver.GetAvailable()
	if err != nil {
		return err
	}

	installedVersions := ver.GetInstalled()
	var versions []string
	for _, version := range args {
		match := config.ValidSemVer.MatchString(version)
		if !match {
			return fmt.Errorf("invalid version '%s'", version)
		}

		if availableVersions[version] == "" {
			return fmt.Errorf("'%s' is not avaliable. Run `gsolc-select versions installable`", version)
		}

		if installedVersions[version] != "" {
			return fmt.Errorf("version '%s' is already installed. Run `gsolc-select versions`", version)
		}
	}

	if all {
		for key, _ := range availableVersions {
			versions = append(versions, key)
		}

		args = versions
	}

	if len(args) == 0 {
		return errors.New("wrong number of args, required at least one or flag `--all/-a`")
	}

	log.Warn("Installing...")
	var installed, notInstalled []string
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()
	if async {
		installed, notInstalled, err = installer.AsyncInstallSolcs(ctx, args)
	} else {
		installed, notInstalled, err = installer.InstallSolcs(ctx, args)
	}

	if err != nil {
		return err
	}

	for _, version := range notInstalled {
		log.Infof("Failed to install version %s.", version)
	}

	for _, version := range installed {
		log.Infof("Version %s installed.", version)
	}

	return nil
}

func init() {
	installCmd.Flags().BoolVarP(&async, "parallel", "p", false, "indicate if you want to install solc versions asynchronously")
	installCmd.Flags().BoolVarP(&all, "all", "a", false, "indicate if you want to install all available solc versions")
	RegisterCmd(rootCmd, installCmd)
}
