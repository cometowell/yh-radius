package main

import "net"

func passwordMiddle(rp RadiusPackage)  {
	if rp.isChap {

	}
}

func authReply(rp RadiusPackage, listener *net.UDPConn, dest *net.UDPAddr) {
	reply := RadiusPackage {
		Code:ACCESS_ACCEPT_CODE,
		Identifier: rp.Identifier,
		Authenticator: rp.Authenticator,
	}

	replyMessage := RadiusAttr{
		AttrType: 18,
		AttrValue: []byte("认证成功"),
	}

	replyMessage.Length()
	reply.AddRadiusAttr(replyMessage)
	reply.PackageLength()

	listener.WriteToUDP(reply.ToByte(), dest)
}
