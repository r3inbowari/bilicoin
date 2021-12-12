package bilicoin

import (
	"bytes"
	"encoding/json"
	"errors"
	. "github.com/r3inbowari/zlog"
	"github.com/sirupsen/logrus"
	"os"
	"regexp"
	"time"
)

/**
 * local config struct
 */
type LocalConfig struct {
	Finger      string     `json:"finger"`    // canvas指纹
	APIAddr     string     `json:"api_addr"`  // API服务ADDR
	CacheTime   time.Time  `json:"-"`         // 缓存时间
	LoggerLevel *string    `json:"log_level"` // 日志等级
	BiU         []BiliUser `json:"biu"`       // 用户集合
	CaKey       string     `json:"ca_key"`    // CA 密钥
	CaCert      string     `json:"ca_crt"`    // CA 证书
}

var config = new(LocalConfig)
var configPath = "bili.json"

// GetConfig 返回配置文件
// imm 立即返回
func GetConfig(imm bool) *LocalConfig {
	if config.CacheTime.Before(time.Now()) || imm {
		if err := LoadConfig(configPath, config); err != nil {
			Log.Fatal("loading file failed")
			return nil
		}
		config.CacheTime = time.Now().Add(time.Second * 60)
	}
	return config
}

// SetConfig set
func (lc *LocalConfig) SetConfig() error {
	fp, err := os.Create(configPath)
	if err != nil {
		Log.WithFields(logrus.Fields{"err": err}).Fatal("loading file failed")
	}
	defer fp.Close()
	data, err := json.Marshal(lc)
	if err != nil {
		Log.WithFields(logrus.Fields{"err": err}).Fatal("marshal file failed")
	}
	n, err := fp.Write(data)
	if err != nil {
		Log.WithFields(logrus.Fields{"err": err}).Fatal("write file failed")
	}
	Log.WithFields(logrus.Fields{"size": n}).Info("[FILE] Update user configuration")
	return nil
}

const configFileSizeLimit = 10 << 20

/**
 * Load File
 * @param path 文件路径
 * @param dist 存放目标
 */
func LoadConfig(path string, dist interface{}) error {
	configFile, err := os.Open(path)
	if err != nil {
		Log.WithFields(logrus.Fields{"path": path, "err": err}).Fatal("Failed to open config file.")
		return err
	}

	fi, _ := configFile.Stat()
	if size := fi.Size(); size > (configFileSizeLimit) {
		Log.WithFields(logrus.Fields{"path": path, "size": size}).Fatal("Config file size exceeds reasonable limited")
		return errors.New("limited")
	}

	if fi.Size() == 0 {
		Log.WithFields(logrus.Fields{"path": path, "size": 0}).Fatal("Config file is empty, skipping")
		return errors.New("empty")
	}

	buffer := make([]byte, fi.Size())
	_, err = configFile.Read(buffer)
	buffer, err = StripComments(buffer)
	if err != nil {
		Log.WithFields(logrus.Fields{"err": err}).Fatal("Failed to strip comments from json")
		return err
	}

	buffer = []byte(os.ExpandEnv(string(buffer)))

	err = json.Unmarshal(buffer, &dist)
	if err != nil {
		Log.WithFields(logrus.Fields{"err": err}).Fatal("Failed unmarshalling json")
		return err
	}
	return nil
}

/**
 * 注释清除
 */
func StripComments(data []byte) ([]byte, error) {
	data = bytes.Replace(data, []byte("\r"), []byte(""), 0)
	lines := bytes.Split(data, []byte("\n"))
	filtered := make([][]byte, 0)

	for _, line := range lines {
		match, err := regexp.Match(`^\s*#`, line)
		if err != nil {
			return nil, err
		}
		if !match {
			filtered = append(filtered, line)
		}
	}
	return bytes.Join(filtered, []byte("\n")), nil
}
