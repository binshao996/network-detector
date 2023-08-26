//go:build !windows
// +build !windows

package detector

func (*MTR) PreHandle() error {
	return nil
}
