package cli

import (
	"fmt"
	"github.com/fabelx/go-solc-select/pkg/versions"
	"github.com/spf13/cobra"
)

var currentCmd = &cobra.Command{
	Use:   "current",
	Short: "Current solc versions",
	Long: `Gsolc-select

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
