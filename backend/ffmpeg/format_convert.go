package ffmpeg

import (
	"bili-audio-downloader/backend/utils"
	"errors"
	"fmt"
	"path/filepath"
)

func ConvertToMP3(input, output string) error {
	err := initFfmpeg(&input, &output)
	if err != nil {
		return err
	}

	// 转码文件
	log, err := utils.RunCommand("ffmpeg", "-i", input, "-c:a", "libmp3lame", "-q:a", "0", output)
	if err != nil {
		return errors.New(fmt.Sprintf("ffmpeg error: %v, output: %s", err, log))
	}

	return nil
}

func ConvertToFlac(input, output string) error {
	err := initFfmpeg(&input, &output)
	if err != nil {
		return err
	}

	// 转码文件
	log, err := utils.RunCommand("ffmpeg", "-i", input, "-c:a", "flac", "-sample_fmt", "s16", output)
	if err != nil {
		return errors.New(fmt.Sprintf("ffmpeg error: %v, output: %s", err, log))
	}

	return nil
}

func initFfmpeg(input, output *string) error {
	// 检查 ffmpeg
	if !CheckExists() {
		return errors.New("ffmpeg not exists")
	}

	// 绝对化路径
	inputAbs, err := filepath.Abs(*input)
	if err != nil {
		return err
	}
	outputAbs, err := filepath.Abs(*output)
	if err != nil {
		return err
	}

	*input = inputAbs
	*output = outputAbs
	return nil
}
