package main

func authMiddleWare(rp RadiusPackage) {
	attrs := rp.RadiusAttrs

	for _, attr := range attrs {
		if attr.AttrName == "User-Password" {

		}

		if attr.AttrName == "CHAP-Password" {

		}
	}
}

func authReply()  {

}

func accountingMiddleWare(rp RadiusPackage) {

}

func accountingReply() {

}