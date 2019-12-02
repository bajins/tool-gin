package utils

import (
	"fmt"
	"time"
)

/**
类型判断
*/
func TypeJudgment(f interface{}) {
	switch vv := f.(type) {
	case string:
		if f != nil {
			fmt.Println("is string ", vv)
		}
	case int:
		if f.(int) > 0 {
			fmt.Println("is int ", vv)
		}
	case int64:
		if f.(int64) > 0 {
			fmt.Println("is int ", vv)
		}
	case time.Time:
		if !f.(time.Time).IsZero() {
			fmt.Println("is time.Time ", vv)
		}
	case []interface{}:
		fmt.Println("is array:")
		for i, j := range vv {
			fmt.Println(i, j)
		}
	}
}
