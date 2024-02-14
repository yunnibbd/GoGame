package login

import (
	"mssgserver/db"
	"mssgserver/net"
	"mssgserver/server/login/controller"
)

var Router = net.NewRouter()

func Init() {
	db.TestDB()
	initRouter()
}

func initRouter() {
	controller.DefaultAccount.Router(Router)
}
