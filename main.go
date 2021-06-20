package main

import (
	"flag"

	"genshin-sign-helper/conf"
	"genshin-sign-helper/service"
	log "genshin-sign-helper/util/logger"
)

func main() {
	isGo := flag.Bool("go", false, "Sign in immediately.")
	svcFlag := flag.String("service", "",
		"Control the system service. accept\"start\", \"stop\", \"restart\", \"install\", \"uninstall\"")
	flag.Parse()

	prg := service.Init()
	if len(*svcFlag) != 0 {
		prg.Control(*svcFlag)
		return
	}

	Init()

	if *isGo {
		service.Task()
	}

	if err := prg.Instance.Run(); err != nil {
		log.Error(err.Error())
	}
}

func Init() {
	conf.Init()
	log.Init(conf.LogFile)
}
