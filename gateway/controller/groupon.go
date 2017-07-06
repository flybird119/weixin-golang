package controller

import (
	"net/http"

	"github.com/goushuyun/weixin-golang/misc/token"

	"github.com/goushuyun/weixin-golang/misc"
	"github.com/goushuyun/weixin-golang/pb"
)

//SharedMajorBatchSave 通用专业批量增加
func SharedMajorBatchSave(w http.ResponseWriter, r *http.Request) {
	req := &pb.SharedMajor{}
	// call RPC to handle request
	misc.CallWithResp(w, r, "bc_groupon", "SharedMajorBatchSave", req)
}

//SharedMajorList 获取专业列表（筛选获取）
func SharedMajorList(w http.ResponseWriter, r *http.Request) {
	req := &pb.SharedMajor{}
	misc.CallWithResp(w, r, "bc_groupon", "SharedMajorList", req)
}

//创建学校的学院
func SaveSchoolInstitute(w http.ResponseWriter, r *http.Request) {
	c := token.Get(r)
	//检测token
	if c == nil || c.StoreId == "" {
		misc.ReturnNotToken(w, r)
		return
	}

	req := &pb.SchoolInstitute{}
	misc.CallWithResp(w, r, "bc_groupon", "SaveSchoolInstitute", req, "school_id", "name")
}

//创建学院专业
func SaveInstituteMajor(w http.ResponseWriter, r *http.Request) {
	c := token.Get(r)
	//检测token
	if c == nil || c.StoreId == "" {
		misc.ReturnNotToken(w, r)
		return
	}

	req := &pb.InstituteMajor{}
	misc.CallWithResp(w, r, "bc_groupon", "SaveInstituteMajor", req, "institute_id", "name")
}

//获取学校学院专业列表
func GetSchoolMajorInfo(w http.ResponseWriter, r *http.Request) {
	c := token.Get(r)
	//检测token
	if c == nil || c.StoreId == "" {
		misc.ReturnNotToken(w, r)
		return
	}
	req := &pb.SchoolMajorInfoReq{StoreId: c.StoreId}
	misc.CallWithResp(w, r, "bc_groupon", "GetSchoolMajorInfo", req)
}
