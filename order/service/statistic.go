package service

import (
	"golang.org/x/net/context"

	"github.com/goushuyun/weixin-golang/misc"

	"github.com/goushuyun/weixin-golang/pb"
	"github.com/wothing/log"
)

//售后订单处理结果
func (s *OrderServiceServer) StoreDailyGoodsSalesStatistic(ctx context.Context, in *pb.Store) (*pb.Void, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "StoreDailyGoodsSalesStatistic", "%#v", in))
	//1.0 首先获取这个云店铺下的所有学校

	//2.0 根据学校统计每日销售额
	return &pb.Void{}, nil
}
