package util

import (
	"log"
	"os/exec"
	"runtime"
	"strings"
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

	stdOut, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	defer func() { _ = stdOut.Close() }()
	stdErr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}
	defer func() { _ = stdErr.Close() }()

	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	go func() {
		output := make([]byte, 1024)
		for {
			n, err := stdOut.Read(output)
			if err != nil {
				break
			}
			if n > 0 {
				log.Println("[stdOut]", strings.TrimSpace(string(output[:n])))
			}
		}
	}()

	go func() {
		output := make([]byte, 1024)
		for {
			n, err := stdErr.Read(output)
			if err != nil {
				break
			}
			if n > 0 {
				log.Println("[stdErr]", strings.TrimSpace(string(output[:n])))
			}
		}
	}()

	err = cmd.Wait()
	if err != nil {
		return nil, err
	}
	// 路径分隔符
	//separator := string(os.PathSeparator)
	// 调整路径分隔符
	//path := strings.ReplaceAll("/usr/local/bin/foo", "/", separator)

	// 执行命令

	return nil, nil
}
