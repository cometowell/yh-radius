package main

func authReply(cxt *Context, replyCode byte, msg string) {
	cxt.Response.Code = replyCode
	replyMessage := RadiusAttr{
		AttrType:  18,
		AttrValue: []byte(msg),
	}
	replyMessage.Length()
	replyMessage.setStandardAttrStringVal()
	attr, _ := ATTRITUBES[AttrKey{Standard, int(replyMessage.AttrType)}]
	replyMessage.AttrName = attr.Name
	cxt.Response.AddRadiusAttr(replyMessage)
	cxt.Response.PackageLength()
	secret := cxt.RadNas.Secret
	authReplyAuthenticator(cxt.Request.Authenticator, cxt.Response, secret)
	cxt.Listener.WriteToUDP(cxt.Response.ToByte(), cxt.Dst)
}