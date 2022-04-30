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

package versions

import (
	"github.com/Masterminds/semver"
	"github.com/fabelx/go-solc-select/internal/errors"
	"github.com/fabelx/go-solc-select/pkg/config"
	"runtime"
	"sort"
)

// GetPlatform Returns a representation of the current operating system platform
func GetPlatform() (Platform, error) {
	switch osName := runtime.GOOS; osName {
	case "darwin":
		return &MacPlatform{Name: config.MacosxAmd64}, nil
	case "linux":
		return &LinuxPlatform{Name: config.LinuxAmd64}, nil
	case "windows":
		return &WindowsPlatform{Name: config.WindowsAmd64}, nil
	default:
		return nil, &errors.UnsupportedPlatformError{Platform: osName}
	}
}

// SortVersions Returns a sorted array of versions
func SortVersions(versions map[string]string) []*semver.Version {
	vs := make([]*semver.Version, 0, len(versions))
	for v := range versions {
		vs = append(vs, semver.MustParse(v))
	}

	sort.Sort(semver.Collection(vs))
	return vs
}
