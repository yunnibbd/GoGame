package controller

import (
	"github.com/mitchellh/mapstructure"
	"log"
	"mssgserver/constant"
	"mssgserver/db"
	"mssgserver/net"
	"mssgserver/server/login/model"
	"mssgserver/server/login/proto"
	"mssgserver/utils"
	"time"
)

var DefaultAccount = &Account{}

type Account struct {
}

func (a *Account) Router(r *net.Router) {
	g := r.Group("account")
	g.AddRoute("login", a.login)
}

func (a *Account) login(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	loginReq := &proto.LoginReq{}
	loginRes := &proto.LoginRsp{}
	mapstructure.Decode(req.Body.Msg, loginReq)
	user := &model.User{}
	ok, err := db.Engine.Table(user).Where("username=?", loginReq.Username).Get(user)
	if err != nil {
		log.Println("用户表查询出错", err)
		return
	}
	if !ok {
		rsp.Body.Code = constant.UserNotExist
		return
	}
	pwd := utils.Password(loginReq.Password, user.Passcode)
	if pwd != user.Passwd {
		rsp.Body.Code = constant.PwdIncorrect
		return
	}
	//jwt
	token, _ := utils.Award(user.UId)

	rsp.Body.Code = constant.OK
	loginRes.UId = user.UId
	loginRes.Username = user.Username
	loginRes.Session = token
	loginRes.Password = ""
	rsp.Body.Msg = loginRes

	//保存用户登录记录
	ul := &model.LoginHistory{
		UId: user.UId, CTime: time.Now(), Ip: loginReq.Ip,
		Hardware: loginReq.Hardware, State: model.Login,
	}
	db.Engine.Table(ul).Insert(ul)
	//最后一次登录对的记录
	ll := &model.LoginLast{}
	ok, _ = db.Engine.Table(ll).Where("uid=?", user.UId).Get(ll)
	if ok {
		ll.IsLogout = 0
		ll.Ip = loginReq.Ip
		ll.LoginTime = time.Now()
		ll.Session = token
		db.Engine.Table(ll).Update(ll)
	} else {
		ll.IsLogout = 0
		ll.Ip = loginReq.Ip
		ll.LoginTime = time.Now()
		ll.Session = token
		ll.UId = user.UId
		db.Engine.Table(ll).Insert(ll)
	}
	//缓存一下 此用户和当前的ws连接
}
