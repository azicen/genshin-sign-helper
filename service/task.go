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

    i := 1

    g := client.NewGenshinClient()
    err := util.ReadFileAllLine(conf.CookieFile, func(s string) {

        time.Sleep(time.Duration(conf.SignInterval) * time.Minute) //SIGN_INTERVAL

        log.Debug(s)
        log.Info("执行cookie行号: %v", i)

        gameRolesList := g.GetUserGameRoles(s)
        currentDay := time.Now().Day()

        for j := 0; j < len(gameRolesList); j++ {
            time.Sleep(10 * time.Second)

            if signTime, ok := conf.SignRecordJSON.Roles[gameRolesList[j].UID]; ok {
                if currentDay == signTime.Time.Day() {
                    log.Debug("Line: %v, UID:%v. 今日已签到.", i, gameRolesList[j].UID)
                    continue
                }
            }

            if g.Sign(s, gameRolesList[j]) {
                time.Sleep(10 * time.Second)
                data := g.GetSignStateInfo(s, gameRolesList[j])
                log.Info("Line: %v, UID:%v, 昵称:%v, 连续签到天数:%v. 签到成功.",
                    i, gameRolesList[j].UID, gameRolesList[j].Name, data.TotalSignDay)
                conf.SignRecordJSON.Roles[gameRolesList[j].UID] = model.RolesJSON{
                    UID:          gameRolesList[j].UID,
                    Name:         gameRolesList[j].Name,
                    Time:         time.Now(),
                    TotalSignDay: data.TotalSignDay,
                }
            } else {
                log.Error("Line: %v, UID:%v, 昵称:%v. 签到失败.",
                    i, gameRolesList[j].UID, gameRolesList[j].Name)
            }
        }

        i++

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
