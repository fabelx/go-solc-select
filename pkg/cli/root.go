/*
	Copyright © 2022 Vladyslav Novotnyi <daprostovseeto@gmail.com>.

	fabelx/go-solc-select is licensed under the
	GNU Affero General Public License v3.0

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU Affero General Public License as
    published by the Free Software Foundation, either version 3 of the
    License, or (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU Affero General Public License for more details.
    You should have received a copy of the GNU Affero General Public License
    along with this program.  If not, see <https://www.gnu.org/licenses/>.

go-solc-select is a tool written in Golang for managing and switching between versions of the Solidity compiler.
*/

package cli

import (
	"fmt"
	"github.com/fabelx/go-solc-select/pkg/config"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "solc-select",
	Short: "Manage multiple Solidity compiler version",
	Long: `solc-select

Allows users to installer and quickly switch between Solidity compiler versions

Example of usage:
  solc-select versions current - get current solc version


`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Checks if there are folders necessary for the application to work
		// - folder with `global-version` file. Dir:<$HomeDir/.solc-select>
		// - folder with compiler files solc. Dir:<$HomeDir/.solc-select/artifacts>
		if _, err := os.Stat(config.SolcArtifacts); os.IsNotExist(err) {
			// Creates folders if they don't exist
			err = os.MkdirAll(config.SolcArtifacts, 0755)
			if err != nil {
				fmt.Println(err) // todo: Exit?
				return
			}
		}
	},
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
