package net

import "strings"

type HandlerFunc func(req *WsMsgReq, rsp *WsMsgRsp)

type group struct {
	prefix     string
	handlerMap map[string]HandlerFunc
}

func (g group) exec(name string, req *WsMsgReq, rsp *WsMsgRsp) {
	h := g.handlerMap[name]
	if h != nil {
		h(req, rsp)
	}
}

func (g *group) AddRoute(name string, handlerFunc HandlerFunc) {
	g.handlerMap[name] = handlerFunc
}

type router struct {
	group []*group
}

func NewRouter() *router {
	return &router{}
}

func (r *router) Run(req *WsMsgReq, rsp *WsMsgRsp) {
	strs := strings.Split(req.Body.Name, ".")
	prefix := ""
	name := ""
	if len(strs) == 2 {
		prefix = strs[0]
		name = strs[1]
	}
	for _, g := range r.group {
		if g.prefix == prefix {
			g.exec(name, req, rsp)
		}
	}
}

func (r *router) Group(prefix string) *group {
	g := &group{
		prefix: prefix,
	}
	r.group = append(r.group, g)
	return g
}
