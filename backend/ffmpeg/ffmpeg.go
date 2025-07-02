package ffmpeg

import "os/exec"

func CheckExists() bool {
	_, err := exec.LookPath("ffmpeg")
	if err != nil {
		return false
	}
	return true
}
