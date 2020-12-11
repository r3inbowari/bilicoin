package bilicoin

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
	"math/rand"
	"net/http"
	"regexp"
	"time"
)

type BiliQRCode struct {
	Data      QRCodeData `json:"data"`   // QR数据
	Timestamp int64      `json:"ts"`     // 时间戳
	Status    bool       `json:"status"` // 扫描状态
	Code      int        `json:"code"`   // 状态码
}

type QRCodeData struct {
	Url      string `json:"url"`      // 地址
	OAuthKey string `json:"oauthKey"` // oa-key
}

type BiliUser struct {
	UUID          string     `json:"_uuid"`                     // CooKie -> UUID
	BuVID         string     `json:"buvid"`                     // CooKie -> Buv跟踪
	OAuth         QRCodeData `json:"oauth"`                     // CooKie -> QR登录信息
	SID           string     `json:"sid"`                       // CooKie -> SID
	BiliJCT       string     `json:"bili_jct"`                  // CooKie -> JCT
	SESSDATA      string     `json:"SESSDATA"`                  // CooKie -> Session
	DedeUserID    string     `json:"uid"`                       // CooKie -> UID
	DedeUserIDMD5 string     `json:"DedeUserID__ckMd5"`         // CooKie -> UID_Decode
	Bi            ABi        `json:"bi"`                        // Config -> 用户信息
	Login         bool       `json:"login"`                     // Config -> 弃用
	Expire        time.Time  `json:"expire"`                    // CooKie -> 失效时间
	DropCoinCount int        `json:"drop_coin_count,omitempty"` // Config -> 当天投币数量
	BlockBVList   []string   `json:"block_bv_list"`             // Config -> 禁止列表
	Cron          string     `json:"cron"`                      // Config -> 执行表达式
	FT            string     `json:"ft"`                        // Config -> 方糖[可选]
	FTSwitch      bool       `json:"ft_switch"`                 // Config -> 方糖开关
}

type ABi struct {
	BCoins    int       `json:"bCoins"`     // B柯拉
	Coins     int       `json:"coins"`      // 硬币
	Nick      string    `json:"uname"`      // 昵称
	Face      string    `json:"face"`       // 头像地址
	Status    string    `json:"userStatus"` // 账号状态
	LevelInfo LevelInfo `json:"level_info"` // 等级
}

type LevelInfo struct {
	CurrentLevel int `json:"current_level"` // 当前等级
	CurrentMin   int `json:"current_min"`   // 最小经验
	CurrentExp   int `json:"current_exp"`   // 当前经验
	NextExp      int `json:"next_exp"`      // 下一等级
}

// bilibili 标准响应
type BiliInfo struct {
	Status    string   `json:"status"`  // 响应状态
	Data      BiliData `json:"data"`    // 数据
	Message   string   `json:"message"` // 消息
	Timestamp int64    `json:"ts"`      // ts
	TTL       int      `json:"ttl"`     // 生存
	Code      int      `json:"code"`    // 响应代码
}

type BiliData struct {
	List []CoinLog `json:"list,omitempty"`
	Like string    `json:"like,omitempty"`
	BVID string    `json:"bvid,omitempty"`
}

type CoinLog struct {
	Time   string `json:"time"`   // 时间
	Delta  int    `json:"delta"`  // 增量
	Reason string `json:"reason"` // 原因
}

var releaseVersion = "v1.0.2 build on 08 12 2020 fb38f8f..53a9431 master" // release tag

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
	biu.Cron = "30 50 23 * * ?"
	biu.FT = ""
	biu.FTSwitch = false
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

		_ = json.NewDecoder(res.Body).Decode(&result)
		cookies := res.Cookies()
		if len(cookies) == 8 {
			cron.Stop()
			biu.DedeUserID = cookies[0].Value
			biu.DedeUserIDMD5 = cookies[2].Value
			biu.SESSDATA = cookies[4].Value
			biu.BiliJCT = cookies[6].Value
			biu.Login = true
			Info("Login Succeed~", logrus.Fields{"UID": biu.DedeUserID})

			biu.InfoUpdate()
		} else {
			Info(result.Message)
		}
	}
}

// 更新配置树
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

// 等待登录扫描
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
		_ = json.NewDecoder(res.Body).Decode(&msg)

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

// bilibili标准头部
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

// 打赏逻辑
func (biu *BiliUser) DropCoin(bv string) {
	aid := BVCovertDec(bv)
	if biu.DropCoinCount > 4 {
		Info("number of coins tossed today >= 5", logrus.Fields{"BVID": bv, "AVID": aid, "UID": biu.DedeUserID, "dropCount": biu.DropCoinCount})
		biu.SendMessage2WeChat(biu.DedeUserID + "打赏上限")
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
			biu.SendMessage2WeChat(biu.DedeUserID + "打赏" + bv + "成功")
		} else if msg.Message == "超过投币上限啦~" {
			Info("Drop coin limited", logrus.Fields{"BVID": bv, "AVID": aid, "UID": biu.DedeUserID, "dropCount": biu.DropCoinCount})
		}
	}
}

// 方糖
func (biu *BiliUser) SendMessage2WeChat(title string, content ...string) {
	ft := biu.FT
	if biu.FTSwitch && ft != "" {
		if len(content) > 0 {
			GetRequest("https://sc.ftqq.com/" + ft + ".send?desp=" + content[0] + "&text=" + title)
		} else {
			GetRequest("https://sc.ftqq.com/" + ft + ".send?text=" + title)
		}
	}
}

func (biu *BiliUser) RandDrop() {
	bvs := GetGuichuBVs()
	rand.Seed(time.Now().UnixNano())
	randIndex := rand.Intn(len(bvs))
	biu.DropCoin(bvs[randIndex])
}

// 投币任务
func (biu *BiliUser) DropTaskStart() {
	c := cron.New()
	Info("cron add task", logrus.Fields{"UID": biu.DedeUserID, "Cron": biu.Cron})
	_ = c.AddFunc(biu.Cron, func() {
		biu.GetBiliCoinLog()
		Info("get coin log", logrus.Fields{"UID": biu.DedeUserID, "Cron": biu.DropCoinCount})
		if biu.DropCoinCount > 4 {
			Info("cron task not need", logrus.Fields{"UID": biu.DedeUserID})
			return
		}
		for true {
			if biu.DropCoinCount > 4 {
				break
			}
			biu.RandDrop()
			time.Sleep(Random(60))
		}
		Info("cron finish", logrus.Fields{"UID": biu.DedeUserID})
	})
	c.Start()
}

// 注册所有投币任务
func CronTaskLoad() {
	bius := GetConfig().BiU
	if len(bius) == 0 {
		Info("not found users")
		return
	}
	for k, _ := range bius {
		bius[k].DropTaskStart()
	}
}

// 获取所有UID实体
func LoadUser() []BiliUser {
	return GetConfig().BiU
}

// 尝试通过UID获取一个UID实体
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

// 尝试删除Config中一个UID配置实体
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

// 获取Config中所有已知UID的String
func GetAllUID() []string {
	users := LoadUser()
	var retVal []string
	for k, _ := range users {
		retVal = append(retVal, users[k].DedeUserID)
	}
	return retVal
}

// bilicoin初始化
func InitBili() {
	fmt.Println("bilicoin " + releaseVersion)
	InitConfig()
	InitLogger()
}
