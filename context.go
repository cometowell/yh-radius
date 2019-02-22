package main

import "net"

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
}

func (cxt *Context) Next() {
	if cxt.index >= len(cxt.Handlers) {
		return
	}
	cxt.index += 1
	cxt.Handlers[cxt.index](cxt)
}
