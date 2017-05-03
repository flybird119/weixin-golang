package service

import (
	"fmt"
	"time"

	"github.com/goushuyun/weixin-golang/pb"

	"github.com/goushuyun/weixin-golang/order/db"
	schoolDB "github.com/goushuyun/weixin-golang/school/db"
	storeDB "github.com/goushuyun/weixin-golang/store/db"
	"github.com/robfig/cron"
	"github.com/wothing/log"
)

//注册时间轮询
func RegisterOrderPolling(cron *cron.Cron) {
	order_close_task_spec := "0 0/1 * * * *"               //时间轮询表达式 每1分钟 执行一次
	order_statistic_before_dawn_task_spec := "0 0 1 * * *" //时间轮询表达式 每天凌晨1：00执行一次
	order_statistic_at_night_task_spec := "0 0 22 * * *"   //时间轮询表达式 每天晚上22：00执行一次
	// order_statistic_before_dawn_task_spec := "0 0/1 * * * *" //时间轮询表达式 每天凌晨1：00执行一次
	// order_statistic_at_night_task_spec := "0 0/3 * * * *"    //时间轮询表达式 每天晚上22：00执行一次
	fmt.Println(order_close_task_spec + order_statistic_before_dawn_task_spec + order_statistic_at_night_task_spec)
	//注册检查即将关闭的订单
	cron.AddFunc(order_close_task_spec, orderCloseHandle)
	//注册商家订单统计--时间点：凌晨
	cron.AddFunc(order_statistic_before_dawn_task_spec, orderStatisticHandle)
	//注册商家订单统计--时间点：晚上
	cron.AddFunc(order_statistic_at_night_task_spec, orderStatisticHandle)
}

//关闭订单核心处理方法
func orderCloseHandle() {
	//首先查找符合条件的订单
	orders, err := db.FindAllExpireOrder()
	if err != nil {
		log.Error()
	}
	log.Debug("=============系统取消订单==================")
	log.Debugf("%+v", orders)
	log.Debug("=========================================")
	for i := 0; i < len(orders); i++ {
		go db.CloseOrder(orders[i])
	}

}

//订单统计核心处理方法
func orderStatisticHandle() {

	log.Debug("===============================")
	log.Debugf("%#v", "订单统计")
	log.Debug("===============================")
	//查找说有云店
	stores, err := storeDB.FindAllStores()
	if err != nil {
		log.Error(err)
		return
	}
	for i := 0; i < len(stores); i++ {
		go StoreOrdersStatistic(stores[i])
	}
}

// 店铺统计
func StoreOrdersStatistic(store *pb.Store) error {
	schools, err := schoolDB.GetSchoolsByStore(store.Id)
	if err != nil {
		log.Error(err)
		log.Warnf("订单统计发生错误，错误原因：%+v", err)
		return err
	}
	for i := 0; i < len(schools); i++ {
		go SchoolOrdersStatistic(schools[i])
	}
	return nil
}

//学校统计
func SchoolOrdersStatistic(school *pb.School) error {
	now := time.Now()
	statisticDate := now.Add(-1 * 24 * time.Hour)
	statisticDateStr := statisticDate.Format("2006-01-02")
	log.Debugf("date:%+v; dateStr:%+v", statisticDate, statisticDateStr)
	//首先要判断改日期的数据
	isExist, err := db.HasThisDayGoodsSalesData(school.Id, statisticDateStr)
	if err != nil {
		log.Error(err)
		return nil
	}
	//如果存在 过滤
	if isExist {
		log.Debugf("school:%s has this day :%s data", school.Id, statisticDateStr)
		return nil
	}
	goodsSalesStatisticModel := &pb.GoodsSalesStatisticModel{}
	goodsSalesStatisticModel.SchoolId = school.Id
	goodsSalesStatisticModel.StoreId = school.StoreId
	goodsSalesStatisticModel.StatisticAt = statisticDateStr
	//统计线上销售金额
	err = db.OnlineGoodsSalesStatistic(goodsSalesStatisticModel)
	if err != nil {
		log.Error(err)
	}
	//统计线下销售金额
	err = db.OfflineGoodsSalesStatistic(goodsSalesStatisticModel)
	if err != nil {
		log.Error(err)
	}
	log.Debug("===========================出炉喽======================================")
	log.Debugf("%+v", goodsSalesStatisticModel)
	log.Debug("=================================================================")
	err = db.AddGoodsSalesStatistic(goodsSalesStatisticModel, statisticDate)
	if err != nil {
		log.Error(err)
	}
	return nil
}
