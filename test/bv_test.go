package test

import (
	"bilicoin"
	"testing"
)

// base58 enc
func TestBVCovertEnc(t *testing.T) {
	println(bilicoin.BVCovertDec("BV15V411z7Jk"))
	println(bilicoin.BVCovertDec("BV1Q541167Qg"))
	println(bilicoin.BVCovertDec("BV1mK4y1C7Bz"))
	println(bilicoin.BVCovertDec("BV15V411z7Jk"))

	println(bilicoin.BVCovertDec("BV1JD4y1U7"))
}

// base58 dec
func TestBVCovertDec(t *testing.T) {
	println(bilicoin.BVCovertEnc("170001"))
	println(bilicoin.BVCovertEnc("455017605"))
	println(bilicoin.BVCovertEnc("882584971"))
}

func TestGetGuichu(t *testing.T) {
	bvs := bilicoin.GetGuichuBVs()
	for _, v := range bvs {
		println(v)
	}
}
