package service

import (
	"errors"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/goushuyun/weixin-golang/errs"

	"github.com/goushuyun/weixin-golang/misc"

	"golang.org/x/net/context"

	"github.com/goushuyun/weixin-golang/books/db"
	"github.com/goushuyun/weixin-golang/books/info-src/douban"
	"github.com/goushuyun/weixin-golang/books/info-src/wanxiang"
	"github.com/goushuyun/weixin-golang/pb"
	"github.com/wothing/log"
)

type BooksServer struct {
}

func (s *BooksServer) GetBookInfo(ctx context.Context, req *pb.Book) (*pb.Book, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "GetBookInfo", "%#v", req))

	err := db.GetBookInfo(req)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	return req, nil
}

func (s *BooksServer) ModifyBookInfo(ctx context.Context, req *pb.Book) (*pb.GetBookInfoResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "SaveBookInfo", "%#v", req))

	// get book by book_id
	old_book := &pb.Book{Id: req.Id}
	if err := db.GetBookInfo(old_book); err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}

	new_book_v := reflect.ValueOf(*req)
	new_book_t := reflect.TypeOf(*req)

	old_book_v := reflect.ValueOf(old_book).Elem()

	for i := 0; i < new_book_v.NumField(); i++ {
		if new_book_t.Field(i).Name == "Id" {

			log.Debugf("The ignore field is %s", new_book_t.Field(i).Name)
			continue //id and level is not care
		}

		// to handle other string and int64 field
		if new_book_v.Field(i).Kind() == reflect.String && new_book_v.Field(i).String() != "" {
			// handle field, which type is string and not null
			field_name := new_book_t.Field(i).Name
			field_val := new_book_v.Field(i).String()

			if old_book_v.Field(i).CanSet() {
				old_book_v.FieldByName(field_name).SetString(field_val)
			}
		}

		if new_book_v.Field(i).Kind() == reflect.Int64 && new_book_v.Field(i).Int() != 0 {
			// handle field, which type is int64 and not 0
			field_name := new_book_t.Field(i).Name
			field_val := new_book_v.Field(i).Int()

			if old_book_v.Field(i).CanSet() {
				old_book_v.FieldByName(field_name).SetInt(field_val)
			}
		}
	}

	log.JSONIndent(old_book)

	// save info, level plus one and return a new ID
	err := db.SaveBook(old_book)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}

	return &pb.GetBookInfoResp{Code: errs.Ok, Message: "ok", Data: old_book}, nil
}

func (s *BooksServer) GetBookInfoByISBN(ctx context.Context, req *pb.Book) (*pb.GetBookInfoResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "GetBookInfoByISBN", "%#v", req))

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

	// get from douban
	douban_book, err := douban.GetBookInfo(req.Isbn)
	// if err != nil && strings.Index(err.Error(), "404") == -1 {
	// 	log.Error(err)
	// 	return nil, errs.Wrap(errors.New(err.Error()))
	// }
	if err != nil {
		log.Error(err)
		// return nil, errs.Wrap(errors.New(err.Error()))
	}

	api_usage++

	if !bookInfoIsOk(douban_book) {
		// get from wanxiang
		wanxiang_book, err := wanxiang.GetBookInfo(req.Isbn)
		if err != nil {
			log.Error(err)
			return nil, errs.Wrap(errors.New(err.Error()))
		}
		api_usage++
		final_book = integreteInfo(douban_book, wanxiang_book)
	}

	if api_usage == 1 {
		final_book = douban_book
	}

	// API 调用之后，未找到该图书，return
	if final_book == nil {
		return &pb.GetBookInfoResp{Code: errs.Ok, Message: "book_not_found"}, nil
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
	return &pb.GetBookInfoResp{Code: errs.Ok, Message: "ok", Data: final_book}, nil
}
