package bilicoin

import (
	"encoding/json"
	"errors"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
	"math/rand"
	"net/http"
	"regexp"
	"time"
)

type BiliQRCode struct {
	Data      QRCodeData `json:"data"`
	Timestamp int64      `json:"ts"`
	Status    bool       `json:"status"`
	Code      int        `json:"code"`
}

type QRCodeData struct {
	Url      string `json:"url"`
	OAuthKey string `json:"oauthKey"`
}

type BiliUser struct {
	UUID          string     `json:"_uuid"`
	BuVID         string     `json:"buvid"`
	OAuth         QRCodeData `json:"oauth"`
	SID           string     `json:"sid"`
	BiliJCT       string     `json:"bili_jct"`
	SESSDATA      string     `json:"SESSDATA"`
	DedeUserID    string     `json:"uid"`
	DedeUserIDMD5 string     `json:"DedeUserID__ckMd5"`
	Bi            ABi        `json:"bi"`
	Login         bool       `json:"login"`
	Expire        time.Time  `json:"expire"`
	DropCoinCount int        `json:"drop_coin_count,omitempty"`
	BlockBVList   []string   `json:"block_bv_list"`
}

type ABi struct {
	BCoins    int       `json:"bCoins"` // B柯拉
	Coins     int       `json:"coins"`  // 硬币
	Nick      string    `json:"uname"`
	Face      string    `json:"face"`
	Status    string    `json:"userStatus"`
	LevelInfo LevelInfo `json:"level_info"`
}

type LevelInfo struct {
	CurrentLevel int `json:"current_level"`
	CurrentMin   int `json:"current_min"`
	CurrentExp   int `json:"current_exp"`
	NextExp      int `json:"next_exp"`
}

type BiliInfo struct {
	Status    string   `json:"status"`
	Data      BiliData `json:"data"`
	Message   string   `json:"message"`
	Timestamp int64    `json:"ts"`
	TTL       int      `json:"ttl"`
	Code      int      `json:"code"`
}

type BiliData struct {
	List []CoinLog `json:"list,omitempty"`
	Like string    `json:"like,omitempty"`
	BVID string    `json:"bvid,omitempty"`
}

type CoinLog struct {
	Time   string `json:"time"`
	Delta  int    `json:"delta"`
	Reason string `json:"reason"`
}

var Version = "v1.0.1 build on 08 12 2020"

// 创建用户
func CreateUser() (*BiliUser, error) {
	var biu BiliUser
	var err error
	if biu.UUID, err = _uuidGenerate(); err != nil {
		return nil, err
	}
	if biu.BuVID, err = _buvidGenerate(); err != nil {
		return nil, err
	}
	return &biu, err
}

// 打印二维码
func (biu *BiliUser) QRCodePrint() {
	QRCPrint(biu.OAuth.Url)
}

// 获取二维码
func (biu *BiliUser) GetQRCode() error {
	url := "https://passport.bilibili.com/qrcode/getLoginUrl"
	res, err := GET(url, func(reqPoint *http.Request) {
		cookie1 := &http.Cookie{Name: "_uuid", Value: biu.UUID}
		cookie3 := &http.Cookie{Name: "buvid3", Value: biu.BuVID}
		cookie4 := &http.Cookie{Name: "PVID", Value: "1"}
		reqPoint.AddCookie(cookie1)
		reqPoint.AddCookie(cookie3)
		reqPoint.AddCookie(cookie4)
	})

	if err != nil {
		return err
	}

	var qr BiliQRCode
	err = json.NewDecoder(res.Body).Decode(&qr)
	if err != nil {
		return err
	}

	biu.OAuth = qr.Data
	biu.SID = res.Cookies()[0].Value
	return nil
}

