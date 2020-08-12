package test

import (
	"bilicoin"
	"fmt"
	"io/ioutil"
	"regexp"
	"testing"
)

func TestBiliHomepage(t *testing.T) {
	res, _ := bilicoin.GET("https://api.bilibili.com/x/web-interface/ranking/region?rid=119&day=3&original=0", nil)

	result, _ := ioutil.ReadAll(res.Body)

	reg := regexp.MustCompile("BV[a-zA-Z0-9_]+")
	g := reg.FindAllString(string(result), -1)
	fmt.Println(g)
}

func TestCopy(t *testing.T) {
	a := []int{0}
	i := 0
	a = append(a[:i], a[i+1:]...)
	println("!1111")
}
