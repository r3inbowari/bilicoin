package main

import (
	"bilicoin"
	"fmt"
)

func main() {

	fmt.Println("Canvas Fingerprinting")
	user, _ := bilicoin.CreateUser()
	user.GetQRCode()
	user.QRCodePrint()
	user.BiliScanAwait()

	select {}
}


