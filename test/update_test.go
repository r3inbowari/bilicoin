package test

import (
	"bilicoin"
	"runtime"
	"testing"
	"time"
)

func TestWrite(t *testing.T) {
	bilicoin.CreateInstallBatch("cmd1")
}

func TestExec(t *testing.T) {

	bilicoin.CheckAndUpdateAndReload()

	time.Sleep(time.Second * 20)
}

func TestDigest(t *testing.T) {
	if runtime.GOOS == "windows" {
		println(bilicoin.DigestVerify("../LICENSE", "a1cb229ccb06000ce4fa2fd9ee6764bf"))
	}
}

func TestGetDigest(t *testing.T) {
	bilicoin.CheckUpdate()
}