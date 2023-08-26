//go:build windows
// +build windows

package main

import (
	"fmt"
	"network-detector/libs/logger"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"

	"golang.org/x/sys/windows/registry"
)

func (*EmptyStruct) PreHandle() (bool, error) {
	// 设置开机启动
	// 获取应用的可执行文件路径
	exePath, err := os.Executable()
	if err != nil {
		logger.Error("Failed to get executable path:", err)
		return false, err
	}

	// 注册表路径
	regPath := `SOFTWARE\Microsoft\Windows\CurrentVersion\Run`

	// 打开注册表项
	key, _, err := registry.CreateKey(registry.CURRENT_USER, regPath, registry.SET_VALUE)
	if err != nil {
		logger.Error("Failed to open registry key:", err)
		return false, err
	}
	defer key.Close()

	// 设置注册表项的值，注意这里的路径
	appName := "network-detector" // 修改为你的应用名称
	logger.Printf("exePath is exePath:%s", exePath)
	if err := key.SetStringValue(appName, exePath); err != nil {
		logger.Error("Failed to set registry value:", err)
		return false, err
	}

	logger.Info("Successfully set up admin auto start on Windows for", appName)
	otherInstanceExists := isOtherInstanceRunning()
	return otherInstanceExists, nil
}

// 如果打开程序时发现已经有实例在后台运行了，先干掉后台实例
// 这是为了在开机启动后，防止点击shorticon窗口一直不出来，因为默认窗口隐藏了
func isOtherInstanceRunning() bool {
	currentProcessID := os.Getpid()
	currentProcessName := "network-detector"

	powershellCommand := `Get-Process | Format-Table -HideTableHeaders -Property Id,Name`
	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-ExecutionPolicy", "Bypass", "-Command", powershellCommand)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    true,
		CreationFlags: 0x08000000,
	}
	output, err := cmd.Output()

	if err != nil {
		return false
	}

	exists := false
	outputStr := string(output)
	lines := strings.Split(outputStr, "\n")
	for _, line := range lines {
		if strings.Contains(line, currentProcessName) {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				processIDStr := fields[0]
				processName := fields[1]
				processID, err := strconv.Atoi(processIDStr)

				if err == nil && strings.EqualFold(processName, currentProcessName) && processID != currentProcessID {
					if err := terminateProcess(processID); err != nil {
						logger.Printf("Error terminating process %d: %v\n", processID, err)
					} else {
						exists = true
						logger.Printf("Other instance exists, Terminated process %d\n", processID)
					}
				}
			}
		}
	}

	return exists
}

func terminateProcess(processID int) error {
	powershellCommand := fmt.Sprintf(`Stop-Process -Force -Id %d`, processID)
	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-ExecutionPolicy", "Bypass", "-Command", powershellCommand)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    true,
		CreationFlags: 0x08000000,
	}

	if err := cmd.Run(); err != nil {
		logger.Printf("Error terminating process %d: %v\n", processID, err)
		return err
	} else {
		logger.Printf("Terminated process %d\n", processID)
	}
	return nil
}
