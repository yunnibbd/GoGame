package login

import (
	"mssgserver/net"
	"mssgserver/server/login/controller"
)

var Router = net.NewRouter()

func Init() {
	initRouter()
}

func initRouter() {
	controller.DefaultAccount.Router(Router)
}
