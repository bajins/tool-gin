/**
 *
 * @Description:
 * @Author: https://www.bajins.com
 * @File: reptile_test.go
 * @Version: 1.0.0
 * @Time: 2019/9/19 11:13
 * @Project: key-gin
 * @Package:
 * @Software: GoLand
 */
package reptile

import (
	"testing"
)

func TestCDP(t *testing.T) {
	SendMail("nmuqr3on@linshiyouxiang.net", "xftp")
}

func TestLinShiYouXiangList(t *testing.T) {
	_, err := DownloadNetsarang("xshell")
	t.Log(err)

}

func TestXshell(t *testing.T) {
	SendMail("", "xshell")
}
