package cli

import (
	"fmt"
	"github.com/fabelx/go-solc-select/pkg/installer"
	"github.com/fabelx/go-solc-select/pkg/switcher"
	"github.com/spf13/cobra"
	"os"
)

var (
	install bool
)

var useCmd = &cobra.Command{
	Use:   "use",
	Short: "Change the version of global solc compiler",
	Long: `Gsolc-select

Switch between installed versions of solc compiler. 
Using the -i / --installer flag automatically installer the required compiler version.
`,
	Args: cobra.MinimumNArgs(1),
	Run:  useCompiler, // todo: Add validation for args
}

func useCompiler(_ *cobra.Command, args []string) {
	var availableVersions, _ = installer.GetAvailableVersions()
	Version := args[0]
	if !installer.Contains(availableVersions, Version) {
		fmt.Printf("%s is not avaliable.\n", Version)
		os.Exit(1)
	}

	if install {
		fmt.Printf("Installing %s ...\n", Version)
		_, err := installer.InstallCompiler(Version)
		if err != nil {
			fmt.Printf("Failed to install version %s.\n Error: %v\n", Version, err)
			os.Exit(1)
		}
		fmt.Printf("Version %s installed.\n", Version)
	}

	_, err := switcher.SwitchCompiler(Version)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("Version %s not installed. Run `gsolc-select install %s.`\n", Version, Version)
		} else {
			fmt.Printf("Failed to switch version to %s.\nError: %v\n", Version, err)
		}

		os.Exit(1)
	}

	fmt.Printf("Switched global version to %s\n", Version)
}

func init() {
	useCmd.Flags().BoolVarP(&install, "install", "i", false, "indicate if you want to automatically installer versions that are not installed")
	RegisterCmd(rootCmd, useCmd)
}
