package service

import (
	"github.com/go-co-op/gocron/v2"
	"github.com/permadao/dnotion/guild"
	"github.com/permadao/dnotion/logger"
)

var log = logger.New("service")

type Service struct {
	scheduler gocron.Scheduler

	guild *guild.Guild
}

func New(g *guild.Guild) *Service {
	return &Service{
		guild: g,
	}
}

func (s *Service) Run() {
	log.Info("service running")

	go s.runJobs()
}

func (s *Service) Close() {
	s.scheduler.Shutdown()

	log.Info("service closed")
}
