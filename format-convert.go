package main

import (
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

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
