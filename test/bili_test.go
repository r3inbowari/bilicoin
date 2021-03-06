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
	bilicoin.InitBili()

	// release()
	bilicoin.CronTaskLoad()
	time.Sleep(time.Hour)
}

// 获取投币日志(数量)
func TestGetLog(t *testing.T) {
	if len(bilicoin.GetConfig().BiU) > 0 {
		bilicoin.GetConfig().BiU[0].GetBiliCoinLog()
		println("size -> ", bilicoin.GetConfig().BiU[0].DropCoinCount)
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
