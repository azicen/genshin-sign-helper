package conf

import (
	_ "embed"
	"genshin-sign-helper/util"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	LogLevel    string // 显示日志等级
	LogFile     string // 日志文件位置
	LogDetailed bool   // 是否显示详细的日志内容
	Cycle       int64  // 签到周期（小时）
)

//go:embed conf.env
var confFile []byte

func Init() {
	InitFile()

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	LogLevel = os.Getenv("LOG_LEVEL")
	LogFile = os.Getenv("LOG_FILE")
	//var err error
	LogDetailed, err = strconv.ParseBool(os.Getenv("LOG_DETAILED"))
	if err != nil {
		LogDetailed = false
	}
	Cycle, err = strconv.ParseInt(os.Getenv("CYCLE"), 10, 8)
	if err != nil {
		Cycle = 1
	}
}

func InitFile() {
	if !util.CheckFileIsExist(".env") {
		osFile, err := os.OpenFile(".env", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			panic(err)
		}
		if _, err = osFile.Write(confFile); err != nil {
			panic(err)
		}
	}
	if !util.CheckFileIsExist("cookie.txt") {
		create, err := os.Create("cookie.txt")
		if err != nil {
			panic(err)
		}
		if err := create.Close(); err != nil {
			panic(err)
		}
	}
}
