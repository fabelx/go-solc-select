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

package config

import (
	"os"
	"path/filepath"
	"regexp"
)

// HomeDir Home directory of the current user
var HomeDir, _ = os.UserHomeDir()

// SolcDir Directory contains go-solc-select (app) files
var SolcDir = filepath.Join(HomeDir, ".gsolc-select")

// SolcArtifacts Directory contains solc compilers
var SolcArtifacts = filepath.Join(SolcDir, "artifacts")

// CurrentVersionFilePath The name of the file that contains the current version
var CurrentVersionFilePath = filepath.Join(SolcDir, "global-version")

// LinuxAmd64 The name of the operating system for generating a link to the repository with solc compilers for Linux
const LinuxAmd64 = "linux-amd64"

// MacosxAmd64 The name of the operating system for generating a link to the repository with solc compilers for Mac
const MacosxAmd64 = "macosx-amd64"

// WindowsAmd64 The name of the operating system for generating a link to the repository with solc compilers for Windows
const WindowsAmd64 = "windows-amd64"

// SoliditylangUrl Url to repository contains current and historical builds of the Solidity Compiler
const SoliditylangUrl = "https://binaries.soliditylang.org"

// OldSolcUrl The initial part of the url to the old Solidity Compiler for Linux platform
const OldSolcUrl = "https://raw.githubusercontent.com/crytic/solc/master/linux/amd64"

// OldSolcListUrl Url to list of available old Solidity Compilers for Linux platform
const OldSolcListUrl = "https://raw.githubusercontent.com/crytic/solc/new-list-json/linux/amd64/list.json"

// GoSolcSelect The go-solc-select version
const GoSolcSelect = "0.1.0"

// ValidSemVer Regular expression for version
var ValidSemVer, _ = regexp.Compile(`^[\d]+(\.[\d]+){1,2}$`)
