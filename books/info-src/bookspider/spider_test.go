//
package bookspider

import (
	"fmt"
	"strings"
	"testing"

	log "github.com/wothing/log"

	"github.com/hu17889/go_spider/core/common/request"
	"github.com/hu17889/go_spider/core/spider"
)

func TestSpiderDangdangList(t *testing.T) {
	isbn := "9787513914536"
	sp := spider.NewSpider(NewDangDangListProcesser(), "spiderDangDangList")
	baseURL := "http://search.dangdang.com/?key=ISBN&act=input&category_path=01.00.00.00.00.00&type=01.00.00.00.00.00"
	url := strings.Replace(baseURL, "ISBN", isbn, -1)
	req := request.NewRequest(url, "html", "", "GET", "", nil, nil, nil, nil)

	pageItems := sp.GetByRequest(req)
	//pageItems := sp.Get("http://baike.baidu.com/view/1628025.htm?fromtitle=http&fromid=243074&type=syn", "html")

	//没爬到数据
	if pageItems == nil || len(pageItems.GetAll()) <= 0 {
		log.Debug("no matches found!")
		return
	}
	for name, value := range pageItems.GetAll() {
		log.Debug(name + "\t:\t" + value)
	}
}

func TestSpiderDangdangDetail(t *testing.T) {
	sp := spider.NewSpider(NewDangDangDetailProcesser(), "spiderDangDangDetail")
	req := request.NewRequest("http://product.dangdang.com/24170700.html", "html", "", "GET", "", nil, nil, nil, nil)

	pageItems := sp.GetByRequest(req)

	url := pageItems.GetRequest().GetUrl()
	log.Debug("-----------------------------------spider.Get---------------------------------")
	log.Debug("url\t:\t" + url)
	for name, value := range pageItems.GetAll() {
		log.Debug(name + "\t:\t" + value)
	}
}

func TestSpiderAmazonList(t *testing.T) {

	sp := spider.NewSpider(NewAmazonListProcesser(), "spiderAmazonList")
	req := request.NewRequest("https://www.amazon.cn/s/ref=nb_sb_noss?__mk_zh_CN=%E4%BA%9A%E9%A9%AC%E9%80%8A%E7%BD%91%E7%AB%99&url=search-alias%3Dstripbooks&field-keywords=9787508672069", "html", "", "GET", "", nil, nil, nil, nil)

	pageItems := sp.GetByRequest(req)
	pageItems.GetItem("")
	for name, value := range pageItems.GetAll() {
		log.Debug(name + "\t:\t" + value)
	}
}

func TestSpiderBookUUList(t *testing.T) {
	isbn := "9787559602404"
	sp := spider.NewSpider(NewBookUUListProcesser(), "BookUUlist")
	baseUrl := "http://search.bookuu.com/AdvanceSearch.php?isbn=ISBN&sm=&zz=&cbs=&dj_s=&dj_e=&bkj_s=&bkj_e=&layer2=&zk=0&cbrq_n=2017&cbrq_y=&cbrq_n1=2017&cbrq_y1=&sjsj=0&orderby=&layer1=1"
	url := strings.Replace(baseUrl, "ISBN", isbn, -1)
	req := request.NewRequest(url, "html", "", "GET", "", nil, nil, nil, nil)

	pageItems := sp.GetByRequest(req)
	for name, value := range pageItems.GetAll() {
		log.Debug(name + "\t:\t" + value)
	}
}

func TestGetBookInfo(t *testing.T) {
	book, _ := GetBookInfoBySpider("9780596001193")

	println("-----------------------------------OOOOOOM---------------------------------")
	fmt.Printf("%#v", book)
	log.Debug("-----------------------------------OOOOOOM---------------------------------")

}
