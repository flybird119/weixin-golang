package service

import (
	"17mei/errs"
	"errors"
	"time"

	"golang.org/x/net/context"

	"github.com/goushuyun/weixin-golang/misc"
	orderDB "github.com/goushuyun/weixin-golang/order/db"
	"github.com/goushuyun/weixin-golang/pb"
	"github.com/wothing/log"
)

type StatisticServiceServer struct{}

//今日销售额 必添字段 store_id 非必填字段 school_id
func (s *StatisticServiceServer) StatisticToday(ctx context.Context, in *pb.GoodsSalesStatisticModel) (*pb.StatisticTodayResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "StatisticToday", "%#v", in))

	//return nil, errs.Wrap(errors.New(err.Error()))
	now := time.Now()
	statisticDateStr := now.Format("2006-01-02")
	in.StatisticAt = statisticDateStr
	err := orderDB.GetOneDaySales(in)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	//
	log.Debug("=========================")
	log.Debugf("%+v", in)
	log.Debug("=========================")
	todayModel := &pb.StatisticTotalModel{}
	todayModel.OnlineTotalSales = in.AlipayOrderFee + in.WechatOrderFee
	todayModel.OfflineTotalSales = in.OfflineNewBookSalesFee + in.OfflineOldBookSalesFee
	return &pb.StatisticTodayResp{Code: "00000", Message: "ok", Data: todayModel}, nil
}

// 总计统计
func (s *StatisticServiceServer) StatisticTotal(ctx context.Context, in *pb.GoodsSalesStatisticModel) (*pb.StatisticTotalResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "StatisticTotal", "%#v", in))

	//return nil, errs.Wrap(errors.New(err.Error()))
	//首先获取昨天销售额
	now := time.Now()
	now = now.Add(-1 * 24 * time.Hour)
	statisticDateStr := now.Format("2006-01-02")
	in.StatisticAt = statisticDateStr
	err := orderDB.GetOneDaySales(in)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	yesterdayModel := &pb.StatisticTotalModel{}
	yesterdayModel.OnlineTotalSales = in.AlipayOrderFee + in.WechatOrderFee
	yesterdayModel.OfflineTotalSales = in.OfflineNewBookSalesFee + in.OfflineOldBookSalesFee
	//统计历史的
	totalModel, err := orderDB.HistoryTotalSales(in)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	return &pb.StatisticTotalResp{Code: "00000", Message: "ok", Data: &pb.StatisticTotalData{YesterdaySales: yesterdayModel, TotalSales: totalModel}}, nil
}

// 统计详情列表
func (s *StatisticServiceServer) StatisticDaliy(ctx context.Context, in *pb.GoodsSalesStatisticModel) (*pb.StatisticDaliyResp, error) {

	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "StatisticDaliy", "%#v", in))
	statisticModels, err := orderDB.HistoryDaliySales(in)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}

	return &pb.StatisticDaliyResp{Code: "00000", Message: "ok", Data: statisticModels}, nil
}

// 月统计
func (s *StatisticServiceServer) StatisticMonth(ctx context.Context, in *pb.GoodsSalesStatisticModel) (*pb.StatisticMonthResp, error) {

	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "StatisticMonth", "%#v", in))
	//return nil, errs.Wrap(errors.New(err.Error()))
	salesModels, err := orderDB.HistoryMonthSales(in)
	if err != nil {
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	return &pb.StatisticMonthResp{Code: "00000", Message: "ok", Data: salesModels}, nil
}
