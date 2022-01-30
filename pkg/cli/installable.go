package cli

import (
	"fmt"
	"github.com/fabelx/go-solc-select/pkg/config"
	"github.com/fabelx/go-solc-select/pkg/versions"
	"github.com/spf13/cobra"
)

var installableCmd = &cobra.Command{
	Use:   "installable",
	Short: "Installable solc versions",
	Long: `Gsolc-select

Prints out installable solc versions and exit.
`,
	Args: cobra.NoArgs,
	Run:  getInstallableVersions,
}

func getInstallableVersions(cmd *cobra.Command, args []string) {

	installableVersions := make(map[string]string)
	var err error
	if windows {
		platform := versions.WindowsPlatform{Name: config.WindowsAmd64}
		installableVersions, err = platform.GetAvailableVersions() // todo: ???
	} else if linux {
		platform := versions.LinuxPlatform{Name: config.LinuxAmd64}
		installableVersions, err = platform.GetAvailableVersions() // todo: ???
	} else if mac {
		platform := versions.MacPlatform{Name: config.MacosxAmd64}
		installableVersions, err = platform.GetAvailableVersions() // todo: ???
	} else {
		installableVersions, err = versions.GetAvailable()
	}
	if err != nil {
		fmt.Println(err)
	}

	// todo: Add sort
	for key, _ := range installableVersions {
		fmt.Printf("%s\n", key)
	}

}

func init() {
	installableCmd.Flags().BoolVarP(&windows, "windows", "w", false, "indicate if you want to get installable solc versions for windows OS")
	installableCmd.Flags().BoolVarP(&linux, "linux", "l", false, "indicate if you want to get installable solc versions for linux platform")
	installableCmd.Flags().BoolVarP(&mac, "mac", "m", false, "indicate if you want to get installable solc versions for mac OS")
	RegisterCmd(versionCmd, installableCmd)
}
