package main

import (
	"fmt"
	"os"

	"github.com/permadao/dnotion/config"
	"github.com/permadao/dnotion/db"
	"github.com/permadao/dnotion/fin"
	"github.com/permadao/dnotion/service"
	log "github.com/sirupsen/logrus"
)

func main() {
	// 1. Read configs
	config.Init("config")

	// 2. init log
	initLog()

	// 3. init db
	db.Init()

	// 4. init finance
	fin.Init()

	// last init Service
	service.StartServe()
}

func initLog() {
	// Init log
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	// save log file
	if config.Config.Log.Save {
		logfile, err := os.OpenFile(config.Config.Log.File, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic(fmt.Sprintf("open log file error: %s", err.Error()))
		}
		log.SetOutput(logfile)
	}

	log.SetLevel(log.DebugLevel)
}
