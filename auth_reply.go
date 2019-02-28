package main

func authReply(cxt *Context, replyCode byte, msg string) {

	cxt.Response = &RadiusPackage{
		Code:          replyCode,
		Identifier:    cxt.Request.Identifier,
		Authenticator: [16]byte{},
	}

	replyMessage := RadiusAttr{
		AttrType:  18,
		AttrValue: []byte(msg),
	}
	replyMessage.Length()
	replyMessage.setStandardAttrStringVal()
	attr, _ := ATTRITUBES[AttrKey{0, int(replyMessage.AttrType)}]
	replyMessage.AttrName = attr.Name
	cxt.Response.AddRadiusAttr(replyMessage)

	for _, attr := range cxt.ReplyRadiusAttrs {
		attr.Length()
		cxt.Response.AddRadiusAttr(attr)
	}

	cxt.Response.PackageLength()
	secret := cxt.RadNas.Secret
	authReplyAuthenticator(cxt.Request.Authenticator, cxt.Response, secret)
	cxt.Listener.WriteToUDP(cxt.Response.ToByte(), cxt.Dst)
}