// Copyright 2017 chai Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

// 返回参数的类型
func TestType(t *testing.T) {
	var i int
	var s string

	if Type(i) != "int" {
		t.Error("type get fail by int")
	}

	if Type(s) != "string" {
		t.Error("type get fail by string")
	}
}

// 判断是否在数组中
func TestInArray(t *testing.T) {
	a := "key"
	aList := []string{"key", "key2", "key3"}
	aList2 := []string{"key2", "key3"}

	if InArray(a, aList) != true {
		t.Error("value is in array")
	}

	if InArray(a, aList2) != false {
		t.Error("value is not in array")
	}

	b := 2
	bList := []int{2, 3, 4, 5}
	if InArray(b, aList) != false {
		t.Error("value is not in array")
	}

	if InArray(b, bList) != true {
		t.Error("value is in array")
	}

}

// 通过scrypt生成密码
func TestNewPass(t *testing.T) {
	p, err := NewPass("123456", "123")
	if err != nil {
		t.Error(err.Error())
	}

	if len(p) != 64 {
		t.Error("password hash fail")
	}
}

func TestHttp(t *testing.T) {
	params := map[string]string{"test": "1", "t": "22"}
	var param string
	for key, value := range params {
		param += key + "=" + value + "&"
	}
	param = param[0 : len(param)-1]
	t.Error(param)
	result, err := HttpReadBodyString(http.MethodPost, "test", "", map[string]string{"test": "1", "t": "22"}, nil)
	t.Log(result, err)
	httpClient := HttpClient{
		Method:      http.MethodPost,
		UrlText:     "test",
		ContentType: ContentTypeMFD,
		Params:      nil,
		Header:      nil,
	}
	t.Log(httpClient.HttpReadBodyJsonMap())
}

func TestSchedulerIntervalsTimer(t *testing.T) {
	SchedulerIntervalsTimer(fmtp, time.Second*5)
}

func TestSchedulerFixedTimer(t *testing.T) {
	SchedulerFixedTicker(fmtp, time.Second*5)
}

func fmtp() {
	fmt.Println(TimeToString(time.Now()))
}
