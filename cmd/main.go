package main

import (
	"bilicoin"
	"github.com/jessevdk/go-flags"
	. "github.com/r3inbowari/zlog"
	"github.com/r3inbowari/zupdate"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
	"os"
	"runtime"
	"strconv"
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
	Log.Blue(v)
}

var (
	GitHash        = "cb0dc838e04e841f193f383e06e9d25a534c5809"
	buildTime      = "Thu Oct 01 00:00:00 1970 +0800"
	goVersion      = runtime.Version()
	ReleaseVersion = "ver[DEV]"
	Major          string
	Minor          string
	Patch          string
)

var Mode = "DEV"

func main() {
	InitGlobalLogger()
	Log.SetScreen(true)

	bilicoin.InitBili(Mode, ReleaseVersion, GitHash, Major, Minor, Patch)

	_ = zupdate.InitUpdater(zupdate.Option{
		Name:        "hola",
		Mode:        zupdate.DEV,
		CheckSource: "http://120.77.33.188:3000/bilicoin/bin/test.json",
		Callback:    nil,
		EntryArgs:   []string{"-a"},
		Log:         &Log.Logger,
	}, nil)

	//up.IncludeFile("5c74bf9c1face2dcb47bae100f2c664cdbd43400", zupdate.File{
	//	Major: 1,
	//	Minor: 0,
	//	Patch: 1,
	//})
	//up.CheckAndUpdateWithGap()

	release()
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
		bilicoin.UserList()
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
				Log.WithFields(logrus.Fields{"UID": uid, "Cron": cronStr}).Info("incorrect cron spec, please check and try again")
				ret = 3
				return
			}
			biu.Cron = cronStr
			biu.InfoUpdate()
		}
		Log.WithFields(logrus.Fields{"UID": uid, "Cron": cronStr}).Info("cron save completed")
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
		Log.WithFields(logrus.Fields{"UID": uid, "Key": key}).Info("ftqq secret save completed")
	} else if opts.Delete {
		// 删除
		if len(args) != 1 {
			println("incorrect number of parameters")
			ret = 2
			return
		}
		Log.WithFields(logrus.Fields{"UID": args[0]}).Info("try to delete user")
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
	} else if opts.Start {
		// 以普通模式运行
		bilicoin.RunningMode = bilicoin.Simple
		AppInfo(GitHash, buildTime, goVersion, ReleaseVersion, "normal")
		bilicoin.CronTaskLoad()
		select {}
	} else if opts.API {
		// 以服务模式运行
		bilicoin.RunningMode = bilicoin.Api
		AppInfo(GitHash, buildTime, goVersion, ReleaseVersion, "server")
		bilicoin.BCApplication()

	} else {
		ret = 10
		showHelp()
		return
	}
}

func AppInfo(gitHash, buildTime, goVersion string, version string, mode string) {
	Log.Blue("  ________  ___  ___       ___  ________  ________  ___  ________")
	Log.Blue(" |\\   __  \\|\\  \\|\\  \\     |\\  \\|\\   ____\\|\\   __  \\|\\  \\|\\   ___  \\         BILICOIN #UNOFFICIAL# " + gitHash[:7] + "..." + gitHash[33:])
	Log.Blue(" \\ \\  \\|\\ /\\ \\  \\ \\  \\    \\ \\  \\ \\  \\___|\\ \\  \\|\\  \\ \\  \\ \\  \\\\ \\  \\        -... .. .-.. .. -.-. --- .. -. " + version)
	Log.Blue("  \\ \\   __  \\ \\  \\ \\  \\    \\ \\  \\ \\  \\    \\ \\  \\\\\\  \\ \\  \\ \\  \\\\ \\  \\       Running mode: " + mode)
	if mode == "server" {
		Log.Blue("   \\ \\  \\|\\  \\ \\  \\ \\  \\____\\ \\  \\ \\  \\____\\ \\  \\\\\\  \\ \\  \\ \\  \\\\ \\  \\      Port: " + bilicoin.GetConfig(false).APIAddr[1:])
	} else {
		Log.Blue("   \\ \\  \\|\\  \\ \\  \\ \\  \\____\\ \\  \\ \\  \\____\\ \\  \\\\\\  \\ \\  \\ \\  \\\\ \\  \\      Port: UNSUPPORTED")
	}
	Log.Blue("    \\ \\_______\\ \\__\\ \\_______\\ \\__\\ \\_______\\ \\_______\\ \\__\\ \\__\\\\ \\__\\     PID: " + strconv.Itoa(os.Getpid()))
	Log.Blue("     \\|_______|\\|__|\\|_______|\\|__|\\|_______|\\|_______|\\|__|\\|__| \\|__|     built at " + buildTime)
	Log.Blue("")
}
