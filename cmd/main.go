package main

import (
	"bilicoin"
	"fmt"
	. "github.com/r3inbowari/zlog"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"runtime"
	"strconv"
	"strings"
)

var (
	run     = kingpin.Flag("run", "Run bilicoin service").Short('a').Bool()
	list    = kingpin.Flag("list", "Show all users").Short('l').Bool()
	login   = kingpin.Flag("new", "login by QRCode").Short('n').Bool()
	del     = kingpin.Flag("delete", "Try to delete provided user by uid. eg: bilicoin -d 10023442").Short('d').String()
	uid     = kingpin.Flag("uid", "uid").Short('u').String()
	ft      = kingpin.Flag("ft", "Try to set and open ftqq service for a provided uid. eg: bilicoin -u 10023442 -f SDF8S12J123AP2139CAI1").Short('f').String()
	cronStr = kingpin.Flag("cron", "Try to update cron spec for a provided uid. eg: bilicoin -u 10023442 -f 30,16,1,*,*,?").Short('c').String()
)

// injected params
var (
	GitHash        = "cb0dc838e04e841f193f383e06e9d25a534c5809"
	buildTime      = "Thu Oct 01 00:00:00 1970 +0800"
	goVersion      = runtime.Version()
	ReleaseVersion = "ver[DEV]"
	Major          string
	Minor          string
	Patch          string
	Mode           = "DEV"
)

func main() {
	InitGlobalLogger()
	Log.SetScreen(true)
	bilicoin.InitBili(Mode, ReleaseVersion, GitHash)
	release()
}

func release() {

	kingpin.Version(fmt.Sprintf("%s git-%s", ReleaseVersion, GitHash))
	kingpin.Parse()

	if *run {
		// 运行
		AppInfo(GitHash, buildTime, goVersion, ReleaseVersion, "server")
		bilicoin.BCApplication()
		return
	}

	if *login {
		// 登录
		user, _ := bilicoin.CreateUser()
		_ = user.GetQRCode()
		user.QRCodePrint()
		user.BiliScanAwait()
		return
	}

	if *list {
		// 用户列表
		bilicoin.UserList()
		return
	}

	if *del != "" {
		// 删除
		Log.WithFields(logrus.Fields{"UID": *del}).Info("try to delete user")
		_ = bilicoin.DelUser(*del)
		return
	}

	if *ft != "" && *uid != "" {
		if biu, ok := bilicoin.GetUser(*uid); ok == nil && biu != nil {
			biu.FT = *ft
			biu.FTSwitch = true
			biu.InfoUpdate()
		}
		Log.WithFields(logrus.Fields{"UID": uid, "Key": *ft}).Info("ftqq secret save completed")
		return
	}

	if *cronStr != "" && *uid != "" {
		spec := strings.ReplaceAll(*cronStr, ",", " ")
		if _, err := cron.Parse(spec); err != nil {
			Log.WithFields(logrus.Fields{"UID": uid, "Cron": cronStr}).Error("incorrect cron spec, please check and try again")
			return
		}

		if biu, ok := bilicoin.GetUser(*uid); ok == nil && biu != nil {
			biu.Cron = spec
			biu.InfoUpdate()
			Log.WithFields(logrus.Fields{"UID": uid, "Cron": spec}).Info("cron save completed")
			return
		}
	}
}

func AppInfo(gitHash, buildTime, goVersion string, version string, mode string) {
	Log.Blue("  ________  ___  ___       ___  ________  ________  ___  ________")
	Log.Blue(" |\\   __  \\|\\  \\|\\  \\     |\\  \\|\\   ____\\|\\   __  \\|\\  \\|\\   ___  \\         BILICOIN #UNOFFICIAL# " + gitHash[:7] + "..." + gitHash[33:])
	Log.Blue(" \\ \\  \\|\\ /\\ \\  \\ \\  \\    \\ \\  \\ \\  \\___|\\ \\  \\|\\  \\ \\  \\ \\  \\\\ \\  \\        -... .. .-.. .. -.-. --- .. -. " + version)
	Log.Blue("  \\ \\   __  \\ \\  \\ \\  \\    \\ \\  \\ \\  \\    \\ \\  \\\\\\  \\ \\  \\ \\  \\\\ \\  \\       Running mode: " + mode)
	Log.Blue("   \\ \\  \\|\\  \\ \\  \\ \\  \\____\\ \\  \\ \\  \\____\\ \\  \\\\\\  \\ \\  \\ \\  \\\\ \\  \\      Port: " + bilicoin.GetConfig(false).APIAddr[1:])
	Log.Blue("    \\ \\_______\\ \\__\\ \\_______\\ \\__\\ \\_______\\ \\_______\\ \\__\\ \\__\\\\ \\__\\     PID: " + strconv.Itoa(os.Getpid()))
	Log.Blue("     \\|_______|\\|__|\\|_______|\\|__|\\|_______|\\|_______|\\|__|\\|__| \\|__|     built at " + buildTime)
	Log.Blue("")
}
