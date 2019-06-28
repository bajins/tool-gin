package utils

import (
	"errors"
	"fmt"
	"os/exec"
	"runtime"
)

/**
 * 执行python脚本
 *
 * @author claer www.bajins.com
 * @date 2019/6/28 14:19
 */
func ExecutePython(script string, args ...string) (result string, err error) {
	if !IsFile(script) {
		err = errors.New(fmt.Sprintf(script, "error：%s", "文件不存在"))
		return "", err
	}
	name := "python"
	if runtime.GOOS == "linux" {
		name = "python3"
	}
	// 把脚本和参数组合到一个字符串数组
	args = append([]string{script}, args...)
	out, err := exec.Command(name, args...).Output()
	// exit status 2 一般是文件没有找到
	// exit status 1 一般是命令执行错误
	if err != nil {
		return string(out), err
	}
	result = string(out)

	return result, nil
}
