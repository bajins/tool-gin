package utils

import (
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
	name := "python"
	if runtime.GOOS == "linux" {
		name = "python3"
	}
	// 把脚本和参数组合到一个字符串数组
	args = append([]string{script}, args...)
	out, err := exec.Command(name, args...).Output()
	if err != nil {
		return string(out), err
	}
	result = string(out)
	//if strings.Index(result, "success") != 0 {
	//	err = errors.New(fmt.Sprintf(script, "error：%s", result))
	//}
	return result, nil
}
