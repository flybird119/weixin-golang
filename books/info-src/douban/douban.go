package douban

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	my_http "github.com/goushuyun/weixin-golang/misc/http"
	"github.com/goushuyun/weixin-golang/pb"
	"github.com/wangjohn/isbnconversion"
	"github.com/wothing/log"
)

var digitsRegexp = regexp.MustCompile(`(\d+\.?\d*)`)

func GetBookInfo(isbn string) (*pb.Book, error) {

	url := fmt.Sprintf("https://api.douban.com/v2/book/isbn/%s", isbn)
	doubanBook := &DoubanBook{}

	err := my_http.GETWithUnmarshal(url, doubanBook)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	log.Debugf("%#v", doubanBook)

	if doubanBook.Msg != "" {
		return nil, errors.New(doubanBook.Msg)
	}
	var book = &pb.Book{}

	//put data from DoubanBook into pb.Book
	if doubanBook.Isbn13 == "" && doubanBook.Isbn10 != "" {
		isbn13, err := isbnconversion.ISBN10to13(doubanBook.Isbn10)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		doubanBook.Isbn13 = isbn13
	}

	book.Isbn = doubanBook.Isbn13
	book.Title = doubanBook.Title
	book.Publisher = doubanBook.Publisher
	book.Author = strings.Join(doubanBook.Author, " ")
	book.Subtitle = doubanBook.Subtitle
	book.Summary = doubanBook.Summary
	book.Image = doubanBook.Images.Large
	book.AuthorIntro = doubanBook.Author_intro
	book.Pubdate = doubanBook.Pubdate

	// handle book price
	if doubanBook.Price != "" {
		priceStr := digitsRegexp.FindString(doubanBook.Price)
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
