package main

import (
	"os/exec"
	"runtime"
	"syscall"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func (a *App) Checkffmpeg() bool {
	return Checkffmpeg()
}

// 检查系统中是否安装 ffmpeg
// （临时方案）
func Checkffmpeg() bool {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("where", "ffmpeg")
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		_, err := cmd.Output()
		return err == nil
	case "darwin": // macOS
		cmd := exec.Command("which", "ffmpeg")
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		_, err := cmd.Output()
		return err == nil
	default:
		return false
	}
}

// 调用 ffmpeg 转码
func ConventFile(inputFile, outputFile string) error {
	stream := ffmpeg.Input(inputFile).Output(outputFile, ffmpeg.KwArgs{"qscale": "0"})
	cmd := stream.Compile()
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	err := cmd.Run()

	if err != nil {
		return err
	}
	return nil
}
