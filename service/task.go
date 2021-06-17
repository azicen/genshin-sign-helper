package service

import (
	"genshin-sign-helper/client"
	"genshin-sign-helper/util"
	log "genshin-sign-helper/util/logger"
)

func Task() {
	g := client.NewGenshinClient()
	err := util.ReadFileAllLine("cookie.txt", func(s string) {
		//log.Debug(s)
		gameRolesList := g.GetUserGameRoles(s)

		for j := 0; j < len(gameRolesList); j++ {
			g.GetSignStateInfo(s, gameRolesList[j])
			g.Sign(s, gameRolesList[j])
		}
	})
	if err != nil {
		log.Error(err.Error())
	}
}
