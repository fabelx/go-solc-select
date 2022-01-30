package cli

import (
	"fmt"
	"github.com/fabelx/go-solc-select/pkg/versions"
	"github.com/spf13/cobra"
)

var (
	windows bool
	linux   bool
	mac     bool
)

var versionCmd = &cobra.Command{
	Use:   "versions",
	Short: "Installed solc versions",
	Long: `Gsolc-select

Prints out all installed solc versions and exit.
`,
	Args: cobra.NoArgs,
	Run:  getVersions,
}

func getVersions(cmd *cobra.Command, args []string) {
	var installedVersions, err = versions.GetInstalled()
	if err != nil {
		fmt.Println(err)
	}

	// todo: Add sort
	for key, _ := range installedVersions {
		fmt.Printf("%s\n", key)
	}
}

func init() {
	RegisterCmd(rootCmd, versionCmd)
}
