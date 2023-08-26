//go:build !windows
// +build !windows

package main

func (*EmptyStruct) PreHandle() (bool, error) {
	return false, nil
}
