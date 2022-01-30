package versions

import (
	"github.com/fabelx/go-solc-select/internal/errors"
	"github.com/fabelx/go-solc-select/pkg/config"
	"runtime"
)

// GetPlatform Returns a representation of the current operating system platform
func GetPlatform() (Platform, error) { // todo: Extend platform support
	switch osName := runtime.GOOS; osName {
	case "darwin":
		return &MacPlatform{Name: config.MacosxAmd64}, nil
	case "linux":
		return &LinuxPlatform{Name: config.LinuxAmd64}, nil
	case "windows":
		return &WindowsPlatform{Name: config.WindowsAmd64}, nil
	default:
		return nil, &errors.UnsupportedPlatformError{Platform: osName} // todo: Unsupported platform Error
	}
}
