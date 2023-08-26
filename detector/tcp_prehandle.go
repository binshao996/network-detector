//go:build !windows
// +build !windows

package detector

func (*TCPConnection) PreHandle() error {
	return nil
}
