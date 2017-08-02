package bookspider

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/wothing/log"

	"github.com/goushuyun/weixin-golang/pb"
	"github.com/hu17889/go_spider/core/common/page_items"
	"github.com/hu17889/go_spider/core/common/request"
	"github.com/hu17889/go_spider/core/spider"
)

// "github.com/goushuyun/weixin-golang/misc/bookspider"

//通过爬虫获取图书信息
func GetBookInfoBySpider(isbn, upload_way string) (book *pb.Book, err error) {

	book = &pb.Book{}
	isbn = strings.Replace(isbn, "-", "", -1)
	isbn = strings.Replace(isbn, " ", "", -1)
	var num int32
	if upload_way == "batch" {
		num = rand.Int31n(3)
	} else {
		num = rand.Int31n(1)
	}
	log.Debugf("===上传类型：%s========停留秒数：%d", upload_way, num)
	time.Sleep(time.Duration(num) * time.Second)
	ip := getProxyIp()

	//首先从当当上获取图书信息
	book.InfoSrc = "dangdang"
	sp := spider.NewSpider(NewDangDangListProcesser(), "spiderDangDangList")
	baseURL := "http://search.dangdang.com/?key=ISBN&ddsale=1"
	url := strings.Replace(baseURL, "ISBN", isbn, -1)
	req := request.NewRequest(url, "html", "", "GET", "", nil, nil, nil, nil)
	log.Debug(url)
	if ip != "" {
		req.AddProxyHost(ip)
	}
	log.Debugf("==========开始执行=======", num)
	pageItems := sp.GetByRequest(req)
	//没爬到数据
	if pageItems == nil || len(pageItems.GetAll()) <= 0 {
		log.Debug("no matches found!")
	} else {
		structData(pageItems, book)
		if book.Isbn != "" && isbn == book.Isbn && book.Price != 0 && book.Title != "" {
			//如果获取到数据，返回
			log.Debugf("%+v", book)
			return
		}

	}
	//从京东上获取图书信息
	book.InfoSrc = "jd"
	sp = spider.NewSpider(NewJDListProcesser(), "spiderJDList")
	baseURL = "https://search.jd.com/Search?keyword=ISBN&enc=utf-8&wq=ISBN&pvid=3d3aefa8a0904ef1b08547fb69f57ae7&wtype=1&click=3"
	url = strings.Replace(baseURL, "ISBN", isbn, -1)
	req = request.NewRequest(url, "html", "", "GET", "", nil, nil, nil, nil)

	if ip != "" {
		req.AddProxyHost(ip)
	}
	pageItems = sp.GetByRequest(req)
	//没爬到数据
	if pageItems == nil || len(pageItems.GetAll()) <= 0 {
		log.Debug("no matches found!")
	} else {
		structData(pageItems, book)
		if book.Isbn != "" && isbn == book.Isbn && book.Price != 0 && book.Title != "" {
			//如果获取到数据，返回
			log.Debugf("%+v", book)
			return
		}

	}

	//如果当当图书信息为空 从bookUU上获取数据
	log.Debug(ip)

	book.InfoSrc = "bookUU"
	sp = spider.NewSpider(NewBookUUListProcesser(), "BookUUlist")
	baseURL = "http://search.bookuu.com/AdvanceSearch.php?isbn=ISBN&sm=&zz=&cbs=&dj_s=&dj_e=&bkj_s=&bkj_e=&layer2=&zk=0&cbrq_n=2017&cbrq_y=&cbrq_n1=2017&cbrq_y1=&sjsj=0&orderby=&layer1=1"
	url = strings.Replace(baseURL, "ISBN", isbn, -1)
	///req := request.NewRequest(url, "html", "", "GET", "", nil, nil, nil, nil)

	req = request.NewRequest(url, "html", "", "GET", "", nil, nil, nil, nil)

	if ip != "" {
		req.AddProxyHost(ip)
	}

	pageItems = sp.GetByRequest(req)
	//没爬到数据
	if pageItems == nil || len(pageItems.GetAll()) <= 0 {
		log.Debug("no matches found!")
	} else {
		structData(pageItems, book)
		if book.Isbn != "" && isbn == book.Isbn && book.Price != 0 && book.Title != "" {
			//如果获取到数据，返回
			log.Debugf("%+v", book)
			return
		}

	}

	//如果bookUU图书信息为空，那么向amazon获取图书信息
	book.InfoSrc = "amazon"
	sp = spider.NewSpider(NewAmazonListProcesser(), "spiderAmazonList")
	baseURL = "https://www.amazon.cn/s/ref=nb_sb_noss?__mk_zh_CN=%E4%BA%9A%E9%A9%AC%E9%80%8A%E7%BD%91%E7%AB%99&url=search-alias%3Dstripbooks&field-keywords=ISBN"
	url = strings.Replace(baseURL, "ISBN", isbn, -1)
	req = request.NewRequest(url, "html", "", "GET", "", nil, nil, nil, nil)
	if ip != "" {
		req.AddProxyHost(ip)
	}
	pageItems = sp.GetByRequest(req)
	if pageItems == nil || len(pageItems.GetAll()) <= 0 {
		log.Debug("no matches found!")
	} else {
		structData(pageItems, book)
		if book.Isbn != "" && isbn == book.Isbn && book.Price != 0 && book.Title != "" {
			//如果获取到数据，返回
			log.Debugf("%+v", book)
			return
		} else {
			log.Debug("not found")
		}
	}
	return nil, nil
}

func structData(items *page_items.PageItems, book *pb.Book) {
	title, _ := items.GetItem("title")
	remark, _ := items.GetItem("remark")
	author, _ := items.GetItem("author")
	publisher, _ := items.GetItem("publisher")
	pubdate, _ := items.GetItem("pubdate")
	price, _ := items.GetItem("price")
	log.Debug(price)
	priceFloat, _ := strconv.ParseFloat(price, 64)
	priceFloat = priceFloat * 100
	isbn, _ := items.GetItem("isbn")
	edition, _ := items.GetItem("edition")
	image_url, _ := items.GetItem("image_url")
	book.Title = title
	book.Isbn = isbn
	book.Price = int64(priceFloat)
	book.Author = author
	book.Publisher = publisher
	book.Pubdate = pubdate
	book.Image = image_url
	book.Summary = remark
	book.Subtitle = edition
	return
}

func getProxyIp() string {
	orderNo := getOrderNo()
	url := "http://api.ip.data5u.com/dynamic/get.html?order=" + orderNo
	resp, err := http.Post(url,
		"application/text/html",
		strings.NewReader("name=cjb"))
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		log.Error(err)
		return ""
	}
	ipStr := string(body)
	reg := regexp.MustCompile("((2[0-4]\\d|25[0-5]|[01]?\\d\\d?)\\.){3}(2[0-4]\\d|25[0-5]|[01]?\\d\\d?)")
	ip := reg.FindString(string(body))

	if ip == "" {
		ipStr = ip
	}
	ipStr = strings.TrimSpace(ipStr)
	log.Debug(ipStr)
	return "http://" + ipStr
}
