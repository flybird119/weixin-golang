package wanxiang

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"

	myhttp "github.com/goushuyun/weixin-golang/misc/http"
	"github.com/goushuyun/weixin-golang/pb"
	"github.com/wangjohn/isbnconversion"
	"github.com/wothing/log"
)

const (
	jd_key = "25880c77b9cf5bdfdd4460d064e8011d"
)

var digitsRegexp = regexp.MustCompile(`(\d+\.?\d*)`)

func GetBookInfo(isbn string) (*pb.Book, error) {
	wanxiangResp := &wanxiang{}
	url := fmt.Sprintf("https://way.jd.com/jisuapi/isbn?isbn=%s&appkey=%s", isbn, jd_key)

	log.Debug(url)

	// to request data
	err := myhttp.GETWithUnmarshal(url, wanxiangResp)
	if err != nil {
		log.Error(err)
		log.Debugf("the book 【%s】 is not found", isbn)
		return nil, errors.New("not_found")
	}

	// 向 pb 对象中填充数据
	book, wanxiangBook := &pb.Book{InfoSrc: "wanxiang"}, wanxiangResp.Result.Result

	book.Title = wanxiangBook.Title
	book.Image = wanxiangBook.Pic
	book.Author = wanxiangBook.Author
	book.Publisher = wanxiangBook.Publisher
	book.Pubdate = wanxiangBook.Pubdate
	book.Subtitle = wanxiangBook.Edition
	book.Summary = wanxiangBook.Summary

	if wanxiangBook.Isbn == "" && wanxiangBook.Isbn10 != "" {
		isbn13, err := isbnconversion.ISBN10to13(wanxiangBook.Isbn10)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		book.Isbn = isbn13
	} else {
		book.Isbn = wanxiangBook.Isbn
	}

	if wanxiangBook.Price != "" {
		priceStr := digitsRegexp.FindString(wanxiangBook.Price)
		if priceInt, err := strconv.ParseFloat(priceStr, 64); err == nil {
			book.Price = int64(priceInt * 100)
		} else {
			log.Error(err)
			return nil, err
		}
	}

	log.JSON(book)
	return book, nil
}
