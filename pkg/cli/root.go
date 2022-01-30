package cli

import (
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "gsolc-select",
	Short: "Allows users to installer and quickly switch between Solidity compiler versions.",
	Long:  `Allows users to installer and quickly switch between Solidity compiler versions`, //  todo: Write long description
}

// RegisterCmd Registers a new command under the root command
func RegisterCmd(rootCommand *cobra.Command, command *cobra.Command) {
	rootCommand.AddCommand(command)
}

// Execute the entrypoint called by main.go
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
