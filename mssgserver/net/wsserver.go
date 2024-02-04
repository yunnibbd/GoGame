package net

import (
	"encoding/json"
	"fmt"
	"github.com/forgoer/openssl"
	"github.com/gorilla/websocket"
	"log"
	"mssgserver/utils"
	"sync"
)

type wsServer struct {
	wsConn       *websocket.Conn
	router       *router
	outChan      chan *WsMsgRsp
	Seq          int64
	property     map[string]interface{}
	propertyLock sync.RWMutex
}

func NewWsServer(wsConn *websocket.Conn) *wsServer {
	return &wsServer{
		wsConn:   wsConn,
		outChan:  make(chan *WsMsgRsp, 1000),
		property: make(map[string]interface{}),
		Seq:      0,
	}
}

func (w *wsServer) Router(router *router) {
	w.router = router
}

func (w *wsServer) SetProperty(key string, value interface{}) {
	w.propertyLock.Lock()
	defer w.propertyLock.Unlock()
	w.property[key] = value
}

func (w *wsServer) GetProperty(key string) (interface{}, error) {
	w.propertyLock.Lock()
	defer w.propertyLock.Unlock()
	return w.property[key], nil
}

func (w *wsServer) RemoveProperty(key string) {
	w.propertyLock.Lock()
	defer w.propertyLock.Unlock()
	delete(w.property, key)
}

func (w *wsServer) Addr() string {
	return w.wsConn.RemoteAddr().String()
}

func (w *wsServer) Push(name string, data interface{}) {
	rsp := &WsMsgRsp{Body: &RspBody{Name: name, Msg: data, Seq: 0}}
	w.outChan <- rsp
}

func (w *wsServer) Start() {
	//启动读写数据的处理逻辑
	go w.readMsgLoop()
	go w.writeMsgLoop()
}

func (w *wsServer) readMsgLoop() {
	//先读到客户端发过来的数据 再进行处理 再回复消息
	//经过路由 实际处理程序
	defer func() {
		if err := recover(); err != nil {
			log.Fatal(err)
			w.Close()
		}
	}()
	for {
		_, data, err := w.wsConn.ReadMessage()
		if err != nil {
			log.Println("收消息出现错误:", err)
		}
		//收到消息 解析消息
		//data解压unzip
		data, err = utils.UnZip(data)
		if err != nil {
			log.Println("解压数据出错，格式不合法:", err)
			continue
		}
		//加密消息 解密
		secretKey, err := w.GetProperty("secretKey")
		if err == nil {
			key := secretKey.(string)
			d, err := utils.AesCBCDecrypt(data, []byte(key), []byte(key), openssl.ZEROS_PADDING)
			if err != nil {
				log.Println("数据格式有误，解密失败:", err)
				//出错后 发起握手
				//w.Handshake()
			} else {
				data = d
			}
		}
		//data 转为 body
		body := &ReqBody{}
		err = json.Unmarshal(data, body)
		if err != nil {
			log.Println("数据格式有误，非法格式:", err)
		} else {
			//获取到前端传递的数据了, 拿上这些数据 去具体的业务进行处理
			req := &WsMsgReq{Conn: w, Body: body}
			rsp := &WsMsgRsp{Body: &RspBody{Name: body.Name, Seq: req.Body.Seq}}
			w.router.Run(req, rsp)
			w.outChan <- rsp
		}
	}
	w.Close()
}

func (w *wsServer) writeMsgLoop() {
	for {
		select {
		case msg := <-w.outChan:
			fmt.Println(msg)
		}
	}
}

func (w *wsServer) Close() {
	_ = w.wsConn.Close()
}
