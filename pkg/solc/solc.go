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

package solc

import (
	"fmt"
	"github.com/fabelx/go-solc-select/pkg/config"
	ver "github.com/fabelx/go-solc-select/pkg/versions"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

// Execute the entrypoint called by main.go
func Execute() {
	args := os.Args[1:]
	var currentVersion, err = ver.GetCurrent()
	if err != nil {
		log.Fatal(err)
	}

	name := fmt.Sprintf("solc-%s", currentVersion)
	filePath := filepath.Join(config.SolcArtifacts, name, name)
	cmd := exec.Command(filePath, args...)
	cmd.Stdin = os.Stdin
	out, err := cmd.CombinedOutput()

	if err == nil {
		fmt.Print(string(out))
	} else if werr, ok := err.(*exec.ExitError); ok {
		if s := werr.Error(); s != "0" {
			fmt.Print(string(out))
		}

	} else {
		log.Fatal(err)
	}
}
