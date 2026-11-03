package ffmpeg

// ConvertToMP3 将音频转码为 MP3
func ConvertToMP3(input, output string) error {
	input, output, err := absPaths(input, output)
	if err != nil {
		return err
	}
	return run("-y", "-i", input, "-c:a", "libmp3lame", "-q:a", "0", "-f", "mp3", output)
}

// ConvertToFlac 将音频转码为 16 位 FLAC
func ConvertToFlac(input, output string) error {
	input, output, err := absPaths(input, output)
	if err != nil {
		return err
	}
	return run("-y", "-i", input, "-c:a", "flac", "-sample_fmt", "s16", "-f", "flac", output)
}
