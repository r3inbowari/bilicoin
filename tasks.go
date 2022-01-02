package bilicoin

import (
	"errors"
	"time"
)

var TaskMap = map[string]Task{
	"drop-coin":     TaskDropCoin,
	"silver-2-coin": TaskSilver2Coin,
}

func TaskSilver2Coin(user *BiliUser) error {
	// 获取用户信息失败
	if err := user.GetBiliWallet(); err != nil {
		return err
	}
	// 使用银瓜子兑换硬币一枚
	if user.ConvertCoin && user.Bi.Silver >= 700 {
		return user.Silver2Coin()
	}
	return errors.New("not enough silver")
}

func TaskDropCoin(user *BiliUser) error {
	// 获取日志失败
	// 过期、未知错误、服务不可达、解析错误
	if err := user.GetBiliCoinLog(); err != nil {
		return err
	}
	for true {
		if user.DropCoinCount > 4 {
			user.InfoUpdate()
			break
		}
		user.RandDrop()
		time.Sleep(Random(60))
	}
	return nil
}
