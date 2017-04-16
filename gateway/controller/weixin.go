package controller

import (
	"net/http"

	"github.com/goushuyun/weixin-golang/misc"
	"github.com/wothing/log"
)

func ReceiveTicket(w http.ResponseWriter, r *http.Request) {
	log.Debug("My name is Wang Kai ..")

	// receive component_verify_ticket from weixin
	misc.RespondMessage(w, r, "success")
}
