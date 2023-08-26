package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"network-detector/config"
	"network-detector/libs/logger"
	"network-detector/timer"
	"runtime"
	"time"
)

// App struct
type App struct {
	ctx   context.Context
	timer *timer.Timer
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.timer = timer.NewTimer()
}

func (a *App) StartTimer(domain string) {
	// 如果是windows，将参数写入到缓存文件中
	if runtime.GOOS == "windows" {
		err := writeArgsToFile(domain)
		if err != nil {
			logger.Error("Error writing args to file:", err)
		}
	}

	if a.timer != nil {
		a.timer.Stop()
		// add some time to wait for the stopping of the task
		time.Sleep(100 * time.Millisecond)
		a.timer.Start(domain)
	}
}

func writeArgsToFile(domain string) error {
	args := fmt.Sprintf("%s\n", domain)
	return ioutil.WriteFile(config.WindowsDetectArgsCachePath, []byte(args), 0644)
}
