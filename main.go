package main

import (
	"embed"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	"network-detector/config"
	"network-detector/detector"
	"network-detector/libs/logger"
)

type EmptyStruct struct {
}

func NewEmptyStruct() *EmptyStruct {
	return &EmptyStruct{}
}

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// windows设置开机启动
	empty := NewEmptyStruct()
	otherInstanceExists, err := empty.PreHandle()
	if err != nil {
		logger.Error("Empty PreHandle error", err)
	}

	// windows第一次手动点探测后会写入参数到本地，第二次重启后默认隐藏窗口，后台运行，做到对用户无感知
	exists, err := checkArgsFile()
	if err != nil {
		fmt.Println("checkArgsFile:", err)
	}
	var hiddenWin = false

	if otherInstanceExists {
		hiddenWin = false
	} else if exists {
		hiddenWin = true
	}

	// Create application with options
	err = wails.Run(&options.App{
		Title:  "网络探测工具",
		Width:  1024,
		Height: 620,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup:   app.startup,
		StartHidden: hiddenWin,
		// close时只隐藏，后台还运行
		HideWindowOnClose: true,
		Bind: []interface{}{
			app,
			detector.NewDNS(),
			detector.NewPing(),
			detector.NewMTR(),
			detector.NewTCPConnection(),
		},
	})

	if err != nil {
		logger.Error("Error:", err.Error())
	}
}

func checkArgsFile() (bool, error) {
	filePath := config.WindowsDetectArgsCachePath
	// 检查文件是否存在
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	// 读取文件内容
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return false, err
	}

	lines := strings.Split(string(content), "\n")
	if len(lines) != 3 {
		return false, nil
	}

	// 判断是否有值
	application := strings.TrimSpace(lines[0])
	host := strings.TrimSpace(lines[1])
	stationName := strings.TrimSpace(lines[2])

	if application == "" || host == "" || stationName == "" {
		return false, nil
	}

	return true, nil
}
