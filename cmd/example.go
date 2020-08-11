package main

import (
	"bilicoin"
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

type cmdOptions struct {
	Help   bool `short:"h" long:"help" description:"show this help message"`
	Start  bool `short:"s" long:"start" description:"start bilibili drop coin server"`
	Delete bool `short:"d" long:"delete" description:"delete user info"`
	New    bool `short:"n" long:"new" description:"create user info by QRCode"`
	List   bool `short:"l" long:"list" description:"show all users"`
}

func pErr(format string, a ...interface{}) {
	fmt.Fprint(os.Stdout, os.Args[0], ": ")
	fmt.Fprintf(os.Stdout, format, a...)
}

func showHelp() {
	const v = `Usage: qrc [OPTIONS] [TEXT]
Options:
  -h, --help
    Show this help message
  -i, --invert
    Invert color
Text examples:
  http://www.example.jp/
  MAILTO:foobar@example.jp
  WIFI:S:myssid;T:WPA;P:pass123;;
`

	os.Stderr.Write([]byte(v))
}

func main() {
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
	if opts.Help {
		showHelp()
		return
	}

	if opts.List {
		bius := bilicoin.GetAllUID()
		println("total:", len(bius))
		for _, v := range bius {
			println(v)
		}
		return
	}

	var text string
	if len(args) == 1 {
		text = args[0]
	}

	if opts.Delete {
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
		bius := bilicoin.GetConfig().BiU
		if len(bius) == 0 {
			println("not found users")
			return
		}
		for k, _ := range bius {
			println("cron - add user " + bius[k].DedeUserID)
			CronDrop(bius[k])
		}
		select {}
	} else {
		ret = 2
		showHelp()
		return
	}
	// add
	//bilicoin.InitLogger()
	//bilicoin.Info("Canvas Fingerprinting " + bilicoin.GetConfig().Finger)
	//user, _ := bilicoin.CreateUser()
	//user.GetQRCode()
	//user.QRCodePrint()
	//user.BiliScanAwait()

	// del
	//_ = bilicoin.DelUser("30722")

	// drop
	// biu, _ := bilicoin.GetUser("30722")
	// biu.RandDrop()
	// time.Sleep(time.Hour)
}

func CronDrop(biu bilicoin.BiliUser) {
	c := cron.New()
	_ = c.AddFunc("30 31 0 * * ?", func() {
		biu.GetBiliCoinLog()
		for i := 0; i < 5; i++ {
			biu.RandDrop()
			time.Sleep(time.Minute)
			bilicoin.Info("cron finish", logrus.Fields{"UID": biu.DedeUserID})
		}
	})
	c.Start()
}
