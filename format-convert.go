package main

import (
	"os/exec"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func (a *App) Checkffmpeg() bool {
	return Checkffmpeg()
}

// 检查系统中是否安装 ffmpeg
// （临时方案）
func Checkffmpeg() bool {
	cmd := exec.Command("ffmpeg", "-version")
	cmd.Stdout = nil
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
	err := ffmpeg.Input(inputFile).
		Output(outputFile, ffmpeg.KwArgs{"qscale": "0"}).
		OverWriteOutput().ErrorToStdOut().Run()
	if err != nil {
		return err
	}
	return nil
}
