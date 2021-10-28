package bilicoin

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var log = logrus.New()

var LogLevel = map[string]logrus.Level{
	"PANIC": logrus.PanicLevel,
	"FATAL": logrus.FatalLevel,
	"ERROR": logrus.ErrorLevel,
	"WARN":  logrus.WarnLevel,
	"INFO":  logrus.InfoLevel,
	"DEBUG": logrus.DebugLevel,
	"TRACE": logrus.TraceLevel,
}

type Ext struct {
}

func fieldParse(obj interface{}) string {
	var ret string
	switch v := obj.(type) {
	case string:
		ret = v
	case float64:
		ret = strconv.FormatFloat(v, 'E', -1, 64)
	case int:
		ret = strconv.Itoa(v)
	case nil:
		ret = "null"
	default:
		ret = "Unsupported"
	}
	return ret
}

func newEntry(msg string, level logrus.Level) Entry {
	t := time.Now().Format("2006-01-02 15:04:05")
	var str string
	if level == logrus.InfoLevel {
		str = "[INFO] " + t + " " + msg
	} else if level == logrus.WarnLevel {
		str = "[WARN] " + t + " " + msg
	} else if level == logrus.FatalLevel {
		str = "[FATAL] " + t + " " + msg
	}
	return Entry{msg: str, level: level}
}

func (en *Entry) withFields(ext logrus.Fields) {
	en.msg += " | "
	for k, v := range ext {
		en.msg += k + " " + fieldParse(v) + " | "
	}
}

type Entry struct {
	msg   string
	level logrus.Level
}

func (en *Entry) Print() {
	switch en.level {
	case logrus.InfoLevel:
		fmt.Printf("\x1b[%dm"+en.msg+" \x1b[0m\n", 32)
	case logrus.WarnLevel:
		fmt.Printf("\x1b[%dm"+en.msg+" \x1b[0m\n", 33)
	case logrus.FatalLevel:
		fmt.Printf("\x1b[%dm"+en.msg+" \x1b[0m\n", 31)
	}
}

func Info(msg string, ext ...logrus.Fields) {
	en := newEntry(msg, logrus.InfoLevel)
	if len(ext) > 0 {
		en.withFields(ext[0])
	}
	en.Print()
}

func Warn(msg string, ext ...logrus.Fields) {
	en := newEntry(msg, logrus.WarnLevel)
	if len(ext) > 0 {
		en.withFields(ext[0])
	}
	en.Print()
}

func Fatal(msg string, ext ...logrus.Fields) {
	en := newEntry(msg, logrus.FatalLevel)
	if len(ext) > 0 {
		en.withFields(ext[0])
	}
	en.Print()
}

func Blue(msg string) {
	fmt.Printf("\x1b[%dm"+msg+" \x1b[0m\n", 34)
}

func AppInfo(gitHash, buildTime, goVersion string, version string) {
	// RunningMode = mod
	Blue("  ________  ___  ___       ___  ________  ________  ___  ________")
	Blue(" |\\   __  \\|\\  \\|\\  \\     |\\  \\|\\   ____\\|\\   __  \\|\\  \\|\\   ___  \\         BILICOIN #UNOFFICIAL# " + gitHash[:7] + "..." + gitHash[33:])
	Blue(" \\ \\  \\|\\ /\\ \\  \\ \\  \\    \\ \\  \\ \\  \\___|\\ \\  \\|\\  \\ \\  \\ \\  \\\\ \\  \\        -... .. .-.. .. -.-. --- .. -. " + version)
	Blue("  \\ \\   __  \\ \\  \\ \\  \\    \\ \\  \\ \\  \\    \\ \\  \\\\\\  \\ \\  \\ \\  \\\\ \\  \\       Running mode: " + RunningMode)
	if RunningMode == "api" {
		Blue("   \\ \\  \\|\\  \\ \\  \\ \\  \\____\\ \\  \\ \\  \\____\\ \\  \\\\\\  \\ \\  \\ \\  \\\\ \\  \\      Port: " + GetConfig(false).APIAddr[1:])
	} else {
		Blue("   \\ \\  \\|\\  \\ \\  \\ \\  \\____\\ \\  \\ \\  \\____\\ \\  \\\\\\  \\ \\  \\ \\  \\\\ \\  \\      Port: UNSUPPORTED")
	}
	Blue("    \\ \\_______\\ \\__\\ \\_______\\ \\__\\ \\_______\\ \\_______\\ \\__\\ \\__\\\\ \\__\\     PID: " + strconv.Itoa(os.Getpid()))
	Blue("     \\|_______|\\|__|\\|_______|\\|__|\\|_______|\\|_______|\\|__|\\|__| \\|__|     built on " + buildTime)
	Blue("")
}

