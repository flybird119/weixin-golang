package main

import (
	"strconv"
	"testing"
	"time"

	"github.com/wothing/log"

	"github.com/goushuyun/weixin-golang/pb"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var (
	smsClient pb.SMSServiceClient
	ctx       context.Context
)

func init() {
	address := "127.0.0.1" + ":" + strconv.FormatFloat(8850, 'g', -1, 64)
	log.Debug(address)
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithTimeout(time.Second*3))
	if err != nil {
		log.Fatalf("connect to sms service fail")
	}

	smsClient = pb.NewSMSServiceClient(conn)
	ctx = metadata.NewContext(context.Background(), metadata.Pairs("tid", "user-test"))
}

func TestSms(t *testing.T) {
	message := []string{"35421"}
	a := &pb.SMSReq{Type: pb.SMSType_CommonCheckCode, Mobile: "86-13122210065", Message: message}
	smsClient.SendSMS(ctx, a)
}

func TestType(t *testing.T) {
	t.Logf("-----%d----", pb.SMSType_CommonCheckCode)
}
