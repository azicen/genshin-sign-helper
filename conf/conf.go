package conf

import (
	"bufio"
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"

	"genshin-sign-helper/model"
	"genshin-sign-helper/util"
	"genshin-sign-helper/util/constant"
)

var (
	LogLevel    string // 显示日志等级
	LogDetailed bool   // 是否显示详细的日志内容
	Cycle       int    // 签到周期（小时）
	SignTime    int    // 签到时间，用于确定几点签到，如果当前运行时间大于签到时间，且当日未签到，则立刻签到
)

var (
	RunDir     string // 运行目录
	EnvFile    string // .env目录
	RecordFile string // record.json目录
	LogFile    string // log.txt目录
	CookieFile string // cookie.txt目录
)

var SignRecordJSON *model.SignRecordJSON = nil

//go:embed conf.env
var confFile []byte

func Init() {
	InitFile()

	SendEnv()

	ReadRecordJSON()
}

func InitFile() {
	RunDir = util.GetCurrentDir()
	EnvFile = fmt.Sprintf("%v/%v", RunDir, constant.EnvFileName)
	RecordFile = fmt.Sprintf("%v/%v", RunDir, constant.RecordFileName)
	LogFile = fmt.Sprintf("%v/%v", RunDir, constant.LogFileName)
	CookieFile = fmt.Sprintf("%v/%v", RunDir, constant.CookieFileName)

	if !util.CheckFileIsExist(EnvFile) {
		osFile, err := os.OpenFile(EnvFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			panic(err)
		}
		if _, err = osFile.Write(confFile); err != nil {
			panic(err)
		}
	}
	if !util.CheckFileIsExist(CookieFile) {
		create, err := os.Create(CookieFile)
		if err != nil {
			panic(err)
		}
		if err := create.Close(); err != nil {
			panic(err)
		}
	}
	if !util.CheckFileIsExist(RecordFile) {
		_, err := os.Create(RecordFile)
		if err != nil {
			panic(err)
		}
	}
}

//SendEnv 读取.env配置文件
func SendEnv() {
	err := godotenv.Load(EnvFile)
	if err != nil {
		panic(err)
	}

	LogLevel = os.Getenv("LOG_LEVEL")

	//var err error
	LogDetailed, err = strconv.ParseBool(os.Getenv("LOG_DETAILED"))
	if err != nil {
		LogDetailed = false
	}

	cycle64, err := strconv.ParseInt(os.Getenv("CYCLE"), 10, 8)
	if err != nil {
		Cycle = 1
	}
	Cycle = int(cycle64)

	signTime64, err := strconv.ParseInt(os.Getenv("SIGN_TIME"), 10, 8)
	if err != nil {
		SignTime = 7
	}
	SignTime = int(signTime64)
}

//ReadRecordJSON 读取保存签到记录的json
func ReadRecordJSON() {
	body, err := util.ReadFile(RecordFile)
	if err == nil {
		if err := json.Unmarshal(body, &SignRecordJSON); err != nil {
			SignRecordJSON = nil
		}
	}
	if SignRecordJSON == nil {
		SignRecordJSON = &model.SignRecordJSON{
			Time:  time.Unix(1600131600, 0), //1601254800
			Roles: map[string]model.RolesJSON{},
		}
	}
}

//SendRecordJSON 储存保存签到记录的json
func SendRecordJSON() (err error) {
	file, err := os.OpenFile(RecordFile, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	toJSON, err := util.StructToJSON(SignRecordJSON)
	if err != nil {
		return err
	}

	fileIO := bufio.NewWriter(file)
	_, err = fileIO.Write(toJSON)
	err = fileIO.Flush()
	if err != nil {
		return err
	}
	return err
}
