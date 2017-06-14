package bookspider

import (
	"strconv"
	"strings"

	log "github.com/wothing/log"

	"github.com/goushuyun/weixin-golang/pb"
	"github.com/hu17889/go_spider/core/common/page_items"
	"github.com/hu17889/go_spider/core/common/request"
	"github.com/hu17889/go_spider/core/spider"
)

// "github.com/goushuyun/weixin-golang/misc/bookspider"

//通过爬虫获取图书信息
func GetBookInfoBySpider(isbn string) (book *pb.Book, err error) {
	book = &pb.Book{InfoSrc: "dangdang"}

	//首先从当当上获取图书信息
	sp := spider.NewSpider(NewDangDangListProcesser(), "spiderDangDangList")
	baseURL := "http://search.dangdang.com/?key=ISBN&ddsale=1"
	url := strings.Replace(baseURL, "ISBN", isbn, -1)
	req := request.NewRequest(url, "html", "", "GET", "", nil, nil, nil, nil)

	pageItems := sp.GetByRequest(req)
	//没爬到数据
	if pageItems == nil || len(pageItems.GetAll()) <= 0 {
		log.Debug("no matches found!")
	} else {
		structData(pageItems, book)
		if book.Isbn != "" && isbn == book.Isbn {
			//如果获取到数据，返回
			log.Debugf("%+v", book)
			return
		}

	}
	//如果当当图书信息为空 从bookUU上获取数据
	book.InfoSrc = "bookUU"
	sp = spider.NewSpider(NewBookUUListProcesser(), "BookUUlist")
	baseURL = "http://search.bookuu.com/AdvanceSearch.php?isbn=ISBN&sm=&zz=&cbs=&dj_s=&dj_e=&bkj_s=&bkj_e=&layer2=&zk=0&cbrq_n=2017&cbrq_y=&cbrq_n1=2017&cbrq_y1=&sjsj=0&orderby=&layer1=1"
	url = strings.Replace(baseURL, "ISBN", isbn, -1)
	req = request.NewRequest(url, "html", "", "GET", "", nil, nil, nil, nil)

	pageItems = sp.GetByRequest(req)
	//没爬到数据
	if pageItems == nil || len(pageItems.GetAll()) <= 0 {
		log.Debug("no matches found!")
	} else {
		structData(pageItems, book)
		if book.Isbn != "" && isbn == book.Isbn {
			//如果获取到数据，返回
			log.Debugf("%+v", book)
			return
		}

	}
	// //如果当当图书信息为空 从bookUU上获取数据
	// sp := spider.NewSpider(NewBookUUListProcesser(), "BookUUlist")
	// baseURL := "http://search.bookuu.com/AdvanceSearch.php?isbn=ISBN&sm=&zz=&cbs=&dj_s=&dj_e=&bkj_s=&bkj_e=&layer2=&zk=0&cbrq_n=2017&cbrq_y=&cbrq_n1=2017&cbrq_y1=&sjsj=0&orderby=&layer1=1"
	// url := strings.Replace(baseURL, "ISBN", isbn, -1)
	// req := request.NewRequest(url, "html", "", "GET", "", nil, nil, nil, nil)
	//
	// pageItems := sp.GetByRequest(req)
	// //没爬到数据
	// if pageItems == nil || len(pageItems.GetAll()) <= 0 {
	// 	log.Debug("no matches found!")
	// } else {
	// 	structData(pageItems, book)
	// 	if book.Isbn != "" && isbn == book.Isbn {
	// 		//如果获取到数据，返回
	// 		return
	// 	}
	//
	// }
	//
	// //首先从当当上获取图书信息
	// sp = spider.NewSpider(NewDangDangListProcesser(), "spiderDangDangList")
	// baseURL = "http://search.dangdang.com/?key=ISBN&act=input&category_path=01.00.00.00.00.00&type=01.00.00.00.00.00"
	// url = strings.Replace(baseURL, "ISBN", isbn, -1)
	// req = request.NewRequest(url, "html", "", "GET", "", nil, nil, nil, nil)
	//
	// pageItems = sp.GetByRequest(req)
	// //pageItems := sp.Get("http://baike.baidu.com/view/1628025.htm?fromtitle=http&fromid=243074&type=syn", "html")
	// //没爬到数据
	// if pageItems == nil || len(pageItems.GetAll()) <= 0 {
	// 	log.Debug("no matches found!")
	// } else {
	// 	structData(pageItems, book)
	// 	if book.Isbn != "" && isbn == book.Isbn {
	// 		//如果获取到数据，返回
	// 		return
	// 	}
	//
	// }

	//如果bookUU图书信息为空，那么向amazon获取图书信息
	book.InfoSrc = "amazon"
	sp = spider.NewSpider(NewAmazonListProcesser(), "spiderAmazonList")
	baseURL = "https://www.amazon.cn/s/ref=nb_sb_noss?__mk_zh_CN=%E4%BA%9A%E9%A9%AC%E9%80%8A%E7%BD%91%E7%AB%99&url=search-alias%3Dstripbooks&field-keywords=ISBN"
	url = strings.Replace(baseURL, "ISBN", isbn, -1)
	req = request.NewRequest(url, "html", "", "GET", "", nil, nil, nil, nil)
	pageItems = sp.GetByRequest(req)
	if pageItems == nil || len(pageItems.GetAll()) <= 0 {
		log.Debug("no matches found!")
	} else {
		structData(pageItems, book)
		if book.Isbn != "" && isbn == book.Isbn {
			//如果获取到数据，返回
			log.Debugf("%+v", book)
			return
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
