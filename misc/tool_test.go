package misc

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/tealeg/xlsx"
)

func TestGenCheckCode(t *testing.T) {
	code := GenCheckCode(4, KC_RAND_KIND_NUM)
	fmt.Println("====>", code)
}

// 获取0-n之间的所有偶数
func even(a int) (array []int) {
	for i := 0; i < a; i++ {
		if i&1 == 0 { // 位操作符&与C语言中使用方式一样
			array = append(array, i)
		}
	}
	return array
}

// 互换两个变量的值
// 不需要使用第三个变量做中间变量
func swap(a, b int) (int, int) {
	a ^= b // 异或等于运算
	b ^= a
	a ^= b
	return a, b
}

// 左移、右移运算
func shifting(a int) int {
	a = a << 1
	a = a >> 1
	return a
}

// 变换符号
func nagation(a int) int {
	// 注意: C语言中是 ~a+1这种方式
	return ^a + 1 // Go语言取反方式和C语言不同，Go语言不支持~符号。
}

func TestBinary(t *testing.T) {
	fmt.Printf("even: %v\n", even(100))
	a, b := swap(100, 200)
	fmt.Printf("swap: %d\t%d\n", a, b)
	fmt.Printf("shifting: %d\n", shifting(100))
	fmt.Printf("nagation: %d\n", nagation(100))
	fmt.Printf("shifting:%d\n", (3 << 1))
	fmt.Printf("shifting:%d\n", ((5 << 1) & 1))
}

func TestNumFormat(t *testing.T) {
	price := 051
	discount := 0.02
	totalFee := float64(price) * discount
	fmt.Println(totalFee)
	totalPriceStr := fmt.Sprintf("%0.0f", totalFee)
	fmt.Println(totalPriceStr)

}
func TestNumFloat(t *testing.T) {
	price := 51
	discountStr := fmt.Sprintf("%.3f", float64(2)/100)
	fmt.Println(discountStr)
	discount := float64(2) / 1000
	fmt.Println(discount)
	totalFee := float64(price) * discount
	fmt.Println(totalFee)
	totalPriceStr := fmt.Sprintf("%0.0f", totalFee)
	fmt.Println(totalPriceStr)

}

func TestSubString(t *testing.T) {
	cardNo := "62284819919020398000"
	str := SubString(cardNo, len(cardNo)-4, 4)
	fmt.Print(str + "\n")
	fmt.Print(len(cardNo))
}

func TestPointerFunc(t *testing.T) {
	var p *int

	test(&p)
	fmt.Println(*p)
}

func test(p **int) {
	x := 100
	*p = &x
	fmt.Println(**p)
}

func TestDownloadAndAnaly(t *testing.T) {
	res, _ := http.Get("http://image.goushuyun.cn/Exceltest.xls")
	file, _ := os.Create("hello.xls")
	io.Copy(file, res.Body)
	xlFile, err := xlsx.OpenFile("hello.xlsx")
	if err != nil {
		fmt.Printf("err :%+v", err)
	}
	var i int
	for _, sheet := range xlFile.Sheets {
		for _, row := range sheet.Rows {
			if i == 0 {
				i++
				continue
			}
			value, _ := row.Cells[1].String()
			fmt.Printf("%s\n", value)
			if value == "" {
				break
			}
			i = i + 1
			fmt.Printf("%d\n", (i))

			// for _, cell := range row.Cells {
			// 	text, _ := cell.String()
			// 	fmt.Printf("%s\n", text)
			// }
		}
	}
	//os.Remove("hello.xls")
}
func TestUrlSubString(t *testing.T) {
	uri := "http://image.goushuyun.cn/Exceltest.xls"
	splitStringArray := strings.Split(uri, "/")
	fmt.Println(splitStringArray)
	fmt.Println(splitStringArray[len(splitStringArray)-1])

	reg := regexp.MustCompile("\\.xlsx$")
	edition := reg.FindString(uri)
	fmt.Println(edition)
	fmt.Println(edition == "")

}
