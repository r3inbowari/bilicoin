package bilicoin

import (
	"bufio"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"time"
)

func Release() {
	// close listener
	Shutdown(context.TODO())
}

func CreateInstallBatch(name string) {
	file, e := os.OpenFile("install.bat", os.O_CREATE|os.O_WRONLY, 0666)
	if e != nil {
		fmt.Println("failed")
		os.Exit(1004)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	writer.WriteString("taskkill /f /pid " + strconv.Itoa(os.Getpid()) + "\n")
	writer.WriteString("start \"bilicoin\" " + name + ".exe -a\n")
	writer.WriteString("exit\n")
	writer.Flush()
}

func ExecBatchFromWindows(path string) error {
	return exec.Command("cmd.exe", "/c", "start "+path+".bat").Start()
}

func Reload(path string) error {
	if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		exec.Command("chmod", "777", path)
		path = "./" + path
	}
	// init接管
	cmd := exec.Command(path, "-a")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Start()
}

var url = "http://r3in.top:3000/StockPankouSampler/PankouSampler_1.0.1"

func DigestVerify(path string, digestStr string) bool {
	if runtime.GOOS == "windows" {
		path += ".exe"
	}
	file, err := os.Open(path)
	if err != nil {
		Info("[UP] file not exist")
		return false
	}
	md5f := md5.New()
	_, err = io.Copy(md5f, file)
	if err != nil {
		Info("[UP] file open error")
		return false
	}

	ok := digestStr == hex.EncodeToString(md5f.Sum([]byte("")))
	if ok {
		Info("[UP] file digest match", logrus.Fields{"digest": digestStr, "file": hex.EncodeToString(md5f.Sum([]byte("")))})
	} else {
		Info("[UP] file digest mismatch", logrus.Fields{"digest": digestStr, "file": hex.EncodeToString(md5f.Sum([]byte("")))})
	}
	return ok
}

func CheckAndUpdateAndReload() {
	Warn("[UP] Checking for updates")
	defer func() {
		Info("[UP] update check completed")
	}()

	systemType := runtime.GOOS

	name := "bilicoin"
	// systemType = "linux"

	// get md5 digest
	ok, digest, verStr := CheckUpdate()
	if !ok {
		return
	}
	// var digest = "fba41dcef7634ed0b6c92a22c32ea2f8"

	// download
	DownloadExec("bilicoin", verStr)

	// verify
	verify := DigestVerify("bilicoin", digest)
	if !verify {
		return
	}

	// reload
	Info("[UP] reloading", logrus.Fields{"os": runtime.GOOS, "arch": runtime.GOARCH})
	if systemType == "linux" || systemType == "darwin" {
		Release()
		time.Sleep(time.Second * 3)
		Reload(name)
		os.Exit(1010)
	} else if systemType == "windows" {
		CreateInstallBatch(name)
		ExecBatchFromWindows("install")
	}
}

var host = "http://r3in.top:3000/"

func DownloadExec(name, version string) error {

	goarch := runtime.GOARCH
	goos := runtime.GOOS
	dUrl := host + name + "/bin/" + name + "_" + goos + "_" + goarch + "_" + version

	if runtime.GOOS == "windows" {
		dUrl += ".exe"
		name += ".exe"
	}
	Info(dUrl)

	var bar *ProgressBar

	err := Download(dUrl, name, func(fileLength int64) {
		Info("[UP] downloading... collected file size", logrus.Fields{"size": fileLength})
		bar = NewProgressBar(fileLength)
	}, func(length, downLen int64) {
		bar.Play(downLen)
	})
	if err != nil {
		Warn("[UP] download failed...")
		return err
	}
	bar.Finish()
	return nil
}

func Download(url, name string, lenCall func(fileLength int64), fb func(length, downLen int64)) error {
	var (
		fsize   int64
		buf     = make([]byte, 32*1024)
		written int64
	)
	client := new(http.Client)
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	fsize, err = strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 32)
	if err != nil {
		fmt.Println(err)
	}

	lenCall(fsize)

	file, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		return err
	}
	defer file.Close()
	if resp.Body == nil {
		return errors.New("null")
	}
	defer resp.Body.Close()
	for {
		nr, er := resp.Body.Read(buf)
		if nr > 0 {
			nw, ew := file.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
		fb(fsize, written)
	}
	return err
}

type Default struct {
	Name    string   `json:"name"`
	Major   int      `json:"major"`
	Minor   int      `json:"minor"`
	Patch   int      `json:"patch"`
	Types   []string `json:"types"`
	Digests []string `json:"digests"`
}

func CheckUpdate() (bool, string, string) {

	res, err := http.Get(host + "bilicoin/bin/default.json")
	if err != nil {
		return false, "", ""
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return false, "", ""
	}

	var defs Default
	err = json.Unmarshal(body, &defs)
	if err != nil {
		return false, "", ""
	}

	Info("[UP] Current version", logrus.Fields{"major": version.Major, "minor": version.Minor, "patch": version.Patch})
	value := defs.Major<<24 + defs.Minor<<12 + defs.Patch<<0
	now := version.Major<<24 + version.Minor<<12 + version.Patch<<0
	if now < int64(value) {
		Info("[UP] Found new version", logrus.Fields{"major": defs.Major, "minor": defs.Minor, "patch": defs.Patch})
		for k, v := range defs.Types {
			if v == runtime.GOOS+"_"+runtime.GOARCH {
				return true, defs.Digests[k], "v" + strconv.FormatInt(int64(defs.Major), 10) + "." + strconv.FormatInt(int64(defs.Minor), 10) + "." + strconv.FormatInt(int64(defs.Patch), 10)
			}
		}
	}
	return false, "", ""
}
