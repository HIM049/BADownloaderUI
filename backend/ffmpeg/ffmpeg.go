package ffmpeg

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
)

// CheckExists 检查 PATH 中是否有 ffmpeg
func CheckExists() bool {
	_, err := exec.LookPath("ffmpeg")
	return err == nil
}

// run 以隐藏窗口的方式执行 ffmpeg，出错时附带命令输出
func run(args ...string) error {
	cmd := exec.Command("ffmpeg", args...)
	setHideWindow(cmd)

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ffmpeg error: %w, output: %s", err, out.String())
	}
	return nil
}

// absPaths 检查 ffmpeg 是否可用并返回输入输出的绝对路径
func absPaths(input, output string) (string, string, error) {
	if !CheckExists() {
		return "", "", errors.New("ffmpeg not exists")
	}
	input, err := filepath.Abs(input)
	if err != nil {
		return "", "", err
	}
	output, err = filepath.Abs(output)
	if err != nil {
		return "", "", err
	}
	return input, output, nil
}
