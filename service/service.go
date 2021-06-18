package service

import (
	"genshin-sign-helper/conf"
	"time"

	"github.com/kardianos/service"

	log "genshin-sign-helper/util/logger"
)

// Service 服务
// 程序结构。
// 定义启动和停止方法。
type Service struct {
	Config   *service.Config
	Instance service.Service
	exit     chan struct{}
	i        int
}

func Init() (prg *Service) {
	prg = &Service{
		Config: &service.Config{
			Name:        "GenshinSignHelperService",
			DisplayName: "GenshinImpact Sign Helper Service",
			Description: "GenshinImpact mihoyo community sign helper.",
		},
		i: 0,
	}
	s, err := service.New(prg, prg.Config)
	if err != nil {
		log.Fatal(err.Error())
	}
	prg.Instance = s
	errs := make(chan error, 5)

	if err != nil {
		log.Fatal(err.Error())
	}

	go func() {
		for {
			err := <-errs
			if err != nil {
				log.Info(err.Error())
			}
		}
	}()
	return
}

func (p *Service) Control(svcFlag string) {
	if len(svcFlag) != 0 {
		err := service.Control(p.Instance, svcFlag)
		if err != nil {
			log.Fatal("Valid actions: %q\n", service.ControlAction, err)
		}
		return
	}
}

func (p *Service) Start(s service.Service) error {
	if service.Interactive() {
		log.Info("GenshinImpact Sign Helper Service Running in terminal.")
	} else {
		log.Info("GenshinImpact Sign Helper Service Running under service manager.")
	}
	p.exit = make(chan struct{})

	// 开始不应该阻塞。异步执行实际工作。
	go func() {
		err := p.run()
		if err != nil {
			log.Fatal(err.Error())
		}
	}()
	return nil
}

func (p *Service) run() error {
	// 运行间隔
	ticker := time.NewTicker(time.Hour)
	for {
		select {
		case <-ticker.C:
			currentTime := time.Now()
			if currentTime.Hour() == 23 && currentTime.Day() > conf.SignRecordJSON.Time.Day() {
				Task()
				continue
			}
			if conf.SignTime > currentTime.Hour() {
				continue
			}
			if p.i < conf.Cycle {
				p.i++
				continue
			}
			p.i = 0
			Task()
			break
		case <-p.exit:
			ticker.Stop()
			return nil
		}
	}
}

func (p *Service) Stop(s service.Service) error {
	// Stop 中的任何工作都应该很快，通常最多几秒钟。
	log.Info("GenshinImpact Sign Helper Service Stopping!")
	close(p.exit)
	return nil
}
