package service

import (
	"context"
	"errors"
	"path/filepath"
	"strings"

	"google.golang.org/grpc/metadata"

	"github.com/goushuyun/weixin-golang/errs"
	"github.com/goushuyun/weixin-golang/misc"
	"github.com/pborman/uuid"

	"github.com/goushuyun/weixin-golang/books/db"
	"github.com/goushuyun/weixin-golang/books/info-src/bookspider"
	"github.com/goushuyun/weixin-golang/pb"
	"github.com/wothing/log"
)

func integreteInfo(to, from *pb.Book) *pb.Book {
	if to == nil && from != nil {
		return from
	} else if to == nil && from == nil {
		return nil
	}

	if to.Title == "" {
		to.Title = from.Title
	}

	if to.Price == 0 {
		to.Price = from.Price
	}

	if to.Author == "" {
		to.Author = from.Author
	}

	if to.Publisher == "" {
		to.Publisher = from.Publisher
	}

	if to.Image == "" {
		to.Image = from.Image
	}

	if to.Pubdate == "" {
		to.Pubdate = from.Pubdate
	}

	if to.Subtitle == "" {
		to.Subtitle = from.Subtitle
	}

	return to
}

func bookInfoIsOk(book *pb.Book) bool {
	if book == nil {
		return false
	}

	if book.Price != 0 && book.Isbn != "" && book.Title != "" {
		return true
	} else {
		return false
	}
}

func GetBookInfoByISBNWithNoContext(req *pb.Book) (*pb.Book, error) {

	// 查找本地数据库
	err := db.GetBookInfoByISBN(req)

	if err == nil {
		// 数据库中有，就直接返回
		return req, nil
	}

	// 抓取图书信息
	var (
		api_usage  int64
		final_book *pb.Book
	)

	// get from spider
	spider_book, err := bookspider.GetBookInfoBySpider(req.Isbn, req.UploadWay)
	if err != nil {
		spider_book = nil
		log.Error(err)
	}

	api_usage++

	// if !bookInfoIsOk(spider_book) {
	// 	// get from wanxiang
	// 	wanxiang_book, err1 := wanxiang.GetBookInfo(req.Isbn)
	// 	if err1 != nil {
	// 		log.Error(err)
	// 		return nil, errs.Wrap(errors.New(err1.Error()))
	// 	}
	// 	api_usage++
	// 	final_book = integreteInfo(spider_book, wanxiang_book)
	// }

	if api_usage == 1 {
		final_book = spider_book
	}

	// API 调用之后，未找到该图书，或者图书信息不完善 return
	if final_book == nil || !bookInfoIsOk(final_book) {
		return nil, nil
	}

	// 抓取图书图片，存到七牛
	final_book.StoreId = req.StoreId
	if strings.HasPrefix(final_book.Image, "http") {
		fetchImageReq := &pb.FetchImageReq{
			Zone: pb.MediaZone_Test,
			Url:  final_book.Image,
			Key:  final_book.StoreId + "/" + final_book.Isbn + filepath.Ext(final_book.Image),
		}
		mediaResp := &pb.FetchImageResp{}
		ctx := metadata.NewContext(context.Background(), metadata.Pairs("tid", uuid.New()))
		err = misc.CallSVC(ctx, "bc_mediastore", "FetchImage", fetchImageReq, mediaResp)
		if err != nil {
			log.Error(err)
			return nil, errs.Wrap(errors.New(err.Error()))
		}
		final_book.Image = fetchImageReq.Key
	}

	// 数据入库
	err = db.SaveBook(final_book)
	if err != nil {
		log.Error(err)
	}

	// 返回图书信息
	return final_book, nil
}
