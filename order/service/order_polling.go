package service

import (
	"fmt"

	"github.com/robfig/cron"
)

//注册时间轮询
func RegisterOrderPolling(cron *cron.Cron) {
	order_close_task_spec := "0 0/2 * * * *"               //时间轮询表达式 每两分钟 执行一次
	order_statistic_before_dawn_task_spec := "0 0 1 * * *" //时间轮询表达式 每天凌晨1：00执行一次
	order_statistic_at_night_task_spec := "0 0 22 * * *"   //时间轮询表达式 每天晚上22：00执行一次
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
}

//订单统计核心处理方法
func orderStatisticHandle() {

}
