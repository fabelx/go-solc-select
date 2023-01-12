/*
	Copyright © 2022 Vladyslav Novotnyi <daprostovseeto@gmail.com>.

	fabelx/go-solc-select is licensed under the
	GNU Affero General Public License v3.0

	Permissions of this strongest copyleft license are conditioned on making available complete source code of
	licensed works and modifications, which include larger works using a licensed work, under the same license.
	Copyright and license notices must be preserved. Contributors provide an express grant of patent rights.
	When a modified version is used to provide a service over a network, the complete source code of the
	modified version must be made available.

go-solc-select is a tool written in Golang for managing and switching between versions of the Solidity compiler.
*/

package cli

import (
	ver "github.com/fabelx/go-solc-select/pkg/versions"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var versionsCmd = &cobra.Command{
	Use:   "versions",
	Short: "Installed solc versions",
	Long: `gsolc-select

Prints out all installed solc versions and exit.
`,
	Example: `  gsolc-select versions current
  gsolc-select versions installable
  gsolc-select versions installable -w
`,
	Args: cobra.NoArgs,
	RunE: getVersions,
}

func getVersions(cmd *cobra.Command, args []string) error {
	installedVersions := ver.GetInstalled()
	versions := ver.SortVersions(installedVersions)
	for _, version := range versions {
		log.Warnf("%s", version.String())
	}

	return nil
}

func init() {
	RegisterCmd(rootCmd, versionsCmd)
}
