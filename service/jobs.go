package service

import (
	"fmt"
	"time"

	"github.com/go-co-op/gocron/v2"
)

func (s *Service) runJobs() {
	var err error
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		panic(err)
	}

	fmt.Println("using time zone (UTC + 8):", time.Now().In(loc))

	s.scheduler, err = gocron.NewScheduler(gocron.WithLocation(loc))
	if err != nil {
		panic(err)
	}

	// job list
	s.scheduler.NewJob(
		gocron.WeeklyJob(1, gocron.NewWeekdays(time.Saturday), gocron.NewAtTimes(gocron.NewAtTime(4, 0, 0))),
		gocron.NewTask(s.genGrade),
	)
	s.scheduler.NewJob(
		gocron.WeeklyJob(1, gocron.NewWeekdays(time.Saturday), gocron.NewAtTimes(gocron.NewAtTime(4, 30, 0))),
		gocron.NewTask(s.genGuilds),
	)

	s.scheduler.NewJob(
		gocron.WeeklyJob(1, gocron.NewWeekdays(time.Friday), gocron.NewAtTimes(gocron.NewAtTime(12, 0, 0))),
		gocron.NewTask(s.genPromotionsStat),
	)
	startHour := 0
	intervalHours := 2
	endHour := 24
	for hour := startHour; hour <= endHour; hour += intervalHours {
		s.scheduler.NewJob(
			gocron.WeeklyJob(1, gocron.NewWeekdays(time.Saturday), gocron.NewAtTimes(gocron.NewAtTime(uint(hour), 0, 0))),
			gocron.NewTask(s.genIncentiveStat),
		)
	}
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
	startDateOfNews := GetPreviousDate(15 * 7)

	log.Info("genGrade...", "start", start, "end", end)

	// translation guild grade
	if err := s.guild.GenGrade(start, end); err != nil {
		log.Error("genGrade failed", "err", err)
	}

	// developer guild grade
	if err := s.guild.GenDevGrade(last, end); err != nil {
		log.Error("genDevGrade failed", "err", err)
	}

	// news guild grade
	if err := s.guild.GenNewsGrade(startDateOfNews, end); err != nil {
		log.Error("genNewsGrade failed", "err", err)
	}

	log.Info("genGrade done")
}

func (s *Service) genPromotionsStat() {
	// Load latest contributors
	s.guild.LoadContributors()

	end := GetCurrentDate()
	// Automatic settlement of brand promotion points
	if err := s.guild.GenPromotionSettlement(end); err != nil {
		log.Error("Automatic settlement of brand promotion points failed", "err", err)
	}
	log.Info("genPromotionsStat done")
}

func (s *Service) genIncentiveStat() {
	now := GetCurrentDate()
	daysAgo := GetPreviousDate(6)
	//没有记录不执行
	exist, err := s.guild.IsExistRecord(now)
	if !exist || err != nil {
		return
	}
	exist = s.guild.IsExistIncentiveStatRecord(daysAgo)
	if exist {
		return
	}
	success, paymentDateMap, err := s.guild.GenIncentiveStat(now)
	if err != nil {
		log.Error("GenIncentiveStat failed", "err", err)
	}
	if success {
		err = s.guild.GenTotalIncentiveStat(paymentDateMap)
		if err != nil {
			log.Error("GenTotalIncentiveStat failed", "err", err)
		}
	}
}
func GetCurrentDate() (date string) {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		panic(err)
	}
	now := time.Now().In(loc)
	date = now.Format("2006-01-02")
	return
}

func GetPreviousDate(days int) (date string) {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		panic(err)
	}
	now := time.Now().In(loc)
	last := now.AddDate(0, 0, -days)
	date = last.Format("2006-01-02")
	return
}
