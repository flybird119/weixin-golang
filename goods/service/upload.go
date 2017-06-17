package service

import (
	"17mei/errs"
	"errors"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/goushuyun/weixin-golang/misc"
	"github.com/goushuyun/weixin-golang/pb"
	"github.com/tealeg/xlsx"
	"github.com/wothing/log"
	"golang.org/x/net/context"
)

//AddGoods 增加商品
func (s *GoodsServiceServer) GoodsBactchUploadOperate(ctx context.Context, in *pb.GoodsBatchUploadModel) (*pb.NormalResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "GoodsBactchUploadOperate", "%#v", in))

	// 1 首先保存记录

	// 2 下载表格文件
	splitStringArray := strings.Split(in.OriginFile, "/")
	filename := splitStringArray[len(splitStringArray)-1]
	// 3 读取文件 ，并获取列表
	books, err := readExcel(filename)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	log.Debugf("%+v", books)

	//4 分组批量上传 设置 batch_size

	// 最后删除文件
	os.Remove(filename)
	return &pb.NormalResp{Code: "00000", Message: "ok"}, nil
}

//下载文件
func downloadRemoteExcel(originFileUrl string, filename string) {
	res, _ := http.Get(originFileUrl)
	file, _ := os.Create(filename)
	io.Copy(file, res.Body)

}

//读取文件
func readExcel(name string) (books []*pb.Book, err error) {
	//判定 文件格式
	reg := regexp.MustCompile("\\.xlsx$")
	format := reg.FindString(name)
	if format == "" {
		return readExcelByXls(name)
	} else {
		return readExcelByXlsx(name)
	}
}

//xlsx 格式读取文件
func readExcelByXlsx(name string) (books []*pb.Book, err error) {
	xlFile, err := xlsx.OpenFile(name)
	if err != nil {
		log.Error(err)
		return
	}

	sheet := xlFile.Sheets[0]
	for _, row := range sheet.Rows {
		isbn, _ := row.Cells[0].String()
		num, _ := row.Cells[0].String()
		book := &pb.Book{Isbn: isbn, InfoSrc: num}
		books = append(books, book)
	}
	return
}

//xls 格式读取文件
func readExcelByXls(name string) (books []*pb.Book, err error) {

	return
}

// xlsx 写文件

// 上传文件到七牛云
