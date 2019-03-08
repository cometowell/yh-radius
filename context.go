package main

import (
	"github.com/go-xorm/xorm"
	"net"
)

type RadMiddleWare func(cxt *Context)

type Context struct {
	Request          RadiusPackage
	Response         *RadiusPackage
	User             *RadUser
	Listener         *net.UDPConn
	Dst              *net.UDPAddr
	RadNas              RadNas
	index            int
	Handlers         []RadMiddleWare
	ReplyRadiusAttrs []RadiusAttr
	throwPackage     bool
	Session *xorm.Session
}

func (cxt *Context) Next() {
	if cxt.index >= len(cxt.Handlers) - 1 {
		return
	}
	cxt.index += 1
	cxt.Handlers[cxt.index](cxt)
}
