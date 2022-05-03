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

package errors

import "fmt"

type NotInstalledError struct {
	Version string `json:"version"`
}

type UnknownVersionError struct {
	Version string `json:"version"`
}

type UnsupportedPlatformError struct {
	Platform string `json:"platform"`
}

type UnexpectedStatusCode struct {
	StatusCode int    `json:"status_code"`
	Url        string `json:"url"`
}

type ChecksumMismatchError struct {
	HashFunc string `json:"hash_func"`
	Platform string `json:"platform"`
}

type NoCompilerSelected struct{}

func (r *NotInstalledError) Error() string {
	return fmt.Sprintf("Version '%s' not installed. Run `solc-select install %s`.", r.Version, r.Version)
}

func (r *UnknownVersionError) Error() string {
	return fmt.Sprintf("Unknown version: '%s'.", r.Version)
}

func (r *UnsupportedPlatformError) Error() string {
	return fmt.Sprintf("'%s' platform is not currently supported.", r.Platform)
}

func (r *UnexpectedStatusCode) Error() string {
	return fmt.Sprintf("Recieved unexpected status code: '%d' from '%s' request.", r.StatusCode, r.Url)
}

func (r *ChecksumMismatchError) Error() string {
	return fmt.Sprintf("%s checksum mismatch of files for %s platform.", r.HashFunc, r.Platform)
}

func (r *NoCompilerSelected) Error() string {
	return fmt.Sprintln("No compiler version selected.")
}
