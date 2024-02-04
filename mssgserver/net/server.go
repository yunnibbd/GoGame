package net

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type server struct {
	addr   string
	router *router
}

func NewServer(addr string) *server {
	return &server{
		addr: addr,
	}
}

func (s *server) Router(router *router) {
	s.router = router
}

func (s *server) Start() {
	http.HandleFunc("/", s.wsHandler)
	err := http.ListenAndServe(s.addr, nil)
	if err != nil {
		panic(err)
	}
}

var wsUpgrader = websocket.Upgrader{
	//允许所有CORS跨域请求
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (s *server) wsHandler(w http.ResponseWriter, r *http.Request) {
	wsConn, err := wsUpgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Fatal("websocket服务连接出错", err)
	}

	wsServer := NewWsServer(wsConn)
	wsServer.Router(s.router)
	wsServer.Start()
}
