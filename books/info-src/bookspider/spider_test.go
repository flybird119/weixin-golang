//
package bookspider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"testing"

	log "github.com/wothing/log"

	"github.com/hu17889/go_spider/core/common/request"
	"github.com/hu17889/go_spider/core/spider"
)

const ProxyServer = "proxy.abuyun.com:9020"

type ProxyAuth struct {
	License   string
	SecretKey string
}

func (p ProxyAuth) ProxyClient() http.Client {
	proxyURL, _ := url.Parse("http://" + p.License + ":" + p.SecretKey + "@" + ProxyServer)
	return http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)}}
}

func TestSpiderDangdangList(t *testing.T) {
	isbn := "9780596001193"
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

func TestSpiderJDList(t *testing.T) {
	isbn := "9787301091319"
	sp := spider.NewSpider(NewJDListProcesser(), "spiderJDList")
	baseURL := "https://search.jd.com/Search?keyword=ISBN&enc=utf-8&wq=ISBN&pvid=3d3aefa8a0904ef1b08547fb69f57ae7"
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
	book, _ := GetBookInfoBySpider("9787301265017")
	println("-----------------------------------OOOOOOM---------------------------------")
	fmt.Printf("%#v", book)
	log.Debug("-----------------------------------OOOOOOM---------------------------------")

}
func TestRegular(t *testing.T) {
	detailStr := "https://item.jd.com/11020022.html"
	reg := regexp.MustCompile("/\\d*\\.")
	log.Debug(reg.FindString(detailStr))

}
func TestProxyIp(t *testing.T) {
	url := "http://api.ip.data5u.com/dynamic/get.html?order=d64615fa08c3dfea28fa9c0a1fbc3791&sep=3"
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
		return
	}

	reg := regexp.MustCompile("((2[0-4]\\d|25[0-5]|[01]?\\d\\d?)\\.){3}(2[0-4]\\d|25[0-5]|[01]?\\d\\d?)")
	ip := reg.FindString(string(body))
	log.Debug(string(body))
	log.Debug(ip)

}

func TestGolangProxy(t *testing.T) {
	//
	// reg, _ := regexp.Compile("您的IP地址是：\\[.+?\\]")
	// proxy := func(_ *http.Request) (*url.URL, error) {
	// 	return url.Parse("http://114.215.87.9") //根据定义Proxy func(*Request) (*url.URL, error)这里要返回url.URL
	// }
	// transport := &http.Transport{Proxy: proxy}
	// client := &http.Client{Transport: transport}
	// resp, err := client.Get("http://www.ip138.com") //请求并获取到对象,使用代理
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// res, err := http.Get("http://www.ip138.com") //请求并获取到对象
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// dataproxy, err := ioutil.ReadAll(resp.Body) //取出主体的内容
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// data, err := ioutil.ReadAll(res.Body) //取出主体的内容
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// //fmt.Printf("%s",data) //打印
	// //log.Debugf("=====proxy:%s", string(dataproxy))
	// log.Debugf("=====---------------------------:%s", string(data))
	//
	// sproxy := reg.FindString(string(dataproxy))
	// s := reg.FindString(string(data))
	// res.Body.Close()
	// resp.Body.Close()
	// fmt.Printf("不使用代理:%s", s)     //打印
	// fmt.Printf("使用代理:%s", sproxy) //打印
	reg := regexp.MustCompile("((2[0-4]\\d|25[0-5]|[01]?\\d\\d?)\\.){3}(2[0-4]\\d|25[0-5]|[01]?\\d\\d?)")

	res, err := http.Get("http://ip.chinaz.com/") //请求并获取到对象
	if err != nil {
		log.Fatal(err)
	}
	data, err := ioutil.ReadAll(res.Body) //取出主体的内容
	if err != nil {
		log.Fatal(err)
	}
	ip := reg.FindString(string(data))
	fmt.Printf("真实ip：%s\n", ip) //打印
	ipStr := getProxyIp()
	fmt.Printf("proxy ip := %s", ipStr)
	proxy := func(_ *http.Request) (*url.URL, error) {
		return url.Parse(ipStr) //根据定义Proxy func(*Request) (*url.URL, error)这里要返回url.URL
	}
	transport := &http.Transport{Proxy: proxy}
	client := &http.Client{Transport: transport}
	resp, err := client.Get("http://ip.chinaz.com/") //请求并获取到对象,使用代理
	if err != nil {
		log.Fatal(err)
	}
	dataproxy, err := ioutil.ReadAll(resp.Body) //取出主体的内容
	if err != nil {
		log.Fatal(err)
	}
	ip = reg.FindString(string(dataproxy))
	fmt.Printf("\n代理ip：%s\n", ip) //打印

}

func TestAbuyun(t *testing.T) {
	targetURI := "http://ip.chinaz.com/"
	//targetURI := "http://www.abuyun.com/switch-ip"
	//targetURI := "http://www.abuyun.com/current-ip"
	reg := regexp.MustCompile("((2[0-4]\\d|25[0-5]|[01]?\\d\\d?)\\.){3}(2[0-4]\\d|25[0-5]|[01]?\\d\\d?)")

	// 初始化 proxy http client
	client := ProxyAuth{License: "H2YYNX817619N32D", SecretKey: "73FAB0143E36EF3D"}.ProxyClient()

	request, _ := http.NewRequest("GET", targetURI, bytes.NewBuffer([]byte(``)))

	// 切换IP (只支持 HTTP)
	request.Header.Set("Proxy-Switch-Ip", "yes")

	response, err := client.Do(request)

	if err != nil {
		panic("failed to connect: " + err.Error())
	} else {
		bodyByte, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Println("读取 Body 时出错", err)
			return
		}
		response.Body.Close()

		body := string(bodyByte)

		fmt.Println("Response Status:", response.Status)
		fmt.Println("Response Header:", response.Header)
		fmt.Println("Response Body:\n", body)

		ip := reg.FindString(string(body))
		fmt.Printf("\n代理ip：%s\n", ip) //打印
	}

}

func TestJdAnaly(t *testing.T) {
	priceUrl := "http://p.3.cn/prices/mgets?skuIds=J_12460649031"
	// reg := regexp.MustCompile("/\\d*\\.")
	// productId := reg.FindString(productUrl)
	// productId = strings.Replace(productId, ".", "", -1)
	// productId = strings.Replace(productId, "/", "", -1)

	// log.Debug("productId========", productId)
	// priceUrl = strings.Replace(priceUrl, "PRODUCTID", productId, -1)
	log.Debug("priceUrl========", priceUrl)
	resp, err := http.Post(priceUrl,
		"application/text/html",
		strings.NewReader("name=cjb"))
	if err != nil {
		fmt.Println(err)
	}
	var price string
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}
	//获取价格
	var param []map[string]string
	log.Debug(string(body))
	err = json.Unmarshal(body, &param)
	if err != nil {
		log.Debug(err)
		return
	} else {
		price = param[0]["m"]
		if price == "" {
			return
		}
	}

	log.Debug("==============:%s", price)
}
