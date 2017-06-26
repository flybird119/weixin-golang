package bookspider

import (
	"fmt"
	"regexp"
	"strings"

	log "github.com/wothing/log"

	"github.com/hu17889/go_spider/core/common/page"
	"github.com/hu17889/go_spider/core/common/request"
	"github.com/hu17889/go_spider/core/spider"
)

type AmazonListProcesser struct {
}

func NewAmazonListProcesser() *AmazonListProcesser {
	return &AmazonListProcesser{}
}

// Parse html dom here and record the parse result that we want to crawl.
// Package goquery (http://godoc.org/github.com/PuerkitoBio/goquery) is used to parse html.
func (s *AmazonListProcesser) Process(p *page.Page) {
	if !p.IsSucc() {
		log.Debug(p.Errormsg())
		return
	}

	query := p.GetHtmlParser()

	selection := query.Find(".s-access-detail-page")
	log.Debug(selection.Size())
	url, _ := selection.Attr("href")
	url = strings.Trim(url, " \t\n")
	sp := spider.NewSpider(NewAmazonDetailProcesser(), "spiderAmazonList")

	ip := getProxyIp()
	req := request.NewRequest(url, "html", "", "GET", "", nil, nil, nil, nil)
	if ip != "" {
		req.AddProxyHost(ip)
	}
	pageItems := sp.GetByRequest(req)
	//pageItems := sp.Get("http://baike.baidu.com/view/1628025.htm?fromtitle=http&fromid=243074&type=syn", "html")
	if pageItems == nil || pageItems.GetAll() == nil {
		return
	}
	log.Debug("-----------------------------------spider.Get---------------------------------")
	log.Debug("url\t:\t" + url)
	for name, value := range pageItems.GetAll() {
		p.AddField(name, value)
	}

}

func (s *AmazonListProcesser) Finish() {
	fmt.Printf("TODO:before end spider \r\n")
}

type AmazonDetailProcesser struct {
}

func NewAmazonDetailProcesser() *AmazonDetailProcesser {
	return &AmazonDetailProcesser{}
}

// Parse html dom here and record the parse result that we want to crawl.
// Package goquery (http://godoc.org/github.com/PuerkitoBio/goquery) is used to parse html.
func (s *AmazonDetailProcesser) Process(p *page.Page) {
	if !p.IsSucc() {
		log.Debug(p.Errormsg())
		return
	}

	query := p.GetHtmlParser()
	//获取图书名称 简介
	title := query.Find("#productTitle").Text()
	title = strings.Trim(title, " \t\n")

	//获取图书作者，出版社 ，出版时间
	var author, publisher, pubdate string
	author = query.Find("#byline .author").Text()
	publisher = query.Find("#detail_bullets_id .content li:nth-child(1)").Text()
	pubdate = query.Find("#title span").Text()

	author = strings.Trim(author, " \t\n")
	publisher = strings.Trim(publisher, " \t\n")
	author = strings.Replace(author, " ", "", -1)
	author = strings.Replace(author, "\n", "", -1)
	author = strings.Replace(author, "\t", "", -1)
	publisher = strings.Replace(publisher, " ", "", -1)
	publisher = strings.Replace(publisher, "\n", "", -1)
	publisher = strings.Replace(publisher, "\t", "", -1)
	pubdate = strings.Trim(pubdate, " \t\n")
	print(pubdate)
	reg := regexp.MustCompile("\\d{4}[\\p{Han}]{1}\\d+[\\p{Han}]{1}")
	pubdate = reg.FindString(pubdate)
	reg = regexp.MustCompile("第\\d+[版|次]")
	edition := reg.FindString(publisher)
	reg = regexp.MustCompile("出版社:[\\p{Han} ]+;")
	publisher = reg.FindString(publisher)
	publisher = strings.Replace(publisher, "出版社:", "", -1)
	publisher = strings.Replace(publisher, ";", "", -1)
	//获取图书价格
	price := query.Find("#buyBoxInner .a-color-secondary").Text()
	price = strings.Trim(price, " \t\n")

	reg = regexp.MustCompile("\\d+.\\d{2}")
	price = reg.FindString(price)
	//获取isbn
	isbnStr := query.Find("#detail_bullets_id .content li:nth-child(7)").Text()
	reg = regexp.MustCompile("(\\d[- ]*){12}[\\d]")
	isbn := reg.FindString(isbnStr)
	isbn = strings.Replace(isbn, "-", "", -1)
	isbn = strings.Replace(isbn, " ", "", -1)

	//获取图片url
	url, _ := query.Find("#img-canvas img").Attr("data-a-dynamic-image")
	urls := strings.Split(url, ":[")

	if len(urls) > 0 {
		url = urls[0]
		log.Debug(url)
		url = strings.Trim(url, " \t\n")
		log.Debug(url)
		reg = regexp.MustCompile("https://.*\"")
		url = reg.FindString(url)
		url = strings.Replace(url, "\"", "", -1)
		url = strings.Replace(url, "{", "", -1)
	} else {
		url = ""
	}

	p.AddField("title", title)
	p.AddField("remark", "")
	p.AddField("author", author)
	p.AddField("publisher", publisher)
	p.AddField("pubdate", pubdate)
	p.AddField("price", price)
	p.AddField("isbn", isbn)
	p.AddField("image_url", url)
	p.AddField("edition", edition)
}

func (s *AmazonDetailProcesser) Finish() {

}
