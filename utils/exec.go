package utils

import (
	"errors"
	"fmt"
	"os/exec"
)

// 执行python脚本
// string: 要执行的Python脚本，应该是完整的路径
// args: 	参数
func ExecutePython(script string, args ...string) (string, error) {
	if !IsFileExist(script) {
		return "", errors.New(fmt.Sprintf(script, "error：%s", "文件不存在"))
	}
	name := "python"
	// 判断是否同时装了python2.7和python3，优先使用python3
	_, err := Execute("python3", "-V")
	if err == nil {
		name = "python3"
	}
	// 把脚本和参数组合到一个字符串数组
	args = append([]string{script}, args...)
	return Execute(name, args...)
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
