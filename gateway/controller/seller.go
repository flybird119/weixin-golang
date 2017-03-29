package controller

import (
	"fmt"
	"net/http"

	"github.com/wothing/log"
)

func SellerLogin(w http.ResponseWriter, r *http.Request) {
	log.Info("My name is Wang Kai ...")

	fmt.Fprintf(w, "Hello, world")
}
