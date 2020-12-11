package test

import (
	"bilicoin"
	"testing"
)

func TestInit(t *testing.T) {
	bilicoin.InitConfig()
}

func TestFTQQ(t *testing.T) {
	if len(bilicoin.GetConfig().BiU) > 0 {
		bilicoin.GetConfig().BiU[0].SendMessage2WeChat("hello1", "hello1")
		bilicoin.GetConfig().BiU[0].SendMessage2WeChat("hello2")
	} else {
		println("please add a user to test")
	}

}

func TestRandom(t *testing.T) {
	println(bilicoin.Random(123))
}
