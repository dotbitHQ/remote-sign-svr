package main

import (
	"context"
	"fmt"
	"github.com/dotbitHQ/das-lib/http_api"
	"github.com/dotbitHQ/das-lib/http_api/logger"
	"github.com/scorpiotzh/toolib"
	"github.com/urfave/cli/v2"
	"os"
	"remote-sign-svr/config"
	"remote-sign-svr/dao"
	"remote-sign-svr/http_svr"
	"remote-sign-svr/http_svr/handle"
	"sync"
	"time"
)

var (
	log               = logger.NewLogger("main", logger.LevelDebug)
	exit              = make(chan struct{})
	ctxServer, cancel = context.WithCancel(context.Background())
	wgServer          = sync.WaitGroup{}
)

func main() {
	log.Debugf("server start：")
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Load configuration from `FILE`",
			},
		},
		Action: runServer,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func runServer(ctx *cli.Context) error {
	// config file
	configFilePath := ctx.String("config")
	if err := config.InitCfg(configFilePath); err != nil {
		return err
	}

	// config file watcher
	watcher, err := config.AddCfgFileWatcher(configFilePath)
	if err != nil {
		return err
	}
	// ============= service start =============

	//sentry
	if err := http_api.SentryInit(config.Cfg.Notify.SentryDsn); err != nil {
		return fmt.Errorf("SentryInit err: %s", err.Error())
	}
	defer http_api.RecoverPanic()
	// db
	dbDao, err := dao.NewGormDB(config.Cfg.DB.Mysql)
	if err != nil {
		return fmt.Errorf("dao.NewGormDB err: %s", err.Error())
	}

	// http
	httpSvr := http_svr.HttpSvr{
		Ctx:  ctxServer,
		Addr: config.Cfg.Server.HttpAddr,
		H: &handle.HttpHandle{
			Ctx:   ctxServer,
			DbDao: dbDao,
		},
	}
	httpSvr.Run()

	// ============= service end =============
	toolib.ExitMonitoring(func(sig os.Signal) {
		log.Warn("ExitMonitoring:", sig.String())
		if watcher != nil {
			log.Warn("close watcher ... ")
			_ = watcher.Close()
		}
		httpSvr.Shutdown()
		cancel()
		wgServer.Wait()
		log.Warn("success exit server. bye bye!")
		time.Sleep(time.Second)
		exit <- struct{}{}
	})

	<-exit
	return nil
}
