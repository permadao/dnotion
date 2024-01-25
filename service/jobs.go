package service

import (
	"time"
)

func (s *Service) runJobs() {
	var err error
	s.scheduler, err = gocron.NewScheduler(gocron.WithLocation(time.UTC))
	if err != nil {
		panic(err)
	}

	// job list
	s.scheduler.NewJob(
		gocron.WeeklyJob(1, gocron.NewWeekdays(time.Saturday), gocron.NewAtTimes(gocron.NewAtTime(1, 0, 0))),
		gocron.NewTask(s.genGrade),
	)
	s.scheduler.NewJob(
		gocron.WeeklyJob(1, gocron.NewWeekdays(time.Saturday), gocron.NewAtTimes(gocron.NewAtTime(1, 30, 0))),
		gocron.NewTask(s.genGuilds),
	)

	s.scheduler.Start()
}

func (s *Service) genGuilds() {
	date := GetPreviousDate(1)
	log.Info("genGuild...", "date", date)

	s.guild.GenGuilds("AR", date)

	log.Info("genGuild done")
}

func (s *Service) genGrade() {
	end := GetCurrentDate()
	last := GetPreviousDate(7)
	start := GetPreviousDate(4 * 7)
	log.Info("genGrade...", "start", start, "end", end)

	// translation guild grade
	if err := s.guild.GenGrade("e8d79c55c0394cba83664f3e5737b0bd", "d8c270f68a8f44aaa6b24e17c927df2b", start, end); err != nil {
		log.Error("genGrade failed", "err", err)
	}

	// developer guild grade
	if err := s.guild.GenDevGrade("146e1f661ed943e3a460b8cf12334b7b", "623ccfc9fb1443279decf90fb752215d", last, end); err != nil {
		log.Error("genDevGrade failed", "err", err)
	}

	log.Info("genGrade done")
}

func GetCurrentDate() (date string) {
	now := time.Now()
	date = now.Format("2006-01-02")
	return
}

func GetPreviousDate(days int) (date string) {
	now := time.Now()
	last := now.AddDate(0, 0, -days)
	date = last.Format("2006-01-02")
	return
}
