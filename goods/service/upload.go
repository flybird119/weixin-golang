package service

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc/metadata"

	bookService "github.com/goushuyun/weixin-golang/books/service"
	"github.com/goushuyun/weixin-golang/errs"
	"github.com/goushuyun/weixin-golang/goods/db"
	"github.com/goushuyun/weixin-golang/mediastore/service"
	"github.com/goushuyun/weixin-golang/misc"
	"github.com/goushuyun/weixin-golang/pb"
	"github.com/pborman/uuid"
	"github.com/tealeg/xlsx"
	"github.com/wothing/log"
	"golang.org/x/net/context"
)

//AddGoods 增加商品
func (s *GoodsServiceServer) GoodsBactchUploadOperate(ctx context.Context, in *pb.GoodsBatchUploadModel) (*pb.NormalResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "GoodsBactchUploadOperate", "%#v", in))
	// 1 首先保存记录
	err := db.AddBatchUpload(in)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}

	go coreUploadHandler(in)
	return &pb.NormalResp{Code: "00000", Message: "ok"}, nil
}

////获取批量上传数据
func (s *GoodsServiceServer) GoodsBactchUploadList(ctx context.Context, in *pb.GoodsBatchUploadModel) (*pb.GoodsBatchUploadModelListResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "GoodsBactchUploadList", "%#v", in))
	// 1 首先保存记录
	models, err, totalCount := db.GoodsBactchUploadList(in)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	return &pb.GoodsBatchUploadModelListResp{Code: "00000", Message: "ok", Data: models, TotalCount: totalCount}, nil
}

//下载文件
func downloadRemoteExcel(originFileUrl string, filename string) {
	res, _ := http.Get(originFileUrl)
	file, _ := os.Create(filename)
	io.Copy(file, res.Body)

}

func coreUploadHandler(in *pb.GoodsBatchUploadModel) {
	ctx := metadata.NewContext(context.Background(), metadata.Pairs("tid", uuid.New()))

	//声明全局失败记录
	var failedRechord []*pb.Goods
	updateUploadModel := &pb.GoodsBatchUploadModel{Id: in.Id}

	// 2 下载表格文件
	splitStringArray := strings.Split(in.OriginFile, "/")
	filename := splitStringArray[len(splitStringArray)-1]
	downloadRemoteExcel(in.OriginFile, filename)
	// 3 读取文件 ，并获取列表
	goodsList, err := readExcel(filename)
	if err != nil {
		log.Error(err)
		//错误处理
		updateUploadModel.State = 2
		updateUploadModel.ErrorReason = "文件读取失败，请核实后再试"
		db.UpdateBatchUpload(updateUploadModel)
		return
	}
	log.Debug(goodsList)
	//如果商品列表为空，那么文件上传失败
	if len(goodsList) <= 0 {
		//失败操作
		updateUploadModel.State = 2
		updateUploadModel.ErrorReason = "文件无上传数据，请核实后再试"
		db.UpdateBatchUpload(updateUploadModel)
		return
	}
	//获取cpu数量
	cpuNum := runtime.NumCPU() * 4

	size := len(goodsList) / cpuNum
	//4 设置 batch_size 获取批量上传数据列表
	spiltList, _ := splitGoodsList(size, goodsList)

	//定义传输通道 -- 模拟协程信号通道（数据返回）
	goodsChan := make(chan pb.Goods)
	var currentCompleteNum int
	//定义统计通道 -- 用于判断任务有没有处理完成
	statisticChan := make(chan int)
	var wg sync.WaitGroup
	fmt.Println("uploadStart")

	for i := 0; i < len(spiltList); i++ {
		wg.Add(1)
		handleList := spiltList[i]
		//5 ****多协程处理数据
		go func(routineList []*pb.Goods) {
			defer wg.Done()
			for k := 0; k < len(routineList); k++ {
				penddingGoods := routineList[k]
				fmt.Println(penddingGoods == nil)
				if penddingGoods == nil {
					continue
				}
				penddingGoods.StoreId = in.StoreId
				//定时器
				timeout := make(chan bool)
				go func() {
					time.Sleep(30 * time.Second) // 设置查询超时时间
					timeout <- true
				}()
				errChan := make(chan error)
				go handlePenddingGoods(ctx, penddingGoods, in.Discount, in.Type, in.StorehouseId, in.ShelfId, in.FloorId, errChan)
				select {
				case err, _ = <-errChan:
					log.Debug("完成")
					close(errChan)
					break
				case <-timeout:
					log.Debug("timeout")
					close(timeout)
					err = errors.New("上传超时，请注意事后核对该图书")
					break
				}
				if err != nil {
					goodsChan <- pb.Goods{Isbn: penddingGoods.Isbn, StrNum: penddingGoods.StrNum, ErrMsg: err.Error()}
				}
				statisticChan <- 1
			}
		}(handleList)

	}
	//6 通过协程通道 获取错误的返回列表
	for {
		var goods pb.Goods
		var singleValue int
		var ok bool
		select {
		case goods, ok = <-goodsChan:
			if ok {
				failedRechord = append(failedRechord, &goods)
			}
		case singleValue, _ = <-statisticChan:
			currentCompleteNum += singleValue
		}

		if currentCompleteNum == len(goodsList) {
			fmt.Println(currentCompleteNum)
			fmt.Println(len(goodsList))
			close(statisticChan)
			close(goodsChan)
			//完成统计
			break
		}
	}
	wg.Wait()

	fmt.Println("uploadOver")
	//构建错误列表
	if len(failedRechord) > 0 {
		failedFileUrl, err := createFailedExcel(failedRechord)
		if err != nil {
			updateUploadModel.ErrorReason = "保存错误文件失败"
		} else {
			updateUploadModel.ErrorFile = failedFileUrl
		}

	}

	//数量
	updateUploadModel.SuccessNum = int64(len(goodsList) - len(failedRechord))
	updateUploadModel.FailedNum = int64(len(failedRechord))
	//状态更新
	updateUploadModel.State = 3
	//更新操作
	db.UpdateBatchUpload(updateUploadModel)

	// 最后一步：删除临时文件
	os.Remove(filename)
}

