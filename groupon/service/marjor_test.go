package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestAnalyMajors(t *testing.T) {
	file1, err := os.OpenFile("/Users/lixiao/Desktop/text.txt", os.O_RDWR|os.O_APPEND|os.O_CREATE, os.ModeType)
	if err != nil {
		panic(err)
	}
	defer file1.Close()
	// 使用ioutil读取文件所有内容
	b, err := ioutil.ReadAll(file1)
	if err != nil {
		panic(err)
	}
	var rbody []map[string]interface{}
	a := strings.Replace(string(b), "\r\n", "\n", -1)
	arr := strings.Split(a, "\n")
	for i := 0; i < len(arr); i++ {
		if arr[i] == "" {
			continue
		}
		resultStr := arr[i]
		arr1 := strings.Split(resultStr, " ")
		if len(arr1) != 2 {
			continue
		}
		t := make(map[string]interface{})
		no := arr1[0]
		name := arr1[1]
		t["no"] = no
		t["name"] = name
		rbody = append(rbody, t)
	}
	res, err := json.Marshal(rbody)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%#v", string(res))
	//fmt.Printf("%v", rbody)

}
func TestStringsRepeat(t *testing.T) {
	var arr []interface{} = []interface{}{"1", "2", "3"}
	str := strings.Repeat(",%s", len(arr))
	str = str[1:]
	str = fmt.Sprintf(str, arr...)
	fmt.Println(str)
}
