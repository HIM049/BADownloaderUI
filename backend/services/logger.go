package services

import (
	"fmt"
	"os"
)

var Logger *CustomLogger

type CustomLogger struct {
	logFile *os.File
}

// 创建并返回日志记录器
func NewCustomLogger() (*CustomLogger, error) {

	// 创建日志文件
	logFile, err := os.OpenFile("./app.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return nil, fmt.Errorf("创建日志文件失败: %v", err)
	}

	return &CustomLogger{logFile: logFile}, nil
}

func (cl *CustomLogger) Print(message string) {
	cl.logToFile(message)
}

func (cl *CustomLogger) Trace(message string) {
	cl.logToFile("TRACE: " + message)
}

func (cl *CustomLogger) Debug(message string) {
	cl.logToFile("DEBUG: " + message)
}

func (cl *CustomLogger) Info(message string) {
	cl.logToFile("INFO: " + message)
}

func (cl *CustomLogger) Warning(message string) {
	cl.logToFile("WARNING: " + message)
}

func (cl *CustomLogger) Error(message string) {
	cl.logToFile("ERROR: " + message)
}

func (cl *CustomLogger) Fatal(message string) {
	cl.logToFile("FATAL: " + message)
	os.Exit(1)
}

// logToFile 记录日志到文件
func (cl *CustomLogger) logToFile(message string) {
	fmt.Fprintln(cl.logFile, message)
}
