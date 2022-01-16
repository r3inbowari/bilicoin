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
	"sort"
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
	UUID          string     `json:"_uuid"`             // CooKie -> UUID
	BuVID         string     `json:"buvid"`             // CooKie -> Buv跟踪
	OAuth         QRCodeData `json:"oauth"`             // CooKie -> QR登录信息
	SID           string     `json:"sid"`               // CooKie -> SID
	BiliJCT       string     `json:"bili_jct"`          // CooKie -> JCT
	SESSDATA      string     `json:"SESSDATA"`          // CooKie -> Session
	DedeUserID    string     `json:"uid"`               // CooKie -> UID
	DedeUserIDMD5 string     `json:"DedeUserID__ckMd5"` // CooKie -> UID_Decode
	Bi            ABi        `json:"bi"`                // Config -> 用户信息
	Login         bool       `json:"login"`             // Config -> 弃用
	Expire        time.Time  `json:"expire"`            // CooKie -> 失效时间
	DropCoinCount int        `json:"-"`                 // Mem -> 当天投币数量
	BlockBVList   []string   `json:"block_bv_list"`     // Config -> 禁止列表
	Cron          string     `json:"cron"`              // Config -> 执行表达式
	FT            string     `json:"ft"`                // Config -> 方糖[可选]
	FTSwitch      bool       `json:"ft_switch"`         // Config -> 方糖开关
	Tasks         []string   `json:"tasks"`             // Config -> 任务列表
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
	TTL       int      `json:"ttl"`     // 跳数
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
			if cron != nil {
				cron.Stop()
			}
			biu.DedeUserID = cookies[0].Value
			biu.DedeUserIDMD5 = cookies[2].Value
			biu.SESSDATA = cookies[4].Value
			biu.BiliJCT = cookies[6].Value
			biu.Login = true
			Log.WithFields(logrus.Fields{"UID": biu.DedeUserID}).Info("login Succeed~")
			biu.Tasks = []string{"drop-coin"}
			biu.InfoUpdate()
		} else {
			if result.Message != "" {
				Log.Info(result.Message)
			}
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
			Log.Info("login process will exit after 5 seconds")
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

// DropCoin 打赏逻辑
func (biu *BiliUser) DropCoin(bv Video) {
	// aid := BVCovertDec(bv)
	if biu.DropCoinCount > 4 {
		Log.WithFields(logrus.Fields{"BVID": bv.Bvid, "AVID": bv.Aid, "UID": biu.DedeUserID, "dropCount": biu.DropCoinCount}).Warn("number of coins tossed today >= 5")
		biu.SendMessage2WeChat(biu.DedeUserID + "打赏上限")
		return
	}

	url := fmt.Sprintf("https://api.bilibili.com/x/web-interface/coin/add?aid=%d&multiply=1&select_like=1&cross_domain=true&csrf=%s", bv.Aid, biu.BiliJCT)
	res, err := Post(url, func(reqPoint *http.Request) {
		biu.NormalAuthHeader(reqPoint)
		reqPoint.Header.Add("origin", "https://account.bilibili.com")
		reqPoint.Header.Add("referer", fmt.Sprintf("https://www.bilibili.com/video/%s/?spm_id_from=333.788.videocard.0", bv.Bvid))
	})

	if res != nil && err == nil {
		var msg BiliInfo
		_ = json.NewDecoder(res.Body).Decode(&msg)
		if msg.TTL > 0 {
			// 预扣1币
			biu.DropCoinCount++
		}
		if msg.Message == "0" {
			Log.WithFields(logrus.Fields{"TTL": msg.TTL, "M": msg.Message, "BVID": bv.Bvid, "AVID": bv.Aid, "UID": biu.DedeUserID, "dropCount": biu.DropCoinCount}).Info("[TASK] drop succeed")
			biu.SendMessage2WeChat(biu.DedeUserID + "打赏" + bv.Bvid + "成功")
		} else if msg.Message == "超过投币上限啦~" {
			// 投币失败，退回1
			biu.DropCoinCount--
			Log.WithFields(logrus.Fields{"TTL": msg.TTL, "M": msg.Message, "BVID": bv.Bvid, "AVID": bv.Aid, "UID": biu.DedeUserID, "dropCount": biu.DropCoinCount}).Info("[TASK] drop limited")
		} else if msg.Message == "账号未登录" {
			Log.WithFields(logrus.Fields{"TTL": msg.TTL, "M": msg.Message, "BVID": bv.Bvid, "AVID": bv.Aid, "UID": biu.DedeUserID, "dropCount": biu.DropCoinCount}).Warn("[TASK] expired token")
		} else {
			Log.WithFields(logrus.Fields{"TTL": msg.TTL, "M": msg.Message, "BVID": bv.Bvid, "AVID": bv.Aid, "UID": biu.DedeUserID, "dropCount": biu.DropCoinCount}).Info("[TASK] drop unknown status")
		}
	} else {
		Log.WithFields(logrus.Fields{"cnt": biu.DropCoinCount, "id": biu.DedeUserID, "url": url}).Error("[DBG] coin error.res nil")
		if err != nil {
			Log.WithFields(logrus.Fields{"cnt": biu.DropCoinCount, "id": biu.DedeUserID, "url": url}).Error("[DBG] coin error")
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
	bvs := GetPopulars()
	rand.Seed(time.Now().UnixNano())
	randIndex := rand.Intn(len(bvs))
	biu.DropCoin(bvs[randIndex])
}

// LoadUsers 获取所有UID实体
func LoadUsers() []BiliUser {
	return GetConfig(false).BiU
}

// GetUser 尝试通过UID获取一个UID实体
func GetUser(uid string) (*BiliUser, error) {
	users := LoadUsers()

	if userIndex := sort.Search(len(users), func(index int) bool { return users[index].DedeUserID == uid }); userIndex == len(users) {
		return nil, errors.New("not found user")
	} else {
		return &users[userIndex], users[userIndex].GetBiliCoinLog()
	}
}

// DelUser 尝试删除Config中一个UID配置实体
func DelUser(uid string) error {
	users := LoadUsers()

	if userIndex := sort.Search(len(users), func(index int) bool { return users[index].DedeUserID == uid }); userIndex == len(users) {
		return errors.New("not found user")
	} else {
		users = append(users[:userIndex], users[userIndex+1:]...)
		config.BiU = users
		return config.SetConfig()
	}
}

func searchBiliUser(users []BiliUser, uid string) *BiliUser {
	if userIndex := sort.Search(len(users), func(index int) bool { return users[index].DedeUserID == uid }); userIndex == len(users) {
		return nil
	} else {
		return &users[userIndex]
	}
}

// GetAllUID 获取Config中所有已知UID的String
func GetAllUID() []string {
	users := LoadUsers()
	var retVal []string
	for k, _ := range users {
		retVal = append(retVal, users[k].DedeUserID)
	}
	return retVal
}

func UserList() {
	users := LoadUsers()
	fmt.Println("")
	fmt.Println("total: " + strconv.Itoa(len(users)))
	fmt.Println()
	fmt.Println("|      UID \t|         Cron     \t|      FTQQ\t|")
	for _, v := range users {
		fmt.Printf("| %s \t| %s \t| %s \t|\n", v.DedeUserID, v.Cron, strconv.FormatBool(v.FTSwitch))
	}
	fmt.Println("")
}

var BuildMode common.Mode

// InitBili bilicoin初始化
func InitBili(buildMode string, ver, hash string) {
	BuildMode = common.Modes[buildMode]
	releaseVersion = ver
	releaseTag = hash
	InitConfig()
	if lp := GetConfig(false).LoggerLevel; lp != nil {
		parseLevel, _ := logrus.ParseLevel(*lp)
		Log.SetLevel(parseLevel)
	}
}

type Task func(user *BiliUser) error

var cronTask sync.Map

func BiliExecutor(user *BiliUser) {
	c := cron.New()
	uuid := CreateUUID()
	cronTask.Store(uuid, c)
	c.Start()
	for _, t := range user.Tasks {
		if task, ok := TaskMap[t]; ok {
			Log.WithFields(logrus.Fields{"UID": user.DedeUserID, "Cron": user.Cron, "TaskID": uuid, "TaskName": t}).Info("[CRON] add task")
			_ = c.AddFunc(user.Cron, func() {
				if err := task(user); err != nil {
					Log.WithFields(logrus.Fields{"UID": user.DedeUserID, "TaskID": uuid, "err": err.Error()}).Warn("[CRON] task failed")
					return
				}
				Log.WithFields(logrus.Fields{"UID": user.DedeUserID}).Info("[CRON] task completed")
			})
		}
	}
}

// CronTaskLoad 注册所有投币任务
func CronTaskLoad() {
	// warning: not found user
	bius := GetConfig(true).BiU
	if len(bius) == 0 {
		Log.Warn("[USER] make sure that at least one user exists in bili.json")
		Log.Warn("[USER] tip: use '-n' option to login by bilibili-mobile-client")
		return
	}
	for i := range bius {
		BiliExecutor(&bius[i])
	}
}
