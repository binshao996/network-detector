//go:build windows
// +build windows

package detector

import (
	"fmt"
	"os/exec"
	"syscall"
)

// 因为windows下调用exec.Command("powershell")会出现powershell窗口闪屏，很不友好，而且windows防火墙默认是允许ICMPv4通过的，所以这里不做这个处理
// 在 Windows 上通过 PowerShell 命令创建一个防火墙规则，允许 ICMPv4 流量通过防火墙
// open firewall, check https://github.com/microsoft/ethr
func (*MTR) PreHandle() error {
	cmd := exec.Command(
		"powershell",
		"-nologo",
		"-noprofile",
		"-noninteractive",
		"-c",
		`New-NetFirewallRule -DisplayName "ICMP_Allow_Any" -Direction Inbound -Protocol ICMPv4 -IcmpType Any -Action Allow  -Profile Any -RemotePort Any`,
	)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    true,
		CreationFlags: 0x08000000,
	}
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("打开ICMP防火墙失败，因为'%w', 请联系研发检查", err)
	}
	return nil
}
