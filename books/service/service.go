package service

import (
	"errors"
	"strings"

	"github.com/goushuyun/weixin-golang/errs"

	"github.com/goushuyun/weixin-golang/misc"

	"golang.org/x/net/context"

	"github.com/goushuyun/weixin-golang/books/db"
	"github.com/goushuyun/weixin-golang/books/info-src/douban"
	"github.com/goushuyun/weixin-golang/pb"
	"github.com/wothing/log"
)

type BooksServer struct {
}

func (s *BooksServer) SaveBookInfo(ctx context.Context, req *pb.Book) (*pb.GetBookInfoResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "GetBookInfo", "%#v", req))

	// save book info, level plus one and return a new ID
	err := db.SaveBook(req)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}

	return &pb.GetBookInfoResp{Code: errs.Ok, Message: "ok", Data: req}, nil
}

func (s *BooksServer) GetBookInfoByISBN(ctx context.Context, req *pb.Book) (*pb.GetBookInfoResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "GetBookInfo", "%#v", req))

	// 查找本地数据库
	err := db.GetBookInfoByISBN(req)

	if err == nil {
		// 数据库中有，就直接返回
		return &pb.GetBookInfoResp{Code: errs.Ok, Message: "ok", Data: req}, nil
	}

	// 抓取图书信息
	var (
		api_usage  int64
		final_book *pb.Book
	)
	douban_book, err := douban.GetBookInfo(req.Isbn)
	if err != nil && strings.Index(err.Error(), "404") == -1 {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	api_usage++

	// requeat other API

	if api_usage > 1 {
		// 将多API 返回图书信息整合
		// final_book = integreteInfo(douban_book, douban_book)
	} else {
		final_book = douban_book
	}

	// API 调用之后，未找到该图书，return
	if final_book == nil {
		return &pb.GetBookInfoResp{Code: errs.Ok, Message: "book_not_found"}, nil
	}

	// 数据入库
	final_book.StoreId = req.StoreId
	err = db.SaveBook(final_book)
	if err != nil {
		log.Error(err)
	}

	// 返回图书信息
	return &pb.GetBookInfoResp{Code: errs.Ok, Message: "ok", Data: final_book}, nil
}
