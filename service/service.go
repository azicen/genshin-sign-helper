package service

import (
    "genshin-sign-helper/conf"
    "genshin-sign-helper/util"
    "time"

    "github.com/kardianos/service"

    log "genshin-sign-helper/util/logger"
)

// Service 服务
// 程序结构。
// 定义启动和停止方法。
type Service struct {
    Config       *service.Config
    Instance     service.Service
    exit         chan struct{}
    i            int
    startingTime time.Time
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

    if err != nil {
        log.Fatal(err.Error())
    }

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
    //Init()

    if service.Interactive() {
        log.Info("GenshinImpact Sign Helper Service Running in terminal.")
    } else {
        log.Info("GenshinImpact Sign Helper Service Running under service manager.")
    }
    p.startingTime = time.Now()
    p.exit = make(chan struct{})

    // 开始不应该阻塞。异步执行实际工作。
    go func() {
        Task()
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
            log.Debug("go i:%v", p.i)
            currentTime := time.Now()
            if util.GetCurrentDayTimestamp() > util.GetDayTimestamp(conf.SignRecordJSON.Time) {
                if currentTime.Hour() == 23 {
                    log.Debug("大于23时，且今日为签到，立刻签到 i:%v", p.i)
                    p.i = 0
                    Task()
                } else if currentTime.Hour() > conf.SignTime {
                    log.Debug("已到设置签到时间，开始签到 i:%v", p.i)
                    p.i = 0
                    Task()
                }
                continue
            }
            if p.i < conf.Cycle {
                p.i++
                continue
            }
            log.Debug("定时签到任务开始 i:%v", p.i)
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
    t := time.Since(p.startingTime)
    h := t / time.Hour
    t -= h * time.Hour
    m := t / time.Minute
    log.Info("Total working time %vH %vM", int64(h), int64(m))
    log.Info("GenshinImpact Sign Helper Service Stopping!")
    close(p.exit)
    return nil
}
