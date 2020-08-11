package bilicoin

import (
	"encoding/json"
	"errors"
	qrcodeTerminal "github.com/Baozisoftware/qrcode-terminal-go"
	"github.com/robertkrimen/otto"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
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
	List []interface{} `json:"list"`
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
	res, err := GET("https://passport.bilibili.com/qrcode/getLoginUrl", func(reqPoint *http.Request) {

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
		reqPoint.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.105 Safari/537.36")
	}, data)

	if res != nil && err == nil {
		var result BiliInfo

		json.NewDecoder(res.Body).Decode(&result)
		// json.NewDecoder(res.Body).Decode(&biu.Bi)
		cookies := res.Cookies()
		if len(cookies) == 8 {
			biu.DedeUserID = cookies[0].Value
			biu.DedeUserIDMD5 = cookies[2].Value
			biu.SESSDATA = cookies[4].Value
			biu.BiliJCT = cookies[6].Value
			cron.Stop()
			biu.Login = true
			Info("Login Succeed~", logrus.Fields{"UID": biu.DedeUserID})
			biu.GetBiliCoinLog()
		} else {
			Info(result.Message)
		}
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
		reqPoint.Header.Add("origin", "https://account.bilibili.com")
		reqPoint.Header.Add("referer", "https://account.bilibili.com/account/coin")
		reqPoint.Header.Add("sec-fetch-dest", "empty")
		reqPoint.Header.Add("sec-fetch-mode", "cors")
		reqPoint.Header.Add("sec-fetch-site", "same-origin")

		reqPoint.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.105 Safari/537.36")

	})

	if res != nil && err == nil {
		var msg BiliInfo
		json.NewDecoder(res.Body).Decode(&msg)
		println(len(msg.Data.List))

		str := msg.Data.List[0].(CoinLog)
		time1, _ := time.ParseInLocation("2006-01-02 15:04:05", str.Time, time.Local)
		println(time1.String())

		println(time1.String())
		println(time1.String())
	}

}
