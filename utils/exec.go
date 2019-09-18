package utils

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"runtime"
)

// 执行python脚本
// string: 要执行的Python脚本，应该是完整的路径
// args: 	参数
func ExecutePython(script string, args ...string) (string, error) {
	if !IsFile(script) {
		err := errors.New(fmt.Sprintf(script, "error：%s", "文件不存在"))
		return "", err
	}
	name := "python"
	if runtime.GOOS == "linux" {
		name = "python3"
	}
	// 把脚本和参数组合到一个字符串数组
	args = append([]string{script}, args...)
	out, err := Execute(name, args...)
	if err != nil {
		if err.Error() == "exit status 1" {
			// 获取当前绝对路径
			dir, err := os.Getwd()
			if err == nil {
				p := path.Join(dir, "pyutils", "requirements.txt")
				Execute("pip", "install", "-r", p)
			}
		}
	}
	return out, err
}

// 执行dos或shell命令
// program: 程序名称
// args: 	参数
func Execute(program string, args ...string) (string, error) {
	// exit status 2 一般是文件没有找到
	// exit status 1 一般是命令执行错误
	out, err := exec.Command(program, args...).Output()
	return string(out), err
}