//读取文件
func readExcel(name string) (books []*pb.Goods, err error) {
	//判定 文件格式
	//
	//-------------暂时解除对xls文件的解析---------------//

	// reg := regexp.MustCompile("\\.xlsx$")
	// format := reg.FindString(name)
	// if format == "" {
	// 	return readExcelByXls(name)
	// } else {
	// 	return readExcelByXlsx(name)
	// }
	//-----------------------------------------------//

	return readExcelByXlsx(name)

}

//xlsx 格式读取文件
func readExcelByXlsx(name string) (books []*pb.Goods, err error) {
	xlFile, err := xlsx.OpenFile(name)
	if err != nil {
		log.Error(err)
		return
	}

	sheet := xlFile.Sheets[0]
	var index int
	for _, row := range sheet.Rows {
		if index == 0 {
			index++
			continue
		}
		isbn, _ := row.Cells[0].String()
		numStr, _ := row.Cells[1].String()
		if isbn == "" || numStr == "" {
			break
		}
		book := &pb.Goods{Isbn: isbn, StrNum: numStr}
		books = append(books, book)
		index++
	}
	return
}

//xls 格式读取文件
func readExcelByXls(name string) (books []*pb.Goods, err error) {

	return
}

// xlsx 写文件 并上传到七牛云

func createFailedExcel(failedRechord []*pb.Goods) (fileUrl string, err error) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row

	file = xlsx.NewFile()
	sheet, err = file.AddSheet("列表")
	if err != nil {
		log.Error(err)
		return
	}
	row = sheet.AddRow()
	head := []string{"ISBN", "数量", "错误原因"}
	row.WriteSlice(&head, len(head))

	for _, goods := range failedRechord {
		row = sheet.AddRow()
		row.AddCell().SetString(goods.Isbn)
		row.AddCell().SetString(goods.StrNum)
		row.AddCell().SetString(goods.ErrMsg)
	}

	fileUrl, err = service.PutExcelFile(file)
	if err != nil {
		log.Debug(err)
		return
	}
	log.Debug(fileUrl)

	//上传到七牛云

	return
}

// goodsList 分组
func splitGoodsList(batchSize int, goodsList []*pb.Goods) (splitList [][]*pb.Goods, err error) {

	for i := 0; i < len(goodsList); i += batchSize {
		if i+batchSize >= len(goodsList) {
			splitList = append(splitList, goodsList[i:])
		} else {
			splitList = append(splitList, goodsList[i:i+batchSize])
		}
	}
	return
}

//handlePenddingGoods 处理每条数据
func handlePenddingGoods(ctx context.Context, goods *pb.Goods, discount, goodsType int64, storeId, shelfId, floorId string, errChan chan error) error {
	//转化 数量
	if goods == nil {
		errChan <- errors.New("数据错误")
		return errors.New("数据错误")
	}
	fmt.Printf("num ===== '%s'\n", goods.StrNum)
	num, err := strconv.ParseInt(goods.StrNum, 10, 64)
	log.Debug("===========num:%d", num)
	if err != nil {
		log.Debug("图书数量不合法")
		errChan <- errors.New("图书数量不合法")
		return errors.New("图书数量不合法")
	}
	//校验isbn是否正确
	reg := regexp.MustCompile("^(\\d[- ]*){12}[\\d]$")
	isbn := reg.FindString(goods.Isbn)
	isbn = strings.Replace(isbn, "-", "", -1)
	isbn = strings.Replace(isbn, " ", "", -1)
	if isbn == "" {
		log.Debug("isbn不正确")
		errChan <- errors.New("isbn不正确")
		return errors.New("isbn不正确")
	}
	//查找图书信息

	book, err := bookService.GetBookInfoByISBNWithNoContext(&pb.Book{Isbn: goods.Isbn, UploadWay: "batch"})
	if err != nil {
		log.Error(err)
		return err
	}
	if book == nil {
		log.Debug("没找到图书")
		errChan <- errors.New("未找到该图书，请手动上传")
		return errors.New("未找到该图书，请手动上传")
	}
	goods.BookId = book.Id
	//计算图书价格
	var serviceDiscount = float64(discount) / 100
	withdrawalFeeStr := fmt.Sprintf("%0.0f", float64(book.Price)*(serviceDiscount))
	price, err := strconv.ParseInt(withdrawalFeeStr, 10, 64)
	if err != nil {
		log.Debug(err)
		errChan <- err
		return err
	}
	//构造位置信息
	var locations []*pb.GoodsLocation
	location := &pb.GoodsLocation{Type: goodsType, StorehouseId: storeId, ShelfId: shelfId, FloorId: floorId, Amount: num, Price: price}
	locations = append(locations, location)
	goods.Location = locations

	log.Debugf("%#v", goods)
	//增加商品
	err = db.AddGoods(goods)
	if err != nil {
		log.Error(err)
		errChan <- err
		return err
	}
	errChan <- nil
	return nil
}