// 获取登录信息
func (biu *BiliUser) GetBiliLoginInfo(cron *cron.Cron) {
	url := "https://passport.bilibili.com/qrcode/getLoginInfo"

	data := "oauthKey=" + biu.OAuth.OAuthKey + "&gourl=" + "https%3A%2F%2Fwww.bilibili.com%2F"
	res, err := Post2(url, func(reqPoint *http.Request) {

		cookie1 := &http.Cookie{Name: "_uuid", Value: biu.UUID}
		cookie2 := &http.Cookie{Name: "buvid3", Value: biu.BuVID}
		cookie3 := &http.Cookie{Name: "sid", Value: biu.SID}
		cookie4 := &http.Cookie{Name: "finger", Value: GetConfig().Finger}
		cookie0 := &http.Cookie{Name: "PVID", Value: "4"}

		reqPoint.AddCookie(cookie0)
		reqPoint.AddCookie(cookie1)
		reqPoint.AddCookie(cookie2)
		reqPoint.AddCookie(cookie3)
		reqPoint.AddCookie(cookie4)

		reqPoint.Header.Add("accept", "application/json, text/javascript, */*; q=0.01")
		reqPoint.Header.Add("accept-encoding", "deflate, br")
		reqPoint.Header.Add("origin", "https://passport.bilibili.com")
		reqPoint.Header.Add("referer", "https://passport.bilibili.com/login")
		reqPoint.Header.Add("sec-fetch-dest", "empty")
		reqPoint.Header.Add("sec-fetch-mode", "cors")
		reqPoint.Header.Add("sec-fetch-site", "same-origin")

		reqPoint.Header.Add("x-requested-with", "XMLHttpRequest")
		reqPoint.Header.Add("content-type", "application/x-www-form-urlencoded; charset=UTF-8")
	}, data)

	if res != nil && err == nil {
		var result BiliInfo

		json.NewDecoder(res.Body).Decode(&result)
		cookies := res.Cookies()
		if len(cookies) == 8 {
			cron.Stop()
			biu.DedeUserID = cookies[0].Value
			biu.DedeUserIDMD5 = cookies[2].Value
			biu.SESSDATA = cookies[4].Value
			biu.BiliJCT = cookies[6].Value
			biu.Login = true
			Info("Login Succeed~", logrus.Fields{"UID": biu.DedeUserID})
			// biu.GetBiliCoinLog()

			biu.InfoUpdate()
		} else {
			Info(result.Message)
		}
	}
}

// 更新信息
func (biu *BiliUser) InfoUpdate() {
	conf := GetConfig()
	for k, _ := range conf.BiU {
		if biu.DedeUserID == conf.BiU[k].DedeUserID {
			conf.BiU[k] = *biu
			if err := conf.SetConfig(); err != nil {
				Warn("error json setting")
			}
			return
		}
	}
	conf.BiU = append(conf.BiU, *biu)
	if err := conf.SetConfig(); err != nil {
		Warn("error json setting")
	}
}

// 等待扫描
func (biu *BiliUser) BiliScanAwait() {
	i := 0
	c := cron.New()
	spec := "*/4 * * * * ?"
	_ = c.AddFunc(spec, func() {
		i++
		biu.GetBiliLoginInfo(c)
	})
	c.Start()
}

// 硬币日志获取
func (biu *BiliUser) GetBiliCoinLog() {
	url := "https://api.bilibili.com/x/member/web/coin/log?jsonp=jsonp"
	res, err := GET(url, func(reqPoint *http.Request) {
		biu.NormalAuthHeader(reqPoint)

		reqPoint.Header.Add("origin", "https://account.bilibili.com")
		reqPoint.Header.Add("referer", "https://account.bilibili.com/account/coin")
	})

	if res != nil && err == nil {
		var msg BiliInfo
		json.NewDecoder(res.Body).Decode(&msg)

		bl, _ := json.Marshal(msg)
		reg := regexp.MustCompile("BV[a-zA-Z0-9_]+")
		g := reg.FindAllString(string(bl), -1)
		biu.BlockBVList = g
		biu.InfoUpdate()

		for k, _ := range msg.Data.List {
			if isToday(msg.Data.List[k].Time) {
				reg := regexp.MustCompile("BV[a-zA-Z0-9_]+")
				g := reg.FindAllString(msg.Data.List[k].Reason, -1)
				if len(g) > 0 {
					biu.DropCoinCount -= msg.Data.List[k].Delta
				}
			} else {
				break
			}
		}
	}
	Info("coin log", logrus.Fields{"dropCount": biu.DropCoinCount, "UID": biu.DedeUserID})
}

