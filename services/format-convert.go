package services

import (
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"os/exec"
	"runtime"
)

// 检查系统中是否安装 ffmpeg
// （临时方案）
func Checkffmpeg() bool {
	switch runtime.GOOS {
	case "windows":
		return checkffmpegOnWindows()
	case "darwin":
		return checkffmpegOnMacOS()
	default:
		return false
	}
}

// windows
func checkffmpegOnWindows() bool {
	cmd := exec.Command("where", "ffmpeg")
	setHideWindow(cmd)
	_, err := cmd.Output()
	return err == nil
}

// MacOS
func checkffmpegOnMacOS() bool {
	cmd := exec.Command("which", "ffmpeg")
	_, err := cmd.Output()
	return err == nil
}

// 调用 ffmpeg 转码
func ConventFile(inputFile, outputFile string) error {
	stream := ffmpeg.Input(inputFile).Output(outputFile, ffmpeg.KwArgs{"qscale": "0"})
	cmd := stream.Compile()
	setHideWindow(cmd)
	err := cmd.Run()

	if err != nil {
		return err
	}
	return nil
}
