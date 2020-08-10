package bilicoin

import (
	"encoding/json"
	qrcodeTerminal "github.com/Baozisoftware/qrcode-terminal-go"
	"net/http"
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

func QRCodePrint(content string) {
	obj := qrcodeTerminal.New2(qrcodeTerminal.ConsoleColors.BrightBlue, qrcodeTerminal.ConsoleColors.BrightGreen, qrcodeTerminal.QRCodeRecoveryLevels.Low)
	obj.Get(content).Print()
}

func GetQRCode() error {
	res, err := API("https://passport.bilibili.com/qrcode/getLoginUrl", func(reqPoint *http.Request) {

		cookie1 := &http.Cookie{Name: "_uuid", Value: "E874B060-657B-4131-91EF-53E11ED6AC1792382infoc"}
		cookie2 := &http.Cookie{Name: "sid", Value: "ix1spr4e"}
		cookie3 := &http.Cookie{Name: "buvid3", Value: "C3D5DE5F-F2DA-4C7D-93FB-70AFBD1C1A0C143080infoc"}
		cookie4 := &http.Cookie{Name: "PVID", Value: "1"}

		reqPoint.AddCookie(cookie1)
		reqPoint.AddCookie(cookie2)
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

	QRCodePrint(qr.Data.Url)
	return nil
}
