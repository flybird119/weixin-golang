package service

import (
	"17mei/errs"
	"errors"
	"fmt"
	"time"

	"github.com/goushuyun/weixin-golang/misc"
	"github.com/goushuyun/weixin-golang/pb"

	accountDB "github.com/goushuyun/weixin-golang/account/db"
	"github.com/goushuyun/weixin-golang/order/db"
	schoolDB "github.com/goushuyun/weixin-golang/school/db"
	storeDB "github.com/goushuyun/weixin-golang/store/db"

	"github.com/robfig/cron"
	"github.com/wothing/log"
)

//注册时间轮询
func RegisterOrderPolling(cron *cron.Cron) {
	order_close_task_spec := "0 0/5 * * * *"               //时间轮询表达式 每1分钟 执行一次
	order_statistic_before_dawn_task_spec := "0 0 1 * * *" //时间轮询表达式 每天凌晨1：00执行一次
	order_statistic_at_night_task_spec := "0 30 22 * * *"  //时间轮询表达式 每天晚上22：00执行一次
	order_system_confirm_task_spec := "0 20 0 * * *"       //系统自动确认订单轮询 每日凌晨 0:20
	// order_system_confirm_task_spec := "0 0/1 * * * *" //系统自动确认订单轮询 每日凌晨 0:20
	// order_statistic_before_dawn_task_spec := "0 0/1 * * * *" //时间轮询表达式 每天凌晨1：00执行一次
	// order_statistic_at_night_task_spec := "0 0/3 * * * *"    //时间轮询表达式 每天晚上22：00执行一次
	fmt.Println(order_close_task_spec + order_statistic_before_dawn_task_spec + order_statistic_at_night_task_spec + order_system_confirm_task_spec)
	//注册检查即将关闭的订单
	cron.AddFunc(order_close_task_spec, orderCloseHandle)
	//注册商家订单统计--时间点：凌晨
	cron.AddFunc(order_statistic_before_dawn_task_spec, orderStatisticHandle)
	//注册商家订单统计--时间点：晚上
	cron.AddFunc(order_statistic_at_night_task_spec, orderStatisticHandle)
	//系统定时任务-订单到时自动完成
	cron.AddFunc(order_system_confirm_task_spec, OrderConirmBySystemHandle)

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

//订单到时自动确认
func OrderConirmBySystemHandle() {
	//首先获取所有的可以系统判定完成的订单
	orders, err := db.GetAllWillCompletOrder()
	if err != nil {
		log.Error(err)
		return
	}
	log.Debugf("orders :%+v", orders)
	for i := 0; i < len(orders); i++ {
		order := orders[i]
		ConfirmOrderBySystem(order)
	}
}

// 确认订单（微信端）
func ConfirmOrderBySystem(in *pb.Order) error {
	//1.0 首先要检验 订单的状态 未发货的订单不能点击成功
	searchOrder := &pb.Order{Id: in.Id}
	err := db.GetOrderBaseInfo(searchOrder)
	if err != nil {
		log.Error(err)
		return err
	}
	if searchOrder.OrderStatus != 3 || searchOrder.ConfirmAt != 0 {
		return nil
	}
	//订单成功——>修改订单状态 +4
	in.OrderStatus = 4
	in.CompleteAt = 1
	err = db.UpdateOrder(in)
	if err != nil {
		log.Error(err)
		return nil
	}
	searchOrder.OrderStatus = 4
	//商家账户更改
	OrderCompleteAccountHandle(searchOrder)
	return nil
}

//订单完成账户操作
// -> 订单完成包含两种两种情况
// -> > 1:订单经过正常流程完成（用户下单-用户支付-商家发货-用户确认订单或者系统确认订单)
// -> > 2:订单售后完成，并且在订单成功(订单状态在用户确认订单或者系统确认订单)之前
func OrderCompleteAccountHandle(in *pb.Order) error {

	//首先查看相对应的记录
	hasSellerWithdrwalAccountItem, err := accountDB.HasExistAcoount(&pb.AccountItem{StoreId: in.StoreId, OrderId: in.Id, ItemType: 1})
	if err != nil {
		go misc.LogErrOrder(in, "订单完成-检查记录是否存在发生错误,影响商户可提现和待结算的转换问题", err)
		log.Error(err)
		return errs.Wrap(errors.New(err.Error()))
	}
	hasSellerBalanceAccountItem, err := accountDB.HasExistAcoount(&pb.AccountItem{StoreId: in.StoreId, OrderId: in.Id, ItemType: 17})
	if err != nil {
		go misc.LogErrOrder(in, "订单完成-检查记录是否存在发生错误 ,影响商户可提现和待结算的转换问题", err)
		log.Error(err)
		return errs.Wrap(errors.New(err.Error()))
	}
	if hasSellerWithdrwalAccountItem || hasSellerBalanceAccountItem {
		go misc.LogErrOrder(in, "订单完成-需要检查订单，影响商户可提现和待结算的转换问题", errors.New("修改失败，改数据已经在账单项中存在"))
		return errs.Wrap(errors.New("修改失败，改数据已经在账单项中存在"))
	}
	//1 需要更改两个值 1 ：待体现金额  2 可提现金额
	//2 需要记录两条记录 1 ：待体现金额资金流向记录  2 ： 可提现金额资金流入记录
	//1.1	待体现金额

	sellerAccountWithdrawal := &pb.Account{StoreId: in.StoreId, UnsettledBalance: -in.WithdrawalFee}
	err = accountDB.ChangeAccountWithdrawalFee(sellerAccountWithdrawal)
	if err != nil {
		go misc.LogErrOrder(in, "订单完成-待体现转化成可提现发生错误 ,影响商户可提现和待结算的转换", err)
		log.Error(err)
		return errs.Wrap(errors.New(err.Error()))
	}
	sellerAccountWithdrawalItem := &pb.AccountItem{UserType: 1, StoreId: in.StoreId, OrderId: in.Id, ItemType: 1, Remark: "已由待结算金额转换为可提现金额", ItemFee: -in.WithdrawalFee, AccountBalance: sellerAccountWithdrawal.UnsettledBalance}

	//1.2 	待体现金额资金流向记录
	err = accountDB.AddAccountItem(sellerAccountWithdrawalItem)
	if err != nil {
		go misc.LogErrAccount(sellerAccountWithdrawalItem, "订单完成-增加待结算转化可结算时发生错误 ,影响下一步更改可体现金额和增加资金流向记录的操作", err)
		log.Error(err)
		return errs.Wrap(errors.New(err.Error()))
	}
	sellerAccountBalance := &pb.Account{StoreId: in.StoreId, Balance: in.WithdrawalFee}
	//2.1：	待体现金额资金流向记录
	err = accountDB.ChangAccountBalance(sellerAccountBalance)
	if err != nil {
		go misc.LogErrOrder(in, "订单完成-更改可结算金额操作时发生错误 ,影响商户可提现和待结算的转换", err)
		log.Error(err)
		return errs.Wrap(errors.New(err.Error()))
	}
	//2.2	可提现金额资金流入记录
	sellerAccountBalanceItem := &pb.AccountItem{UserType: 1, StoreId: in.StoreId, OrderId: in.Id, ItemType: 17, Remark: "由待结算金额转换为可提现金额", ItemFee: in.WithdrawalFee, AccountBalance: sellerAccountBalance.Balance}

	err = accountDB.AddAccountItem(sellerAccountBalanceItem)
	if err != nil {
		go misc.LogErrAccount(sellerAccountWithdrawalItem, "订单完成-记录可结算操作时发生错误 ,影响操作日志", err)
		log.Error(err)
		return errs.Wrap(errors.New(err.Error()))
	}

	return nil
}
