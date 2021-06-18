package conf

import (
	"bufio"
	_ "embed"
	"encoding/json"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"

	"genshin-sign-helper/model"
	"genshin-sign-helper/util"
)

var (
	LogLevel    string // 显示日志等级
	LogFile     string // 日志文件位置
	LogDetailed bool   // 是否显示详细的日志内容
	Cycle       int    // 签到周期（小时）
	SignTime    int    // 签到时间，用于确定几点签到，如果当前运行时间大于签到时间，且当日未签到，则立刻签到
)

var SignRecordJSON *model.SignRecordJSON = nil

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
	if !util.CheckFileIsExist("record.json") {
		create, err := os.Create("record.json")
		if err != nil {
			panic(err)
		}
		if err := create.Close(); err != nil {
			panic(err)
		}
	}
	body, err := util.ReadFile("record.json")
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

func SendRecordJSON() (err error) {
	file, err := os.OpenFile("record.json", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
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
