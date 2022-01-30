package cli

import (
	"fmt"
	"github.com/fabelx/go-solc-select/pkg/installer"
	"github.com/spf13/cobra"
	"os"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "List and install available solc versions",
	Long: `Gsolc-select

Prints out all installable solc versions. 
Installs specific versions of the solc compiler.
You can specify multiple versions separated by spaces or 'all', which will install all available versions of the compiler.
`,
	Run: installCompiler, // todo: Add validation for args
}

func installCompiler(cmd *cobra.Command, args []string) {
	var availableVersions, _ = installer.GetAvailableVersions() // todo: Add logs
	if len(args) == 0 {
		for _, Version := range availableVersions {
			fmt.Println(Version.String())
		}

		os.Exit(0)
	}

	for _, version := range args {
		if !installer.Contains(availableVersions, version) {
			fmt.Printf("%s is not avaliable.\n", version)
			os.Exit(1)
		}

		fmt.Printf("Installing %s ...\n", version)
		_, err := installer.InstallCompiler(version)
		if err != nil {
			fmt.Printf("Failed to install version %s.\n Error: %v\n", version, err)
			os.Exit(1)
		}
		fmt.Printf("Version %s installed.", version)
	}
}

func init() {
	RegisterCmd(rootCmd, installCmd)
}
