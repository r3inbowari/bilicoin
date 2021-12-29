package bilicoin

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/r3inbowari/common"
	. "github.com/r3inbowari/zlog"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
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
	ConvertCoin   bool       `json:"sw_convert_coin"`           // Config -> 银瓜子兑换硬币
	Job           *cron.Cron
}

type ABi struct {
	Gold      int       `json:"gold"`       // 金瓜子
	Silver    int       `json:"silver"`     // 银瓜子
	Coins     float64   `json:"coins"`      // 硬币
	Uname     string    `json:"uname"`      // 昵称
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

// BiliInfo bilibili 标准响应
type BiliInfo struct {
	Status    string   `json:"status"`  // 响应状态
	Data      BiliData `json:"data"`    // 数据
	Message   string   `json:"message"` // 消息
	Timestamp int64    `json:"ts"`      // ts
	TTL       int      `json:"ttl"`     // 生存
	Code      int      `json:"code"`    // 响应代码
}

type BiliData struct {
	List     []CoinLog `json:"list,omitempty"`
	Like     string    `json:"like,omitempty"`
	BVID     string    `json:"bvid,omitempty"`
	BillCoin float64   `json:"billCoin,omitempty"`
	Gold     int       `json:"gold,omitempty"`
	Silver   int       `json:"silver,omitempty"`
	Face     string    `json:"face"`
	Uname    string    `json:"uname"`
}

type CoinLog struct {
	Time   string `json:"time"`   // 时间
	Delta  int    `json:"delta"`  // 增量
	Reason string `json:"reason"` // 原因
}

var releaseVersion = "v1.0.0" // release tag
var releaseTag = "b4ac8f4..6599638 @master"

type Version struct {
	Patch int64
	Minor int64
	Major int64
}

var version Version

var RunningMode = ""

const (
	Simple = "simple"
	Api    = "api"
)

// CreateUser 创建用户
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

// QRCodePrint 打印二维码
func (biu *BiliUser) QRCodePrint() {
	QRCPrint(biu.OAuth.Url)
}

// GetQRCode 获取二维码
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

// GetBiliLoginInfo 获取登录信息 Cron
func (biu *BiliUser) GetBiliLoginInfo(cron *cron.Cron) {
	url := "https://passport.bilibili.com/qrcode/getLoginInfo"

	data := "oauthKey=" + biu.OAuth.OAuthKey + "&gourl=" + "https%3A%2F%2Fwww.bilibili.com%2F"
	res, err := Post2(url, func(reqPoint *http.Request) {

		cookie1 := &http.Cookie{Name: "_uuid", Value: biu.UUID}
		cookie2 := &http.Cookie{Name: "buvid3", Value: biu.BuVID}
		cookie3 := &http.Cookie{Name: "sid", Value: biu.SID}
		cookie4 := &http.Cookie{Name: "finger", Value: GetConfig(false).Finger}
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
			Log.WithFields(logrus.Fields{"UID": biu.DedeUserID}).Info("Login Succeed~")

			biu.InfoUpdate()
		} else {
			if result.Message != "" {
				Log.Info(result.Message)
			}
		}
	}
}

// LoginCallback 获取登录信息
func (biu *BiliUser) LoginCallback(callback func(isLogin bool)) {
	url := "https://passport.bilibili.com/qrcode/getLoginInfo"

	data := "oauthKey=" + biu.OAuth.OAuthKey + "&gourl=" + "https%3A%2F%2Fwww.bilibili.com%2F"
	res, err := Post2(url, func(reqPoint *http.Request) {

		cookie1 := &http.Cookie{Name: "_uuid", Value: biu.UUID}
		cookie2 := &http.Cookie{Name: "buvid3", Value: biu.BuVID}
		cookie3 := &http.Cookie{Name: "sid", Value: biu.SID}
		cookie4 := &http.Cookie{Name: "finger", Value: GetConfig(false).Finger}
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
		reqPoint.Header.Add("sec-fetch-RunningMode", "cors")
		reqPoint.Header.Add("sec-fetch-site", "same-origin")

		reqPoint.Header.Add("x-requested-with", "XMLHttpRequest")
		reqPoint.Header.Add("content-type", "application/x-www-form-urlencoded; charset=UTF-8")
	}, data)

	if res != nil && err == nil {
		var result BiliInfo

		_ = json.NewDecoder(res.Body).Decode(&result)
		cookies := res.Cookies()
		if len(cookies) == 8 {
			biu.DedeUserID = cookies[0].Value
			biu.DedeUserIDMD5 = cookies[2].Value
			biu.SESSDATA = cookies[4].Value
			biu.BiliJCT = cookies[6].Value
			biu.Login = true
			biu.InfoUpdate()
			callback(true)
		} else {
			callback(false)
		}
	}
}

// InfoUpdate 更新配置树
func (biu *BiliUser) InfoUpdate() {
	conf := GetConfig(false)
	for k, _ := range conf.BiU {
		if biu.DedeUserID == conf.BiU[k].DedeUserID {
			conf.BiU[k] = *biu
			if err := conf.SetConfig(); err != nil {
				Log.Warn("[FILE] error json setting")
			}
			return
		}
	}
	conf.BiU = append(conf.BiU, *biu)
	if err := conf.SetConfig(); err != nil {
		Log.Warn("[FILE] error json setting")
	}
}

// BiliScanAwait 等待登录扫描
func (biu *BiliUser) BiliScanAwait() {
	i := 0
	c := cron.New()
	spec := "*/4 * * * * ?"
	_ = c.AddFunc(spec, func() {
		i++
		biu.GetBiliLoginInfo(c)
	})
	c.Start()

	for true {
		if biu.DedeUserID != "" {
			Log.Info("Login process will exit after 5 seconds")
			time.Sleep(time.Second * 5)
			os.Exit(0)
		}
	}
}

func (biu *BiliUser) Silver2Coin() error {
	url := "https://api.live.bilibili.com/xlive/revenue/v1/wallet/silver2coin"

	//payload := strings.NewReader()

	csrf := `csrf_token=7d2e2fbbac654107f0a577d12cd7a48e&csrf=7d2e2fbbac654107f0a577d12cd7a48e`

	res, err := Post2(url, func(reqPoint *http.Request) {
		biu.NormalAuthHeader(reqPoint)

		// raw 提交
		reqPoint.Header.Add("content-type", "application/x-www-form-urlencoded")
		reqPoint.Header.Add("origin", "https://link.bilibili.com")
		reqPoint.Header.Add("referer", "https://link.bilibili.com/p/center/index")
	}, csrf)

	if res != nil && err == nil {
		var msg BiliInfo
		_ = json.NewDecoder(res.Body).Decode(&msg)

		if msg.Message == "兑换成功" {
			Log.Info("[TASK] use 700 silver to one bili coin", logrus.Fields{"remain silver": msg.Data.Silver})
		} else {
			Log.Warn("[TASK] bili coin convert failed")
		}
	}
	return nil
}

func (biu *BiliUser) GetBiliWallet() error {
	url := "https://api.live.bilibili.com/xlive/web-ucenter/user/get_user_info"
	res, err := GET(url, func(reqPoint *http.Request) {
		biu.NormalAuthHeader(reqPoint)

		reqPoint.Header.Add("origin", "https://live.bilibili.com")
		reqPoint.Header.Add("referer", "https://live.bilibili.com")
	})

	if res != nil && err == nil {
		var msg BiliInfo
		_ = json.NewDecoder(res.Body).Decode(&msg)

		if msg.TTL < 1 {
			return errors.New("the service is unreachable or an unknown error has occurred")
		}

		biu.Bi.Gold = msg.Data.Gold
		biu.Bi.Coins = msg.Data.BillCoin
		biu.Bi.Silver = msg.Data.Silver
		biu.Bi.Face = msg.Data.Face
		biu.Bi.Uname = msg.Data.Uname
		biu.InfoUpdate()
	}
	return nil
}

// GetBiliCoinLog 硬币日志获取
func (biu *BiliUser) GetBiliCoinLog() error {
	url := "https://api.bilibili.com/x/member/web/coin/log?jsonp=jsonp"
	res, err := GET(url, func(reqPoint *http.Request) {
		biu.NormalAuthHeader(reqPoint)

		// reqPoint.Header.Add("accept-encoding", "gzip, deflate, br")

		reqPoint.Header.Add("origin", "https://account.bilibili.com")
		reqPoint.Header.Add("referer", "https://account.bilibili.com/account/coin")
	})

	if res != nil && err == nil {
		var msg BiliInfo

		_ = json.NewDecoder(res.Body).Decode(&msg)

		// fix: 循环投币的惨剧可能发生
		// reason: 在部分地区的分发服务器可能默认采用了gzip压缩导致乱码
		// 判断投币历史记录是否成功获取,即ttl跳数>0必然成功
		if msg.TTL < 1 {
			return errors.New("the service is unreachable or an unknown error has occurred")
		}

		// block already dropped
		bl, _ := json.Marshal(&msg)
		reg := regexp.MustCompile("BV[a-zA-Z0-9_]+")
		g := reg.FindAllString(string(bl), -1)
		biu.BlockBVList = g
		biu.InfoUpdate()
		biu.DropCoinCount = 0 // init and recount

		for k, _ := range msg.Data.List {
			if isToday(msg.Data.List[k].Time) {
				if strings.HasSuffix(msg.Data.List[k].Reason, "打赏") {
					biu.DropCoinCount -= msg.Data.List[k].Delta
				}
			} else {
				break
			}
		}
	}
	Log.WithFields(logrus.Fields{"dropCount": biu.DropCoinCount, "UID": biu.DedeUserID}).Info("[USER] Drop history")
	return nil
}

// NormalAuthHeader bilibili标准头部
func (biu *BiliUser) NormalAuthHeader(reqPoint *http.Request) {
	cookie3 := &http.Cookie{Name: "sid", Value: biu.SID}
	cookie1 := &http.Cookie{Name: "_uuid", Value: biu.UUID}
	cookie2 := &http.Cookie{Name: "buvid3", Value: biu.BuVID}
	cookie4 := &http.Cookie{Name: "finger", Value: GetConfig(false).Finger}
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
	// reqPoint.Header.Add("accept-encoding", "deflate, br")
	reqPoint.Header.Add("sec-fetch-dest", "empty")
	reqPoint.Header.Add("sec-fetch-mode", "cors")
	reqPoint.Header.Add("sec-fetch-site", "same-origin")
}

// DropCoin 打赏逻辑
func (biu *BiliUser) DropCoin(bv string) {
	// TODO fix: panic if error-bv inputted
	aid := BVCovertDec(bv)
	if biu.DropCoinCount > 4 {
		Log.WithFields(logrus.Fields{"BVID": bv, "AVID": aid, "UID": biu.DedeUserID, "dropCount": biu.DropCoinCount}).Warn("number of coins tossed today >= 5")
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
			Log.WithFields(logrus.Fields{"BVID": bv, "AVID": aid, "UID": biu.DedeUserID, "dropCount": biu.DropCoinCount}).Info("[TASK] Drop succeed")
			biu.SendMessage2WeChat(biu.DedeUserID + "打赏" + bv + "成功")
		} else if msg.Message == "超过投币上限啦~" {
			Log.WithFields(logrus.Fields{"BVID": bv, "AVID": aid, "UID": biu.DedeUserID, "dropCount": biu.DropCoinCount}).Info("[TASK] Drop limited")
		}
	}
}

