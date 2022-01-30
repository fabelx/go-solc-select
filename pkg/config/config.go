package config

import (
	"fmt"
	"os"
	"path/filepath"
)

var HomeDir, _ = os.UserHomeDir()
var SolcDir = filepath.Join(HomeDir, ".gsolc-select")
var SolcArtifacts = filepath.Join(SolcDir, "artifacts")
var CurrentVersionFilePath = filepath.Join(SolcDir, "global-version")

// EnvVariable The name of the environment variable by which the current version of the solc compiler is obtained
const EnvVariable = "GSOLC_VERSION"

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
var OldSolcListUrl = fmt.Sprintf("%s/list.json", OldSolcUrl)
