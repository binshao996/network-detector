//go:build windows
// +build windows

package detector

import (
	"fmt"
	"os/exec"
	"syscall"
)

// 在 Windows 上通过 PowerShell 命令创建一个防火墙规则，允许TCP443端口通过防火墙
// open firewall, check https://github.com/microsoft/ethr
func (*TCPConnection) PreHandle() error {
	// 创建允许 TCP 443 端口通过的防火墙规则
	tcp443Cmd := exec.Command(
		"powershell",
		"-nologo",
		"-noprofile",
		"-noninteractive",
		"-c",
		`New-NetFirewallRule -DisplayName "TCP_443_Allow" -Direction Inbound -Protocol TCP -Action Allow -Profile Any -LocalPort 443`,
	)

	tcp443Cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    true,
		CreationFlags: 0x08000000,
	}

	if err := tcp443Cmd.Start(); err != nil {
		return fmt.Errorf("打开TCP443防火墙失败，因为 '%w', 请联系研发协助", err)
	}

	return nil
}
