//go:build windows

package utils

import "os/exec"

// OpenFolder 在资源管理器中打开文件夹
func OpenFolder(path string) error {
	return exec.Command("explorer", path).Start()
}
