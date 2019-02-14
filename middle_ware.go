package main

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"net"
)

func passwordMiddle(rp RadiusPackage)  {
	if rp.isChap {

	}
}

func authReply(rp RadiusPackage, listener *net.UDPConn, dest *net.UDPAddr) {
	reply := RadiusPackage {
		Code:ACCESS_ACCEPT_CODE,
		Identifier: rp.Identifier,
		Authenticator:[16]byte{},
	}

	replyMessage := RadiusAttr{
		AttrType: 18,
		AttrValue: []byte("authenticator success"),
	}

	replyMessage.Length()
	reply.AddRadiusAttr(replyMessage)
	reply.PackageLength()

	// TODO secret
	reply.Authenticator = authReplyAuthenticator(rp.Authenticator, reply, "111111", []RadiusAttr{replyMessage})

	listener.WriteToUDP(reply.ToByte(), dest)
}

// ResponseAuth = MD5(Code+ID+Length+RequestAuth+Attributes+Secret)
func authReplyAuthenticator(authAuthenticator [16]byte, reply RadiusPackage, secret string, replyAttrs []RadiusAttr) [16]byte {
	md5hash := md5.New()
	var buf bytes.Buffer
	buf.WriteByte(reply.Code)
	buf.WriteByte(reply.Identifier)

	var lengthBytes = make([]byte, 2)
	binary.BigEndian.PutUint16(lengthBytes, reply.Length)
	buf.Write(lengthBytes)

	buf.Write(authAuthenticator[:])

	for _, attr := range replyAttrs {
		if attr.AttrLength == 0 {
			continue
		}
		buf.Write(attr.toBytes())
	}

	buf.Write([]byte(secret))

	md5hash.Write(buf.Bytes())
	sum := md5hash.Sum(nil)

	return getSixteenBytes(sum)
}