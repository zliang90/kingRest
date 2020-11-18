package main

import (
	"context"
	"flag"
	"github.com/zliang90/kingRest/internal/app/conf"
	"github.com/zliang90/kingRest/internal/app/db"
	restApi "github.com/zliang90/kingRest/internal/restful"
	"github.com/zliang90/kingRest/pkg/log"
	"os"
	"os/signal"
	"runtime"
	"time"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)

	cfgPath := flag.String("c", "config/config.yaml", "configuration file path")
	flag.Parse()

	var err error

	// config
	if err = conf.LoadConfig(*cfgPath); err != nil {
		log.Fatal(err)
		return
	}
	log.Infof("load config: %s", conf.GetConfig().String())
	// log level
	log.SetLevel(conf.GetConfig().LogLevel)
	// db
	log.Info("init db connections")
	if err := db.Init(conf.GetConfig()); err != nil {
		log.Fatal(err)
		return
	}

	// root context
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
		time.Sleep(2 * time.Second)

		log.Info("close db connections")
		db.Close()
	}()

	// restful api
	go restApi.New(conf.GetConfig()).Run(ctx)

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, os.Kill)

	<-quit
}
