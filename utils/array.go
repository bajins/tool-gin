/**
 *
 * @Description:
 * @Author: https://www.bajins.com
 * @File: array.go
 * @Version: 1.0.0
 * @Time: 2019/9/25 10:26
 * @Project: tool-gin
 * @Package:
 * @Software: GoLand
 */
package utils

import (
	"math/rand"
	"time"
)

// 随机打乱数组顺序
func Shuffle(slice []interface{}) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for len(slice) > 0 {
		n := len(slice)
		randIndex := r.Intn(n)
		slice[n-1], slice[randIndex] = slice[randIndex], slice[n-1]
		slice = slice[:n-1]
	}
}
