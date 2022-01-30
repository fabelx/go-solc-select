package cli

import (
	"fmt"
	"github.com/fabelx/go-solc-select/pkg/config"
	"github.com/fabelx/go-solc-select/pkg/installer"
	"github.com/fabelx/go-solc-select/pkg/versions"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install available solc versions",
	Long: `Gsolc-select

Installs specific versions of the solc compiler.
You can specify multiple versions separated by spaces or 'all', which will install all available versions of the compiler.
`,
	Example: "gsolc-select install 0.8.1",
	Args:    cobra.MinimumNArgs(1),
	Run:     installCompilers, // todo: Add validation for args
}

func installCompilers(cmd *cobra.Command, args []string) {
	var availableVersions, _ = versions.GetAvailable() // todo: Add logs
	// todo: Using goroutines
	// todo: External args validation

	for _, version := range args {
		match := config.ValidSemVer.MatchString(version)
		if !match {
			fmt.Printf("Invalid version '%s'.\n", version)
			continue
		}

		if availableVersions[version] == "" {
			fmt.Printf("'%s' is not avaliable. Run `gsolc-select versions installable`.\n", version)
			continue
		}

		fmt.Printf("Installing %s ...\n", version)
		err := installer.InstallSolc(version)
		if err != nil {
			fmt.Printf("Failed to install version %s.\n Error: %v\n", version, err)
			continue
		}

		fmt.Printf("Version %s installed.\n", version)
	}
}

func init() {
	RegisterCmd(rootCmd, installCmd)
}
