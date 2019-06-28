package pyutils

import (
	"key-go/utils"
	"testing"
)

//test测试
func TestCmdPython(t *testing.T) {
	//result, err := utils.ExecutePython("xshell_key.py", "Xshell Plus", "6")
	//result, err := utils.ExecutePython("moba_xterm_Keygen.py",  utils.OsPath(),"11.1")
	result, err := utils.ExecutePython("reg_workshop_keygen.py", "10")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("转换成功", result)
}
