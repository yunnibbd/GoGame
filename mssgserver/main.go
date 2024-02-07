package main

import (
	"mssgserver/config"
	"mssgserver/net"
	"mssgserver/server/login"
)

func main() {
	host := config.File.MustValue("LoginServer", "host", "127.0.0.1")
	port := config.File.MustValue("LoginServer", "port", "8004")
	s := net.NewServer(host + ":" + port)
	login.Init()
	s.Router(login.Router)
	s.Start()
}
