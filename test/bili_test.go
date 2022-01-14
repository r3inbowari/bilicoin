package test

import (
	"bilicoin"
	"fmt"
	"io/ioutil"
	"regexp"
	"testing"
	"time"
)

func TestBili(t *testing.T) {
	bilicoin.InitBili("dev", "v1.0.0", "cb0dc838e04e841f193f383e06e9d25a534c5809")

	// release()
	bilicoin.CronTaskLoad()
	time.Sleep(time.Hour)
}

// 获取投币日志(数量)
func TestGetLog(t *testing.T) {
	if len(bilicoin.GetConfig(false).BiU) > 0 {
		bilicoin.GetConfig(false).BiU[0].GetBiliCoinLog()
		println("size -> ", bilicoin.GetConfig(false).BiU[0].DropCoinCount)
	} else {
		println("please add a user to test")
	}
}

// 获取钱包数据
func TestGetWallet(t *testing.T) {
	if len(bilicoin.GetConfig(false).BiU) > 0 {
		bilicoin.GetConfig(false).BiU[0].GetBiliWallet()
		//println("size -> ", bilicoin.GetConfig(false).BiU[0].DropCoinCount)
	} else {
		println("please add a user to test")
	}
}

// 获取钱包数据
func TestConvert2Coin(t *testing.T) {
	if len(bilicoin.GetConfig(false).BiU) > 0 {
		bilicoin.GetConfig(false).BiU[0].Silver2Coin()
		//println("size -> ", bilicoin.GetConfig(false).BiU[0].DropCoinCount)
	} else {
		println("please add a user to test")
	}
}

// 投币测试
func TestDrop(t *testing.T) {
	if len(bilicoin.GetConfig(false).BiU) > 0 {
		//bilicoin.GetConfig().BiU[0].GetBiliCoinLog()
		//println("size -> ", bilicoin.GetConfig().BiU[0].DropCoinCount)
		bilicoin.GetConfig(false).BiU[0].DropCoin(bilicoin.Video{
			Aid:    123123123,
			Bvid:   "BV1sQ4y1X7WK",
			Author: "",
			Coins:  0,
			Title:  "",
		})
	} else {
		println("please add a user to test")
	}
}

// 获取BVs
func TestGetBVs(t *testing.T) {
	res, _ := bilicoin.GET("https://api.bilibili.com/x/web-interface/ranking/region?rid=119&day=3&original=0", nil)

	result, _ := ioutil.ReadAll(res.Body)

	reg := regexp.MustCompile("BV[a-zA-Z0-9_]+")
	g := reg.FindAllString(string(result), -1)
	for _, v := range g {
		fmt.Println(v)
	}
}

func TestCopy(t *testing.T) {
	a := []int{0}
	i := 0
	a = append(a[:i], a[i+1:]...)
	println("!1111")
}

func MainT() {
	// example:
	// add
	// bilicoin.InitLogger()
	// bilicoin.Info("Canvas Fingerprinting " + bilicoin.GetConfig().Finger)
	// user, _ := bilicoin.CreateUser()
	// user.GetQRCode()
	// user.QRCodePrint()
	// user.BiliScanAwait()
	// del
	// _ = bilicoin.DelUser("30722")
	// drop
	// biu, _ := bilicoin.GetUser("30722")
	// biu.RandDrop()
	// time.Sleep(time.Hour)
}