// SendMessage2WeChat 方糖
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

var taskMap sync.Map

// DropTaskStart 投币任务
func (biu *BiliUser) DropTaskStart() {
	c := cron.New()
	uuid := CreateUUID()
	taskMap.Store(uuid, c)
	Log.WithFields(logrus.Fields{"UID": biu.DedeUserID, "Cron": biu.Cron, "TaskID": uuid}).Info("[CRON] Add task")
	_ = c.AddFunc(biu.Cron, func() {
		// TODO fix: loop drop may be happen...
		_ = biu.GetBiliCoinLog()
		_ = biu.GetBiliWallet()

		// 使用银瓜子兑换硬币一枚
		if biu.ConvertCoin && biu.Bi.Silver >= 700 {
			biu.Silver2Coin()
		}

		for true {
			if biu.DropCoinCount > 4 {
				biu.InfoUpdate()
				break
			}
			biu.RandDrop()
			time.Sleep(Random(60))
		}
		Log.WithFields(logrus.Fields{"UID": biu.DedeUserID}).Info("[CRON] Task Completed")
	})
	c.Start()
}

// CronTaskLoad 注册所有投币任务
func CronTaskLoad() {
	// panic: if not biu in config file
	// exit code 1001
	bius := GetConfig(true).BiU
	if len(bius) == 0 && RunningMode == Simple {
		Log.Warn("[CRON] biu not found: please make sure that at least one user cookies exists in bili.json file")
		Log.Warn("[CRON] tip: use '-n' option to create a new user cookies by bilibili-mobile-client QR Login")
		Log.Info("[CRON] EXIT 1001")
		time.Sleep(time.Second * 5)
		os.Exit(1001)
	}
	if len(bius) == 0 {
		Log.Info("[USER] Not found users")
		return
	}
	for k, _ := range bius {
		// 投币日志
		err := bius[k].GetBiliCoinLog()
		if err != nil {
			// 当前账号获取日志失败，跳过此账号
			// 过期、未知错误、服务不可达、解析错误
			Log.WithFields(logrus.Fields{"err": err.Error(), "id": bius[k].DedeUserID}).Warn("user load error, can not get coin log")
			continue
		}
		// 账号信息获取
		err = bius[k].GetBiliWallet()
		if err != nil {
			// 当前账号获取用户信息失败，跳过此账号
			// 过期、未知错误、服务不可达、解析错误
			Log.WithFields(logrus.Fields{"err": err.Error(), "id": bius[k].DedeUserID}).Warn("user load error, can not get coin log")
			continue
		}
		// 投币任务启动
		bius[k].DropTaskStart()
	}
}

// LoadUser 获取所有UID实体
func LoadUser() []BiliUser {
	return GetConfig(false).BiU
}

// GetUser 尝试通过UID获取一个UID实体
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

// DelUser 尝试删除Config中一个UID配置实体
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

// GetAllUID 获取Config中所有已知UID的String
func GetAllUID() []string {
	users := LoadUser()
	var retVal []string
	for k, _ := range users {
		retVal = append(retVal, users[k].DedeUserID)
	}
	return retVal
}

func UserList() {
	users := LoadUser()
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
}

var BuildMode common.Mode

// InitBili bilicoin初始化
func InitBili(buildMode string, ver, hash string, major, minor, patch string) {
	BuildMode = common.Modes[buildMode]
	releaseVersion = ver
	releaseTag = hash
	version.Major, _ = strconv.ParseInt(major, 10, 64)
	version.Minor, _ = strconv.ParseInt(minor, 10, 64)
	version.Patch, _ = strconv.ParseInt(patch, 10, 64)
	InitConfig()
	if lp := GetConfig(false).LoggerLevel; lp != nil {
		parseLevel, _ := logrus.ParseLevel(*lp)
		Log.SetLevel(parseLevel)
	}
}
