package logger

import (
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var log *logrus.Logger

func init() {
	// 创建日志文件目录
	logDir := getLogDir()
	if err := os.MkdirAll(logDir, 0755); err != nil {
		logrus.Fatal("Failed to create log directory:", err)
	}

	// 设置logrus输出为文件
	logFilePath := filepath.Join(logDir, "network-detect.log")

	// 检查日志文件是否存在，如果不存在则创建
	if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
		file, err := os.Create(logFilePath)
		if err != nil {
			log.Fatal("Failed to create log file:", err)
		}
		file.Close()
	}

	logFile := &lumberjack.Logger{
		Filename:   logFilePath,
		MaxSize:    30,   // 单位：MB
		MaxBackups: 1,    // 最多保留1个备份文件
		MaxAge:     2,    // 最多保留2天
		Compress:   true, // 是否压缩备份文件
	}

	log = logrus.New()
	log.SetOutput(logFile)
	log.SetLevel(logrus.DebugLevel)
	log.SetFormatter(&logrus.TextFormatter{
		DisableColors:   true,
		DisableSorting:  true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	if !isWindows() {
		err := os.Chmod(logFilePath, 0666)
		if err != nil {
			log.Fatal("Failed to set log file permission:", err)
		}
	}
}

// getLogDir 返回日志文件目录
func getLogDir() string {
	if isWindows() {
		// Windows 环境下的日志文件目录（获取当前执行程序的路径（不包含程序名））
		return filepath.Join("C:\\", "logs")
	}
	// 其他平台环境下的日志文件目录
	homeDir, err := os.UserHomeDir()
	if err != nil {
		logrus.Fatal("Failed to get user home directory:", err)
	}
	return filepath.Join(homeDir, "logs")
}

// isWindows 判断是否为 Windows 环境
func isWindows() bool {
	return os.PathSeparator == '\\'
}

// Info 记录信息级别的日志
func Info(args ...interface{}) {
	log.Info(args...)
}

// LogPrintf logs a formatted message using logrus.Printf
func Printf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

// Error 记录错误级别的日志
func Error(args ...interface{}) {
	log.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

// GetLogger 返回日志记录器实例
func GetLogger() *logrus.Logger {
	return log
}
