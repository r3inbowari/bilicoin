package test

import (
	"fmt"
	"math/rand"
	"regexp"
	"testing"
)

func TestCoinVerify(t *testing.T) {
	reg := regexp.MustCompile("BV[a-zA-Z0-9_]+")
	str := "给视频 BV16h411d7B8 打赏"
	g := reg.FindAllString(str, -1)
	fmt.Println(g)
}

func TestRand(t *testing.T) {

	println(rand.Intn(10))
	println(rand.Intn(10))
	println(rand.Intn(10))
	println(rand.Intn(10))
	println(rand.Intn(10))
	println(rand.Intn(10))
	println(rand.Intn(10))
	println(rand.Intn(10))
	println(rand.Intn(10))
	println(rand.Intn(10))
	println(rand.Intn(10))
	println(rand.Intn(10))
	println(rand.Intn(10))
	println(rand.Intn(10))
	println(rand.Intn(10))
	println(rand.Intn(10))
	println(rand.Intn(10))
	println(rand.Intn(10))
	println(rand.Intn(10))
	println(rand.Intn(10))
	println(rand.Intn(10))
	println(rand.Intn(10))
	println(rand.Intn(10))
	println(rand.Intn(10))
	println(rand.Intn(10))
	println(rand.Intn(10))
	println(rand.Intn(10))
	println(rand.Intn(10))
	println(rand.Intn(10))
	println(rand.Intn(10))

}
