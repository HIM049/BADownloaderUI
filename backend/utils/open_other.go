//go:build !windows

package utils

import (
	"os/exec"
	"runtime"
)

// OpenFolder 用系统默认的文件管理器打开文件夹
func OpenFolder(path string) error {
	name := "xdg-open"
	if runtime.GOOS == "darwin" {
		name = "open"
	}
	return exec.Command(name, path).Start()
}
