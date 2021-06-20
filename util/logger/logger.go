package logger

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"

	"genshin-sign-helper/conf"
)

const (
	//LevelError 错误
	LevelError = iota
	//LevelWarning 警告
	LevelWarning
	//LevelInfo 提示
	LevelInfo
	//LevelDebug 除错
	LevelDebug
)

var level int
var detailed bool
var fileIO *bufio.Writer

//Println 打印
func Println(level, content string, v ...interface{}) {
	//处理content中需要替换的值
	count := strings.Count(content, "%v")
	var v0 []interface{}
	if count <= len(v) {
		v0 = v[:count]
		v = v[count:]
	} else {
		v0 = v
		v = nil
	}
	content = fmt.Sprintf(content, v0...)

	//获取详细日志数据
	if detailed {
		_, file, line, ok := runtime.Caller(2)
		if ok {
			content = fmt.Sprintf("%s:%d:%s", file, line, content)
		} else {
			content = fmt.Sprintf("%s:%s", "it was not possible to recover the information", content)
		}
	}

	mag := fmt.Sprintf("%s %s ", level, content)
	if v != nil {
		mag = fmt.Sprintf(mag, v...)
	}
	log.Println(mag)
	if err := fileIO.Flush(); err != nil {
		panic("fileIO ERROR")
	}
}

//Fatal 极端错误
func Fatal(content string, v ...interface{}) {
	if LevelError > level {
		return
	}
	Println("[FATAL]", content, v...)
	panic(content)
	//os.Exit(0)
}

//Error 错误
func Error(content string, v ...interface{}) {
	if LevelError > level {
		return
	}
	Println("[ERROR]", content, v...)
}

//Warning 警告
func Warning(content string, v ...interface{}) {
	if LevelWarning > level {
		return
	}
	Println("[WARN ]", content, v...)
}

//Info 信息
func Info(content string, v ...interface{}) {
	if LevelInfo > level {
		return
	}
	Println("[INFO ]", content, v...)
}

//Debug 校验
func Debug(content string, v ...interface{}) {
	if LevelDebug > level {
		return
	}
	Println("[DEBUG]", content, v...)
}

//BuildLogger 构建logger
func BuildLogger(l string) {
	level = LevelError
	switch l {
	case "error":
		level = LevelError
	case "warning":
		level = LevelWarning
	case "info":
		level = LevelInfo
	case "debug":
		level = LevelDebug
	}
}

func Init(file string) {
	logFile, err := os.OpenFile(file, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		Fatal(err.Error())
	}
	fileIO = bufio.NewWriter(logFile)
	mw := io.MultiWriter(os.Stdout, fileIO)
	log.SetOutput(mw)
	runtime.Caller(1)
	//读取配置文件，设置日志级别
	BuildLogger(conf.LogLevel)
	detailed = conf.LogDetailed
	Info("log file: %v", file)
}
