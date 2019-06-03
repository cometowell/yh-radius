package radius

import (
	"github.com/go-xorm/xorm"
	"go-rad/model"
	"net"
)

type RadMiddleWare func(cxt *Context)

type Context struct {
	Request          RadiusPackage
	Response         *RadiusPackage
	User             *model.RadUser
	Listener         *net.UDPConn
	Dst              *net.UDPAddr
	RadNas           model.RadNas
	Index            int
	Handlers         []RadMiddleWare
	ReplyRadiusAttrs []RadiusAttr
	throwPackage     bool
	Session          *xorm.Session
}

func (cxt *Context) Next() {
	if cxt.Index >= len(cxt.Handlers) - 1 {
		return
	}
	cxt.Index += 1
	cxt.Handlers[cxt.Index](cxt)
}
