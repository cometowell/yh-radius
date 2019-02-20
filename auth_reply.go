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
	cxt.Response.AddRadiusAttr(replyMessage)

	for _, attr := range cxt.RadiusAttrs {
		attr.Length()
		cxt.Response.AddRadiusAttr(attr)
	}

	cxt.Response.PackageLength()
	// TODO secret
	secret := "111111"
	authReplyAuthenticator(cxt.Request.Authenticator, cxt.Response, secret)
	cxt.Listener.WriteToUDP(cxt.Response.ToByte(), cxt.Dst)
}