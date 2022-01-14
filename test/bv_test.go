package test

import (
	"bilicoin"
	"github.com/robertkrimen/otto"
	"io/ioutil"
	"testing"
)

// base58 enc
func TestBVCovertEnc(t *testing.T) {
	println(bilicoin.BVCovertDec("BV15V411z7Jk"))
	println(bilicoin.BVCovertDec("BV1Q541167Qg"))
	println(bilicoin.BVCovertDec("BV1mK4y1C7Bz"))
	println(bilicoin.BVCovertDec("BV15V411z7Jk"))
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

func TestPopular(t *testing.T) {
	a := bilicoin.GetPopulars()
	println(a[0].Title)
}

func TestUUID(t *testing.T) {
	bytes, err := ioutil.ReadFile("bili.js")
	vm := otto.New()
	_, err = vm.Run(bytes)
	uuid, err := vm.Call("generateUuid", nil)
	println(uuid.String())
	println(err)
}

var inlineJSCode = `function generateUuidPart(c){for(var a="",b=0;b<c;b++){a+=parseInt(16*Math.random()).toString(16).toUpperCase()}return formatNum(a,c)}function generateUuid(){var d=generateUuidPart(8),b=generateUuidPart(4),c=generateUuidPart(4),g=generateUuidPart(4),f=generateUuidPart(12),a=(new Date).getTime();return d+"-"+b+"-"+c+"-"+g+"-"+f+formatNum((a%100000).toString(),5)+"infoc"}function formatNum(c,a){var b="";if(c.length<a){for(var d=0;d<a-c.length;d++){b+="0"}}return b+c};`

func TestUUIDInline(t *testing.T) {
	vm := otto.New()
	_, err := vm.Run(inlineJSCode)
	uuid, err := vm.Call("generateUuid", nil)
	println(uuid.String())
	println(err)
}
