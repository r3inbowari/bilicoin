package bilicoin

import (
	"github.com/gorilla/mux"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
	"github.com/wuwenbao/gcors"
	"net/http"
	"sync"
	"time"
)

func BCApplication() {
	Info("[BCS] BILICOIN api mode running")
	reset()
	Info("[BCS] Listened on " + GetConfig().APIAddr)
	r := mux.NewRouter()

	r.HandleFunc("/{uid}/ft", HandleFT)
	r.HandleFunc("/{uid}/cron", HandleCron)
	r.HandleFunc("/version", HandleVersion)
	r.HandleFunc("/users", HandleUsers)
	r.HandleFunc("/user", HandleUserAdd).Methods("POST")
	r.HandleFunc("/user", HandleUserDel).Methods("DELETE")

	// allow CORS
	cors := gcors.New(
		r,
		gcors.WithOrigin("*"),
		gcors.WithMethods("POST, GET, PUT, DELETE, OPTIONS"),
		gcors.WithHeaders("Authorization"),
	)
	log.Fatal(http.ListenAndServe(GetConfig().APIAddr, cors))

	time.Sleep(time.Hour)
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
		Info("[BCS] FTQQ secret key save completed", logrus.Fields{"UID": uid, "Key": key})
	}
	ResponseCommon(w, "try succeed", "ok", 1, http.StatusOK, 0)
}

func HandleCron(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	vars := mux.Vars(r)
	uid := vars["uid"]
	cronStr := r.Form.Get("spec")

	if biu, ok := GetUser(uid); ok == nil && biu != nil {
		if _, err := cron.Parse(cronStr); err != nil {
			ResponseCommon(w, "[BCS] incorrect cron spec, please check and try again", "ok", 1, http.StatusOK, 0)
			Info("[BCS] incorrect cron spec, please check and try again", logrus.Fields{"UID": uid, "Cron": cronStr})
			return
		}
		biu.Cron = cronStr
		biu.InfoUpdate()
		Info("[BCS] Cron save completed by web", logrus.Fields{"UID": uid, "Cron": cronStr})
	}
	reset()
	ResponseCommon(w, "try succeed", "ok", 1, http.StatusOK, 0)
}

func HandleVersion(w http.ResponseWriter, r *http.Request) {
	ResponseCommon(w, releaseVersion+" "+releaseTag, "ok", 1, http.StatusOK, 0)
}

func HandleUsers(w http.ResponseWriter, r *http.Request) {
	users := LoadUser()
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
	ResponseCommon(w, ret, "ok", len(users), http.StatusOK, 0)
}

var loginMap sync.Map

func HandleUserAdd(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	if r.Form.Get("oauth") == "" {
		// 提供回调
		user, _ := CreateUser()
		_ = user.GetQRCode()
		println(user.OAuth.OAuthKey)
		println(user.OAuth.Url)
		loginMap.Store(user.OAuth.OAuthKey, user)
		time.AfterFunc(time.Minute*3, func() {
			loginMap.Delete(user.OAuth.OAuthKey)
		})
		ResponseCommon(w, user.OAuth.OAuthKey, "ok", 1, http.StatusOK, 0)
	} else {
		if user, ok := loginMap.Load(r.Form.Get("oauth")); ok {
			biliUser := user.(*BiliUser)
			biliUser.LoginCallback(func(isLogin bool) {
				ResponseCommon(w, isLogin, "ok", 1, http.StatusOK, 0)
				if isLogin {
					// reset
					reset()
				}
			})
			// ResponseCommon(w, oauth, "ok", 1, http.StatusOK, 0)
		}
	}
}

func HandleUserDel(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	uid := r.Form.Get("uid")
	Info("[BCS] Try to delete user", logrus.Fields{"UID": uid})
	_ = DelUser(uid)
	reset()
	ResponseCommon(w, "try succeed", "ok", 1, http.StatusOK, 0)
}

func reset() {
	Warn("[BSC] Release task resource")
	taskMap.Range(func(key, value interface{}) bool {
		Info("[BSC] Try to stop cron", logrus.Fields{"TaskID": key.(string)})
		value.(*cron.Cron).Stop()
		taskMap.Delete(key)
		return true
	})
	CronTaskLoad()
}
