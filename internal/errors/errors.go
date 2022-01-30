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

func (r *NotInstalledError) Error() string {
	return fmt.Sprintf("Version '%s' not installed. Run `gsolc-select install %s`.", r.Version, r.Version)
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
