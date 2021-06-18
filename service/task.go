package service

import (
	"time"

	"genshin-sign-helper/client"
	"genshin-sign-helper/conf"
	"genshin-sign-helper/model"
	"genshin-sign-helper/util"
	log "genshin-sign-helper/util/logger"
)

func Task() {
	log.Info("开始签到任务...")

	g := client.NewGenshinClient()
	err := util.ReadFileAllLine("cookie.txt", func(s string) {
		//log.Debug(s)
		gameRolesList := g.GetUserGameRoles(s)
		currentDay := time.Now().Day()

		for j := 0; j < len(gameRolesList); j++ {
			time.Sleep(5 * time.Second)
			if currentDay == conf.SignRecordJSON.Roles[gameRolesList[j].UID].Time.Day() {
				continue
			}
			if g.Sign(s, gameRolesList[j]) {
				data := g.GetSignStateInfo(s, gameRolesList[j])
				log.Info("UID:%v 昵称:%v 连续签到天数:%v 签到成功.",
					gameRolesList[j].UID, gameRolesList[j].Name, data.TotalSignDay)
				conf.SignRecordJSON.Roles[gameRolesList[j].UID] = model.RolesJSON{
					UID:          gameRolesList[j].UID,
					Name:         gameRolesList[j].Name,
					Time:         time.Now(),
					TotalSignDay: data.TotalSignDay,
				}
			} else {
				log.Info("UID:%v 昵称:%v 签到失败.",
					gameRolesList[j].UID, gameRolesList[j].Name)
			}
		}

		conf.SignRecordJSON.Time = time.Now()
		err := conf.SendRecordJSON()
		if err != nil {
			return
		}
	})
	if err != nil {
		log.Error(err.Error())
	}

	log.Info("签到结束.")
}
