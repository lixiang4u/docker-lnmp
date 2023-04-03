package util

import (
	"os/exec"
	"runtime"
)

func Exec(shString string) ([]byte, error) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		// Windows 上的命令
		cmd = exec.Command("cmd", "/C", shString)
	} else {
		// 非 Windows 上的命令
		cmd = exec.Command("sh", "-c", shString)
	}

	// 路径分隔符
	//separator := string(os.PathSeparator)
	// 调整路径分隔符
	//path := strings.ReplaceAll("/usr/local/bin/foo", "/", separator)

	// 执行命令
	return cmd.Output()
}
