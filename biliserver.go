package bilicoin

import (
	"github.com/gorilla/mux"
	. "github.com/r3inbowari/zlog"
	"github.com/r3inbowari/zserver"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
	"net/http"
	"sync"
	"time"
)

func BCApplication() {
	Log.Info("[BCS] BiliCoin is running")
	reset()

	d := zserver.DefaultServer(zserver.Options{
		Log:    &Log.Logger,
		Addr:   GetConfig(false).APIAddr,
		Mode:   BuildMode,
		CaCert: GetConfig(false).CaCert,
		CaKey:  GetConfig(false).CaKey,
	})

	d.Map("/{uid}/ft", HandleFT)
	d.Map("/{uid}/cron", HandleCron)
	d.Map("/version", HandleVersion)
	d.Map("/users", HandleUsers)
	d.Map("/user", HandleUserAdd, http.MethodPost)
	d.Map("/user", HandleUserDel, http.MethodGet)
	d.Start()
}

type FilterBiliUser struct {
	UID       string `json:"uid"`
	Cron      string `json:"cron"`
	FT        string `json:"ft"`
	FTSwitch  bool   `json:"ftSwitch"`
	DropCount int    `json:"dropCount"`
}

func HandleFT(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	vars := mux.Vars(r)
	uid := vars["uid"]
	key := r.Form.Get("key")
	sw := r.Form.Get("sw")

	if biu, ok := GetUser(uid); ok == nil && biu != nil {
		if key != "" {
			biu.FT = key
			biu.FTSwitch = true
		}
		if sw == "0" {
			biu.FTSwitch = false
		} else {
			biu.FTSwitch = true
		}
		biu.InfoUpdate()
		Log.WithFields(logrus.Fields{"UID": uid, "Key": key}).Info("[BCS] FTQQ secret key save completed")
	}
	zserver.ResponseCommon(w, "try succeed", "ok", 1, http.StatusOK, 0)
}

func HandleCron(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	vars := mux.Vars(r)
	uid := vars["uid"]
	cronStr := r.Form.Get("spec")

	if biu, ok := GetUser(uid); ok == nil && biu != nil {
		if _, err := cron.Parse(cronStr); err != nil {
			zserver.ResponseCommon(w, "[BCS] incorrect cron spec, please check and try again", "ok", 1, http.StatusOK, 0)
			Log.WithFields(logrus.Fields{"UID": uid, "Cron": cronStr}).Info("[BCS] incorrect cron spec, please check and try again")
			return
		}
		biu.Cron = cronStr
		biu.InfoUpdate()
		Log.WithFields(logrus.Fields{"UID": uid, "Cron": cronStr}).Info("[BCS] Cron save completed by web")
	}
	reset()
	zserver.ResponseCommon(w, "try succeed", "ok", 1, http.StatusOK, 0)
}

func HandleVersion(w http.ResponseWriter, r *http.Request) {
	zserver.ResponseCommon(w, releaseVersion+" "+releaseTag, "ok", 1, http.StatusOK, 0)
}

func HandleUsers(w http.ResponseWriter, r *http.Request) {
	users := LoadUsers()
	var ret []FilterBiliUser
	for _, v := range users {
		ret = append(ret, FilterBiliUser{
			UID:       v.DedeUserID,
			Cron:      v.Cron,
			FT:        v.FT,
			FTSwitch:  v.FTSwitch,
			DropCount: v.DropCoinCount,
		})
	}
	zserver.ResponseCommon(w, ret, "ok", len(users), http.StatusOK, 0)
}

var loginMap sync.Map

func HandleUserAdd(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	if r.Form.Get("oauth") == "" {
		// 提供回调
		user, _ := CreateUser()
		_ = user.GetQRCode()
		Log.WithFields(logrus.Fields{"oauth": user.OAuth.OAuthKey}).Info("[BCS] qrcode created")
		loginMap.Store(user.OAuth.OAuthKey, user)
		time.AfterFunc(time.Minute*3, func() {
			loginMap.Delete(user.OAuth.OAuthKey)
		})
		zserver.ResponseCommon(w, user.OAuth.OAuthKey, "ok", 1, http.StatusOK, 0)
	} else {
		if user, ok := loginMap.Load(r.Form.Get("oauth")); ok {
			biliUser := user.(*BiliUser)
			biliUser.GetBiliLoginInfo(nil)
			if biliUser.Login {
				zserver.ResponseCommon(w, biliUser.Login, "ok", 1, http.StatusOK, 0)
				reset()
			}
		} else {
			zserver.ResponseCommon(w, "non exist", "ok", 1, http.StatusOK, 0)
		}
	}
}

func HandleUserDel(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	uid := r.Form.Get("uid")
	Log.WithFields(logrus.Fields{"UID": uid}).Info("[BCS] Try to delete user")
	_ = DelUser(uid)
	reset()
	zserver.ResponseCommon(w, "try succeed", "ok", 1, http.StatusOK, 0)
}

func reset() {
	cronTask.Range(func(key, value interface{}) bool {
		Log.WithFields(logrus.Fields{"TaskID": key.(string)}).Info("[BSC] release task")
		value.(*cron.Cron).Stop()
		cronTask.Delete(key)
		return true
	})
	CronTaskLoad()
}
