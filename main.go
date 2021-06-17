package main

import (
	"flag"
	"genshin-sign-helper/conf"
	"genshin-sign-helper/service"
	log "genshin-sign-helper/util/logger"
)

func main() {
	conf.Init()

	log.Init()

	isGo := flag.Bool("go", true, "Sign in immediately.")
	svcFlag := flag.String("service", "", "Control the system service.")
	flag.Parse()

	if *isGo {
		service.Task()
	}

	prg := service.Init()
	prg.Control(*svcFlag)

	if err := prg.Instance.Run(); err != nil {
		log.Error(err.Error())
	}
}
