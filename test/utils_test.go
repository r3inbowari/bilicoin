package test

import (
	"bilicoin"
	"testing"
)

func TestInit(t *testing.T) {
	bilicoin.InitConfig()
}

func TestFTQQ(t *testing.T) {
	bilicoin.SendMessage2WeChat("hello1", "hello1")
	bilicoin.SendMessage2WeChat("hello2")
}

func TestRandom(t *testing.T) {
	println(bilicoin.Random(123))
}