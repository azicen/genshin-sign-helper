package client

import (
	"fmt"

	"genshin-sign-helper/model"
	"genshin-sign-helper/util/constant"
	log "genshin-sign-helper/util/logger"
)

//GetUserGameRoles 获取用户游戏角色
func (g *GenshinClient) GetUserGameRoles(cookie string) (rolesList []model.GameRolesData) {
	var info model.GameRolesInfo
	err := g.SendGetMessage(
		cookie,
		constant.GetUserGameRolesByCookie,
		"game_biz=hk4e_cn",
		false,
		&info,
	)
	if err != nil {
		log.Error("unable to send http massage.", err)
		return nil
	}

	switch info.Code {
	case 0:
		log.Debug("get user game roles success.")
		break
	default:
		log.Error("get user game roles error(%v). request failure.", info.Code, info.Msg, info.Data)
		break
	}
	return info.Data.List
}

//GetSignStateInfo 获取签到信息
func (g *GenshinClient) GetSignStateInfo(cookie string, roles model.GameRolesData) (data *model.SignStateData) {
	var info model.SignStateInfo
	err := g.SendGetMessage(
		cookie,
		constant.GetBbsSignRewardInfo,
		fmt.Sprintf("act_id=%s&region=%s&uid=%s", constant.ActID, roles.Region, roles.UID),
		false,
		&info,
	)
	if err != nil {
		log.Error("unable to send http massage.", err)
		return nil
	}
	switch info.Code {
	case 0:
		log.Debug("get sign reward info success.")
		break
	default:
		log.Error("get sign reward info error(%v). request failure.", info.Code, info.Msg, info.Data)
		break
	}
	return &info.Data
}

func (g *GenshinClient) Sign(cookie string, roles model.GameRolesData) bool {
	data := map[string]string{
		"act_id": constant.ActID,
		"region": roles.Region,
		"uid":    roles.UID,
	}
	var info model.SignResponseInfo
	err := g.SendPostMessage(cookie, constant.PostSignInfo, "", true, data, &info)
	if err != nil {
		log.Error("unable to send http massage.", err)
		return false
	}
	switch info.Code {
	case 0:
		log.Info("roles(%v:%v) sign success.", roles.UID, roles.Name)
		return true
	case -5003:
		log.Info("roles(%v:%v) sign info(%v). request failure. %v", roles.UID, roles.Name, info.Code, info.Msg, info.Data)
		return false
	default:
		log.Error("roles(%v:%v) sign error(%v). request failure. %v", roles.UID, roles.Name, info.Code, info.Msg, info.Data)
		return false
	}
}
