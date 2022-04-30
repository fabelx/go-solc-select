# About

`Go-solc-select` - is a simple program that installs `the Solidity compiler` 
and switches between them. This can be a useful tool for managing 
different versions of the Solidity compiler, checking available versions 
for a particular operating system. It is designed to be easy to install 
and use.

The tool is split into two CLI utilities:
- `gsolc-select`: manages installing and setting different `solc` compiler versions
- `solc`: wrapper around `solc` which picks the right version according to what was set via `solc-select`

The `solc` binaries are downloaded from https://binaries.soliditylang.org/ which contains
official artifacts for many historial and modern `solc` versions for Linux and macOS.

The downloaded binaries are stored in `~/.gsolc-select/artifacts/`.

# Platforms

`Go-solc-select` is designed for use on _Unix/Linux/POSIX_ systems as a command line tool.

# Installation
...

# Usage
...

# Contributions
...

# License
`Go-solc-select` is released under the GNU Affero General Public License v3.0.
See the [LICENSE](LICENSE) file for license information.
