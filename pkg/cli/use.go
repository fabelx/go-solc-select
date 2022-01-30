package cli

import (
	"fmt"
	"github.com/fabelx/go-solc-select/pkg/config"
	"github.com/fabelx/go-solc-select/pkg/switcher"
	"github.com/fabelx/go-solc-select/pkg/versions"
	"github.com/spf13/cobra"
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
	Args: cobra.ExactArgs(1),
	Run:  useCompiler, // todo: Add validation for args
}

func useCompiler(cmd *cobra.Command, args []string) {
	version := args[0]
	match := config.ValidSemVer.MatchString(version)
	if !match {
		fmt.Printf("Invalid version '%s'.\n", version)
		return
	}
	var availableVersions, _ = versions.GetAvailable()
	if availableVersions[version] == "" {
		fmt.Printf("'%s' is not avaliable. Run `gsolc-select versions installable`.\n", version)
		return
	}

	if install {
		installCompilers(cmd, args)
	}

	err := switcher.SwitchSolc(version)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Switched global version to '%s'.\n", version)
}

func init() {
	useCmd.Flags().BoolVarP(&install, "install", "i", false, "indicate if you want to automatically installer versions that are not installed")
	RegisterCmd(rootCmd, useCmd)
}
