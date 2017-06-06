//
package misc

import (
	"fmt"
	"strings"
	"testing"

	"github.com/goushuyun/weixin-golang/misc/bookspider"
	"github.com/hu17889/go_spider/core/common/request"
	"github.com/hu17889/go_spider/core/spider"
)

func TestSpiderDangdangList(t *testing.T) {
	isbn := "9787513914536"
	sp := spider.NewSpider(bookspider.NewDangDangListProcesser(), "spiderDangDangList")
	baseURL := "http://search.dangdang.com/?key=ISBN&act=input&category_path=01.00.00.00.00.00&type=01.00.00.00.00.00"
	url := strings.Replace(baseURL, "ISBN", isbn, -1)
	req := request.NewRequest(url, "html", "", "GET", "", nil, nil, nil, nil)

	pageItems := sp.GetByRequest(req)
	//pageItems := sp.Get("http://baike.baidu.com/view/1628025.htm?fromtitle=http&fromid=243074&type=syn", "html")

	//没爬到数据
	if pageItems == nil || len(pageItems.GetAll()) <= 0 {
		println("no matches found!")
		return
	}
	for name, value := range pageItems.GetAll() {
		println(name + "\t:\t" + value)
	}
}

func TestSpiderDangdangDetail(t *testing.T) {
	sp := spider.NewSpider(bookspider.NewDangDangDetailProcesser(), "spiderDangDangDetail")
	req := request.NewRequest("http://product.dangdang.com/24170700.html", "html", "", "GET", "", nil, nil, nil, nil)

	pageItems := sp.GetByRequest(req)

	url := pageItems.GetRequest().GetUrl()
	println("-----------------------------------spider.Get---------------------------------")
	println("url\t:\t" + url)
	for name, value := range pageItems.GetAll() {
		println(name + "\t:\t" + value)
	}
}

func TestSpiderAmazonList(t *testing.T) {

	sp := spider.NewSpider(bookspider.NewAmazonListProcesser(), "spiderAmazonList")
	req := request.NewRequest("https://www.amazon.cn/s/ref=nb_sb_noss?__mk_zh_CN=%E4%BA%9A%E9%A9%AC%E9%80%8A%E7%BD%91%E7%AB%99&url=search-alias%3Dstripbooks&field-keywords=9787508672069", "html", "", "GET", "", nil, nil, nil, nil)

	pageItems := sp.GetByRequest(req)
	pageItems.GetItem("")
	for name, value := range pageItems.GetAll() {
		println(name + "\t:\t" + value)
	}
}

func TestGetBookInfo(t *testing.T) {
	book, _ := GetBookInfoBySpider("9787559602404")
	println("-----------------------------------OOOOOOM---------------------------------")
	fmt.Printf("%#v", book)
	println("-----------------------------------OOOOOOM---------------------------------")

}
