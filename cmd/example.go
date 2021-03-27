package main

import (
	"bilicoin"
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
	"time"
)

type cmdOptions struct {
	Help   bool `short:"h" long:"help" description:"show this help message"`
	Start  bool `short:"s" long:"start" description:"start bilicoin service"`
	Delete bool `short:"d" long:"delete" description:"delete user info"`
	New    bool `short:"n" long:"new" description:"create user info by QRCode"`
	List   bool `short:"l" long:"list" description:"show all users"`
	FT     bool `short:"f" long:"ft" description:"set ftqq secret key"`
	Cron   bool `short:"c" long:"cron" description:"update cron spec"`
	API    bool `short:"a" long:"api" description:"run api server"`
}

func showHelp() {
	println()
	const v = `Usage: bilicoin [OPTIONS] [TEXT]

Options:
[-h] Show this help message
[-s] Start bilicoin service
[-d] Try to delete provided user information
     eg: bilicoin -d 10023442
[-n] Create user info by QRCode
[-l] Show all users
[-f] Try to set and open ftqq service for a provided user
     eg: bilicoin -f 10023442 SDF8S12J123AP2139CAI1
[-c] Try to update cron spec for a provided user
     eg: bilicoin -f 10023442 30 16 1 * * ?
[-a] Run api server
`
	fmt.Printf("\x1b[%dm"+v+" \x1b[0m", 34)
}

func main() {
	bilicoin.InitBili()

	release()
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

func release() {

	ret := 0
	defer func() { os.Exit(ret) }()

	opts := &cmdOptions{}
	optsParser := flags.NewParser(opts, flags.PrintErrors)
	args, err := optsParser.Parse()
	if err != nil || opts.Help || len(args) > 7 {
		ret = 1
		showHelp()
		return
	}

	if opts.List {
		// 用户列表
		users := bilicoin.LoadUser()
		fmt.Println("")
		fmt.Println("total: " + strconv.Itoa(len(users)))
		fmt.Println()
		fmt.Println("| UID | Cron | FTQQ |")
		for _, v := range users {
			print("| ")
			print(v.DedeUserID)
			print(" | ")
			print(v.Cron)
			print(" | ")
			print(v.FTSwitch)
			println(" | ")
		}
		fmt.Println("")
		return
	} else if opts.Cron {
		// Cron
		if len(args) != 7 {
			println("incorrect number of parameters")
			ret = 2
			return
		}
		var uid, cronStr string
		uid = args[0]
		cronStr = args[1] + " " + args[2] + " " + args[3] + " " + args[4] + " " + args[5] + " " + args[6]
		if biu, ok := bilicoin.GetUser(uid); ok == nil && biu != nil {
			if _, err = cron.Parse(cronStr); err != nil {
				bilicoin.Info("incorrect cron spec, please check and try again", logrus.Fields{"UID": uid, "Cron": cronStr})
				ret = 3
				return
			}
			biu.Cron = cronStr
			biu.InfoUpdate()
		}
		bilicoin.Info("cron save completed", logrus.Fields{"UID": uid, "Cron": cronStr})
	} else if opts.FT {
		// 方糖QQ
		if len(args) != 2 {
			println("incorrect number of parameters")
			ret = 2
			return
		}
		var uid, key string
		if len(args[1]) > len(args[0]) {
			uid = args[0]
			key = args[1]
		} else {
			uid = args[1]
			key = args[0]
		}
		if biu, ok := bilicoin.GetUser(uid); ok == nil && biu != nil {
			biu.FT = key
			biu.FTSwitch = true
			biu.InfoUpdate()
		}
		bilicoin.Info("ftqq secret save completed", logrus.Fields{"UID": uid, "Key": key})
	} else if opts.Delete {
		// 删除
		if len(args) != 1 {
			println("incorrect number of parameters")
			ret = 2
			return
		}
		bilicoin.Info("try to delete user", logrus.Fields{"UID": args[0]})
		_ = bilicoin.DelUser(args[0])
	} else if opts.New {
		// 新建
		if len(args) != 0 {
			println("incorrect number of parameters")
			ret = 2
			return
		}
		user, _ := bilicoin.CreateUser()
		_ = user.GetQRCode()
		user.QRCodePrint()
		user.BiliScanAwait()
		for true {
			if user.DedeUserID != "" {
				time.Sleep(time.Second * 5)
				os.Exit(0)
			}
		}
	} else if opts.Start {
		// 以普通模式运行
		bilicoin.AppInfo(bilicoin.Simple)
		bilicoin.CronTaskLoad()
		select {}
	} else if opts.API {
		// 以服务模式运行
		bilicoin.AppInfo(bilicoin.Api)
		bilicoin.BCApplication()

	} else {
		ret = 10
		showHelp()
		return
	}
}
