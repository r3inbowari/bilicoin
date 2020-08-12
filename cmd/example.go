package main

import (
	"bilicoin"
	"github.com/jessevdk/go-flags"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

type cmdOptions struct {
	Help   bool `short:"h" long:"help" description:"show this help message"`
	Start  bool `short:"s" long:"start" description:"start bilicoin service"`
	Delete bool `short:"d" long:"delete" description:"delete user info"`
	New    bool `short:"n" long:"new" description:"create user info by QRCode"`
	List   bool `short:"l" long:"list" description:"show all users"`
	FT     bool `short:"f" long:"ft" description:"set ftqq secret key"`
}

func showHelp() {
	const v = `Usage: bilicoin [OPTIONS] [TEXT]

Options:
  -h, --help
    Show this help message
  -s, --start
    Start bilicoin service
  -d, --delete
    Delete user info
  -n, --new
    Create user info by QRCode
  -l, --list
    Show all users
  -f, --ft
    Set ftqq secret key
`
	os.Stderr.Write([]byte(v))
}

func main() {
	bilicoin.InitConfig()
	bilicoin.InitLogger()

	ret := 0
	defer func() { os.Exit(ret) }()

	opts := &cmdOptions{}
	optsParser := flags.NewParser(opts, flags.PrintErrors)
	args, err := optsParser.Parse()
	if err != nil || len(args) > 1 {
		showHelp()
		ret = 1
		return
	}

	var text string
	if len(args) == 1 {
		text = args[0]
	}

	if opts.Help {
		showHelp()
		return
	} else if opts.List {
		bius := bilicoin.GetAllUID()
		println("total:", len(bius))
		for _, v := range bius {
			println(v)
		}
		return
	} else if opts.FT {
		c := bilicoin.GetConfig()
		c.FT = text
		_ = c.SetConfig()
		bilicoin.Info("ftqq secret save completed")
		bilicoin.Info("current key:" + text)
	} else if opts.Delete {
		_ = bilicoin.DelUser(text)
		bilicoin.Info("try to delete user", logrus.Fields{"UID": text})
	} else if opts.New {
		user, _ := bilicoin.CreateUser()
		user.GetQRCode()
		user.QRCodePrint()
		user.BiliScanAwait()
		for true {
			if user.DedeUserID != "" {
				time.Sleep(time.Second * 5)
				os.Exit(0)
			}
		}
	} else if opts.Start {
		println("bilicoin " + bilicoin.Version)
		bilicoin.CronDropReg()
		select {}
	} else {
		ret = 2
		showHelp()
		return
	}

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
