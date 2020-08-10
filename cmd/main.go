package main

import (
	"bilicoin"
	"fmt"
)

type cmdOptions struct {
	Help    bool `short:"h" long:"help" description:"show this help message"`
	Inverse bool `short:"i" long:"invert" description:"invert color"`
}

func main() {

	fmt.Println("帆布指纹识别")


	bilicoin.GetQRCode()

}
