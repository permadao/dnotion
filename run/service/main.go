package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/permadao/dnotion/config"
	"github.com/permadao/dnotion/db"
	"github.com/permadao/dnotion/guild"
	"github.com/permadao/dnotion/service"
)

func main() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	config := config.New("./config.toml")
	db := db.New(config)
	g := guild.New(config, db)

	s := service.New(g)
	s.Run()

	<-signals
	s.Close()
}
