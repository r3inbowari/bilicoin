package bilicoin

import (
	"encoding/json"
	"errors"
	qrcodeTerminal "github.com/Baozisoftware/qrcode-terminal-go"
	"github.com/robertkrimen/otto"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
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

var inlineJSCode = `func test() {}`

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

func _buvidGenerate() (string, error) {
	url := "https://data.bilibili.com/v/web/web_page_view?mid=null&fts=null&url=https%253A%252F%252Fwww.bilibili.com%252F&proid=3&ptype=2&module=game&title=%E5%93%94%E5%93%A9%E5%93%94%E5%93%A9%20(%E3%82%9C-%E3%82%9C)%E3%81%A4%E3%83%AD%20%E5%B9%B2%E6%9D%AF~-bilibili&ajaxtag=&ajaxid=&page_ref=https%253A%252F%252Fpassport.bilibili.com%252Flogin"
	res, err := GET(url, nil)
	if res != nil {
		cook := res.Cookies()
		return cook[1].Value, err
	}
	return "", errors.New("nil buvid response")
}

func _uuidGenerate() (string, error) {
	bytes, err := ioutil.ReadFile("bili.js")
	vm := otto.New()
	_, err = vm.Run(bytes)
	uuid, err := vm.Call("generateUuid", nil)
	return uuid.String(), err
}

func (biu *BiliUser) QRCodePrint() {
	_QRCPrint(biu.OAuth.Url)
}

func _QRCPrint(content string) {
	obj := qrcodeTerminal.New2(qrcodeTerminal.ConsoleColors.BrightBlue, qrcodeTerminal.ConsoleColors.BrightGreen, qrcodeTerminal.QRCodeRecoveryLevels.Low)
	obj.Get(content).Print()
}

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

	// QRCodePrint(qr.Data.Url)
	biu.OAuth = qr.Data
	biu.SID = res.Cookies()[0].Value
	return nil
}

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
		// json.NewDecoder(res.Body).Decode(&biu.Bi)
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

func isToday(timeStr string) bool {
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", timeStr, time.Local)
	if t.Unix() > GetZeroTime() && t.Unix() < GetLastTime() {
		return true
	}
	return false
}

func GetZeroTime() int64 {
	currentTime := time.Now()
	return time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location()).Unix()
}

func GetLastTime() int64 {
	currentTime := time.Now()
	return time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 23, 59, 59, 0, currentTime.Location()).Unix()
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

func GetGuichuBVs() []string {
	res, _ := GET("https://api.bilibili.com/x/web-interface/ranking/region?rid=119&day=3&original=0", nil)
	result, _ := ioutil.ReadAll(res.Body)
	reg := regexp.MustCompile("BV[a-zA-Z0-9_]{10}")
	return reg.FindAllString(string(result), -1)
}

//func GetGuichuBVs() []string {
//	res, _ := GET("https://api.bilibili.com/x/web-interface/ranking/region?rid=119&day=3&original=0", nil)
//	var info BiliInfo
//	_ = json.NewDecoder(res.Body).Decode(&res)
//	var ret []string
//	return ret
//}

func (biu *BiliUser) DropCoin(bv string) {
	aid := BVCovertDec(bv)
	if biu.DropCoinCount > 4 {
		Info("number of coins tossed today >= 5", logrus.Fields{"BVID": bv, "AVID": aid, "UID": biu.DedeUserID, "dropCount": biu.DropCoinCount})
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
		} else if msg.Message == "超过投币上限啦~" {
			Info("Drop coin limited", logrus.Fields{"BVID": bv, "AVID": aid, "UID": biu.DedeUserID, "dropCount": biu.DropCoinCount})
		}
	}
}

var xor = 177451812
var add = 8728348608
var table = "fZodR9XQDSUm21yCkr6zBqiveYah8bt4xsWpHnJE7jL5VG3guMTKNPAwcF"
var s = []int{11, 10, 3, 8, 4, 6, 2, 9, 5, 7}

func BVCovertDec(bv string) string {
	var tr = make(map[byte]int)
	for i := 1; i < 58; i++ {
		tr[table[i]] = i
	}
	r := 0
	for i := 0; i < 6; i++ {
		r += tr[bv[s[i]]] * int(math.Pow(58.0, float64(i)))
	}
	retAV := (r - add) ^ xor
	return strconv.Itoa(retAV)
}

func BVCovertEnc(av string) string {
	var r = []string{"B", "V", "1", "", "", "4", "", "1", "", "7", "", ""}
	x, _ := strconv.Atoi(av)
	x_ := (x ^ xor) + add
	for i := 0; i < 6; i++ {
		_x := int(math.Pow(58.0, float64(i)))
		aj := int(math.Floor(float64(x_ / _x)))
		r[s[i]] = string(table[aj%58])
	}
	return r[0] + r[1] + r[2] + r[3] + r[4] + r[5] + r[6] + r[7] + r[8] + r[9] + r[10] + r[11]
}

func (biu *BiliUser) RandDrop() {
	bvs := GetGuichuBVs()
	rand.Seed(time.Now().UnixNano())
	randIndex := rand.Intn(len(bvs))
	biu.DropCoin(bvs[randIndex])
}