func (biu *BiliUser) NormalAuthHeader(reqPoint *http.Request) {
	cookie3 := &http.Cookie{Name: "sid", Value: biu.SID}
	cookie1 := &http.Cookie{Name: "_uuid", Value: biu.UUID}
	cookie2 := &http.Cookie{Name: "buvid3", Value: biu.BuVID}
	cookie4 := &http.Cookie{Name: "finger", Value: GetConfig().Finger}
	cookie0 := &http.Cookie{Name: "PVID", Value: "6"}
	cookie8 := &http.Cookie{Name: "SESSDATA", Value: biu.SESSDATA}
	cookie5 := &http.Cookie{Name: "DedeUserID", Value: biu.DedeUserID}
	cookie6 := &http.Cookie{Name: "DedeUserID__ckMd5", Value: biu.DedeUserIDMD5}
	cookie7 := &http.Cookie{Name: "bili_jct", Value: biu.BiliJCT}

	reqPoint.AddCookie(cookie0)
	reqPoint.AddCookie(cookie1)
	reqPoint.AddCookie(cookie2)
	reqPoint.AddCookie(cookie3)
	reqPoint.AddCookie(cookie4)
	reqPoint.AddCookie(cookie5)
	reqPoint.AddCookie(cookie6)
	reqPoint.AddCookie(cookie7)
	reqPoint.AddCookie(cookie8)

	reqPoint.Header.Add("accept", "application/json, text/plain, */*")
	reqPoint.Header.Add("accept-encoding", "deflate, br")
	reqPoint.Header.Add("sec-fetch-dest", "empty")
	reqPoint.Header.Add("sec-fetch-mode", "cors")
	reqPoint.Header.Add("sec-fetch-site", "same-origin")
}

func (biu *BiliUser) DropCoin(bv string) {
	aid := BVCovertDec(bv)
	if biu.DropCoinCount > 4 {
		Info("number of coins tossed today >= 5", logrus.Fields{"BVID": bv, "AVID": aid, "UID": biu.DedeUserID, "dropCount": biu.DropCoinCount})
		SendMessage2WeChat(biu.DedeUserID + "打赏上限")
		return
	}
	url := "https://api.bilibili.com/x/web-interface/coin/add?" + "aid=" + aid + "&multiply=1&select_like=1&cross_domain=true&csrf=" + biu.BiliJCT

	res, err := Post(url, func(reqPoint *http.Request) {
		biu.NormalAuthHeader(reqPoint)
		reqPoint.Header.Add("origin", "https://account.bilibili.com")
		reqPoint.Header.Add("referer", "https://www.bilibili.com/video/"+bv+"/?spm_id_from=333.788.videocard.0")
	})

	if res != nil && err == nil {
		var msg BiliInfo
		_ = json.NewDecoder(res.Body).Decode(&msg)
		if msg.Message == "0" {
			biu.DropCoinCount++
			Info("Drop coin succeed", logrus.Fields{"BVID": bv, "AVID": aid, "UID": biu.DedeUserID, "dropCount": biu.DropCoinCount})
			SendMessage2WeChat(biu.DedeUserID + "打赏" + bv + "成功")
		} else if msg.Message == "超过投币上限啦~" {
			Info("Drop coin limited", logrus.Fields{"BVID": bv, "AVID": aid, "UID": biu.DedeUserID, "dropCount": biu.DropCoinCount})
		}
	}
}

func (biu *BiliUser) RandDrop() {
	bvs := GetGuichuBVs()
	rand.Seed(time.Now().UnixNano())
	randIndex := rand.Intn(len(bvs))
	biu.DropCoin(bvs[randIndex])
}

func LoadUser() []BiliUser {
	return GetConfig().BiU
}

func GetUser(uid string) (*BiliUser, error) {
	users := LoadUser()
	for k, _ := range users {
		if users[k].DedeUserID == uid {
			users[k].GetBiliCoinLog()
			return &users[k], nil
		}
	}
	return nil, errors.New("not found user")
}

func DelUser(uid string) error {
	users := LoadUser()
	for k, _ := range users {
		if users[k].DedeUserID == uid {
			users = append(users[:k], users[k+1:]...)
			config.BiU = users
			_ = config.SetConfig()
		}
	}
	return errors.New("not found user")
}

func GetAllUID() []string {
	users := LoadUser()
	var retVal []string
	for k, _ := range users {
		retVal = append(retVal, users[k].DedeUserID)
	}
	return retVal
}

func CronDrop(biu BiliUser) {
	c := cron.New()
	_ = c.AddFunc(GetConfig().Cron, func() {
		biu.GetBiliCoinLog()
		for i := 0; i < 5; i++ {
			biu.RandDrop()
			time.Sleep(Random(60))
			Info("cron finish", logrus.Fields{"UID": biu.DedeUserID})
		}
	})
	c.Start()
}

func CronDropReg() {
	bius := GetConfig().BiU
	if len(bius) == 0 {
		println("not found users")
		return
	}
	for k, _ := range bius {
		println("cron - add user " + bius[k].DedeUserID)
		CronDrop(bius[k])
	}
}