func InitLogger() {
	log.Out = os.Stdout
	if GetConfig(false).LoggerLevel == nil {
		log.Level = logrus.DebugLevel
	} else {
		log.Level = LogLevel[strings.ToUpper(*GetConfig(false).LoggerLevel)]
	}
	log.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})
	log.Hooks.Add(NewContextHook())
}

//func Info(msg string, fields ...logrus.Fields) {
//	if len(fields) > 0 {
//		log.WithFields(fields[0]).Info(msg)
//	} else {
//		log.Info(msg)
//	}
//}
//
//func Warn(msg string, fields ...logrus.Fields) {
//	if len(fields) > 0 {
//		log.WithFields(fields[0]).Warn(msg)
//	} else {
//		log.Warn(msg)
//	}
//}

func Error(msg string, fields ...logrus.Fields) {
	if len(fields) > 0 {
		log.WithFields(fields[0]).Error(msg)
	} else {
		log.Error(msg)
	}
}

//func Fatal(msg string, fields ...logrus.Fields) {
//	if len(fields) > 0 {
//		log.WithFields(fields[0]).Fatal(msg)
//	} else {
//		log.Fatal(msg)
//	}
//}

func Panic(msg string, fields ...logrus.Fields) {
	if len(fields) > 0 {
		log.WithFields(fields[0]).Panic(msg)
	} else {
		log.Panic(msg)
	}
}

func Trace(msg string, fields ...logrus.Fields) {
	if len(fields) > 0 {
		log.WithFields(fields[0]).Trace(msg)
	} else {
		log.Trace(msg)
	}
}

// ContextHook for log the call context
type contextHook struct {
	Field  string
	Skip   int
	levels []logrus.Level
}

// NewContextHook use to make an hook
// 根据上面的推断, 我们递归深度可以设置到5即可.
func NewContextHook(levels ...logrus.Level) logrus.Hook {
	hook := contextHook{
		Field:  "line",
		Skip:   10,
		levels: levels,
	}
	if len(hook.levels) == 0 {
		hook.levels = logrus.AllLevels
	}
	return &hook
}

// Levels implement levels
func (hook contextHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// Fire implement fire
func (hook contextHook) Fire(entry *logrus.Entry) error {
	entry.Data[hook.Field] = findCaller(hook.Skip)
	return nil
}

// 对caller进行递归查询, 直到找到非logrus包产生的第一个调用.
// 因为filename我获取到了上层目录名, 因此所有logrus包的调用的文件名都是 logrus/...
// 因此通过排除logrus开头的文件名, 就可以排除所有logrus包的自己的函数调用
func findCaller(skip int) string {
	file := ""
	line := 0
	for i := 0; i < 10; i++ {
		file, line = getCaller(skip + i)
		if !strings.HasPrefix(file, "logrus") {
			break
		}
	}
	return fmt.Sprintf("%s:%d", file, line)
}

// 这里其实可以获取函数名称的: fnName := runtime.FuncForPC(pc).Name()
// 但是我觉得有 文件名和行号就够定位问题, 因此忽略了caller返回的第一个值:pc
// 在标准库log里面我们可以选择记录文件的全路径或者文件名, 但是在使用过程成并发最合适的,
// 因为文件的全路径往往很长, 而文件名在多个包中往往有重复, 因此这里选择多取一层, 取到文件所在的上层目录那层.
func getCaller(skip int) (string, int) {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "", 0
	}
	n := 0
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			n++
			if n >= 2 {
				file = file[i+1:]
				break
			}
		}
	}
	return file, line
}
