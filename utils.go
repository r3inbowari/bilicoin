package bilicoin

import (
	"encoding/json"
	"errors"
	"fmt"
	qrcodeTerminal "github.com/Baozisoftware/qrcode-terminal-go"
	"github.com/robertkrimen/otto"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var inlineJSCode = `function generateUuidPart(c){for(var a="",b=0;b<c;b++){a+=parseInt(16*Math.random()).toString(16).toUpperCase()}return formatNum(a,c)}function generateUuid(){var d=generateUuidPart(8),b=generateUuidPart(4),c=generateUuidPart(4),g=generateUuidPart(4),f=generateUuidPart(12),a=(new Date).getTime();return d+"-"+b+"-"+c+"-"+g+"-"+f+formatNum((a%100000).toString(),5)+"infoc"}function formatNum(c,a){var b="";if(c.length<a){for(var d=0;d<a-c.length;d++){b+="0"}}return b+c};`
var xor int64 = 177451812
var add int64 = 8728348608
var table = "fZodR9XQDSUm21yCkr6zBqiveYah8bt4xsWpHnJE7jL5VG3guMTKNPAwcF"
var s = []int64{11, 10, 3, 8, 4, 6, 2, 9, 5, 7}

// bv2av
func BVCovertDec(bv string) string {
	var tr = make(map[byte]int)
	for i := 1; i < 58; i++ {
		tr[table[i]] = i
	}
	var r int64 = 0
	var i int64 = 0
	for ; i < 6; i++ {
		k := int64(math.Pow(58.0, float64(i)))
		r += int64(tr[bv[s[i]]]) * k
	}
	retAV := (r - add) ^ xor
	return strconv.FormatInt(retAV, 10)
}

// av2bv
func BVCovertEnc(av string) string {
	var r = []string{"B", "V", "1", "", "", "4", "", "1", "", "7", "", ""}
	x, _ := strconv.Atoi(av)
	x64 := int64(x)
	x_ := (x64 ^ xor) + add
	for i := 0; i < 6; i++ {
		_x := int64(math.Pow(58.0, float64(i)))
		aj := int64(math.Floor(float64(x_ / _x)))
		r[s[i]] = string(table[aj%58])
	}
	return strings.Join(r, "")
}

// buvid 生成
func _buvidGenerate() (string, error) {
	url := "https://data.bilibili.com/v/web/web_page_view?mid=null&fts=null&url=https%253A%252F%252Fwww.bilibili.com%252F&proid=3&ptype=2&module=game&title=%E5%93%94%E5%93%A9%E5%93%94%E5%93%A9%20(%E3%82%9C-%E3%82%9C)%E3%81%A4%E3%83%AD%20%E5%B9%B2%E6%9D%AF~-bilibili&ajaxtag=&ajaxid=&page_ref=https%253A%252F%252Fpassport.bilibili.com%252Flogin"
	res, err := GET(url, nil)
	if res != nil {
		cook := res.Cookies()
		return cook[1].Value, err
	}
	return "", errors.New("nil buvid response")
}

// 鬼畜区 BVS
func GetGuichuBVs() []string {
	res, _ := GET("https://api.bilibili.com/x/web-interface/ranking/region?rid=119&day=3&original=0", nil)
	result, _ := ioutil.ReadAll(res.Body)
	reg := regexp.MustCompile("BV[a-zA-Z0-9_]{10}")
	return reg.FindAllString(string(result), -1)
}

// _uuid 生成
func _uuidGenerate() (string, error) {
	vm := otto.New()
	_, err := vm.Run(inlineJSCode)
	uuid, err := vm.Call("generateUuid", nil)
	return uuid.String(), err
}

// QRCode 打印
func QRCPrint(content string) {
	obj := qrcodeTerminal.New2(qrcodeTerminal.ConsoleColors.BrightBlue, qrcodeTerminal.ConsoleColors.BrightGreen, qrcodeTerminal.QRCodeRecoveryLevels.Low)
	obj.Get(content).Print()
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

func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func IsFile(path string) bool {
	return !IsDir(path)
}

func InitConfig() {
	Info("[FILE] Init user configuration")
	if !Exists("bili.json") {
		var config LocalConfig
		var l = "debug"
		config.Finger = "1777945899"
		config.LoggerLevel = &l
		config.BiU = []BiliUser{}
		config.APIAddr = ":9090"
		_ = config.SetConfig()
	}
}

func GetRequest(url string) {
	_, err := http.Get(url)
	if err != nil {
		log.Println("[FAIL] 方糖调用错误")
	}
}

func GET(url string, interceptor func(reqPoint *http.Request)) (*http.Response, error) {
	method := "GET"

	client := &http.Client{}

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	if interceptor != nil {
		interceptor(req)
	}
	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.105 Safari/537.36")
	res, err := client.Do(req)
	return res, nil
}

func Post2(url string, interceptor func(reqPoint *http.Request), body string) (*http.Response, error) {
	method := "POST"

	client := &http.Client{}
	payload := strings.NewReader(body)
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return nil, err
	}

	if interceptor != nil {
		interceptor(req)
	}
	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.105 Safari/537.36")
	res, err := client.Do(req)
	return res, err
}

func Post(url string, interceptor func(reqPoint *http.Request)) (*http.Response, error) {
	method := "POST"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	if interceptor != nil {
		interceptor(req)
	}
	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.105 Safari/537.36")
	res, err := client.Do(req)
	return res, err
}

func Random(max int) time.Duration {
	rand.Seed(time.Now().UnixNano())
	return time.Duration(rand.Intn(max)) * time.Second
}

type RequestResult struct {
	Total   int         `json:"total"`
	Data    interface{} `json:"data"`
	Code    int         `json:"code"`
	Message string      `json:"msg"`
}

func ResponseCommon(w http.ResponseWriter, data interface{}, msg string, total int, tag int, code int) {
	var rq RequestResult
	rq.Data = data
	rq.Total = total
	rq.Code = code
	rq.Message = msg
	jsonStr, err := json.Marshal(rq)
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	w.WriteHeader(tag)
	_, _ = fmt.Fprintf(w, string(jsonStr))
}

func CreateUUID() string {
	u1 := uuid.NewV4()
	return u1.String()
}
