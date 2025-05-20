//go:build windows

package services

import (
	"os/exec"
	"syscall"
)

func setHideWindow(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
}
