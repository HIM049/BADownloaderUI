//go:build !windows

package utils

import "os/exec"

func setHideWindow(cmd *exec.Cmd) {}
