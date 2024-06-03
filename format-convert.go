package main

import (
	"os/exec"
	"syscall"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func (a *App) Checkffmpeg() bool {
	return Checkffmpeg()
}

// 检查系统中是否安装 ffmpeg
// （临时方案）
func Checkffmpeg() bool {
	cmd := exec.Command("ffmpeg", "-version")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	err := cmd.Run()

	if err != nil {
		// 未安装 ffmpeg
		return false
	} else {
		// 已安装 ffmpeg
		return true
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
