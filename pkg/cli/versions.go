package cli

import (
	"fmt"
	"github.com/fabelx/go-solc-select/pkg/versions"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "versions",
	Short: "Installed solc versions",
	Long: `Gsolc-select

Prints out all installed solc versions and exit.
`,
	Run: getVersions,
}

func getVersions(cmd *cobra.Command, args []string) {
	var Versions, _ = versions.Get()
	for _, Version := range Versions {
		if Version.Current {
			fmt.Printf("%s (current)\n", Version.Spec)
			continue
		}

		fmt.Printf("%s\n", Version.Spec)
	}
}

func init() {
	RegisterCmd(rootCmd, versionCmd)
}
