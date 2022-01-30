package switcher

import (
	"github.com/fabelx/go-solc-select/internal/errors"
	"github.com/fabelx/go-solc-select/pkg/config"
	"github.com/fabelx/go-solc-select/pkg/versions"
	"os"
)

func SwitchSolc(version string) error {
	installedVersions, err := versions.GetInstalled()
	if err != nil {
		return err
	}

	if installedVersions[version] == "" {
		return &errors.NotInstalledError{Version: version}
	}

	file, err := os.OpenFile(config.CurrentVersionFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}

	defer file.Close()
	_, err = file.WriteString(version)
	if err != nil {
		return err
	}

	os.Setenv("SOLC_VERSION", version)

	return nil
}
