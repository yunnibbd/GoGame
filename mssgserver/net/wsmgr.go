package net

import "sync"

var Mgr = &WsMgr{
	userCache: make(map[int]WSConn),
}

type WsMgr struct {
	uc sync.RWMutex
	userCache map[int]WSConn
}

func (m *WsMgr) UserLogin(conn WSConn, )
