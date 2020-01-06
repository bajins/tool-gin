package utils

import (
	"fmt"
	"math/rand"
	"reflect"
	"strconv"
	"time"
)

// 截取小数位数
func FloatRound(f float64, n int) float64 {
	format := "%." + strconv.Itoa(n) + "f"
	res, _ := strconv.ParseFloat(fmt.Sprintf(format, f), 64)
	return res
}

// 将任何数值转换为int64
func ToInt64(value interface{}) (d int64, err error) {
	val := reflect.ValueOf(value)
	switch value.(type) {
	case int, int8, int16, int32, int64:
		d = val.Int()
	case uint, uint8, uint16, uint32, uint64:
		d = int64(val.Uint())
	default:
		err = fmt.Errorf("ToInt64 need numeric not `%T`", value)
	}

	return
}

// 生成随机数
//
// 首先要初始化随机种子，不然每次生成都是（指每次重新开始）相同的数
// 系统每次都会先用Seed函数初始化系统资源，如果用户不提供seed参数，则默认用seed=1来初始化，这就是为什么每次都输出一样的值的原因
func RandIntn(length int) int {
	// 用一个不确定数来初始化随机种子
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(length)
}
