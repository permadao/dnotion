package service

import (
	"time"

	"github.com/go-co-op/gocron/v2"
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
	if err := s.guild.GenGrade("e8d79c55c0394cba83664f3e5737b0bd", "d8c270f68a8f44aaa6b24e17c927df2b", start, end); err != nil {
		log.Error("genGrade failed", "err", err)
	}

	// developer guild grade
	if err := s.guild.GenDevGrade("146e1f661ed943e3a460b8cf12334b7b", "623ccfc9fb1443279decf90fb752215d", last, end); err != nil {
		log.Error("genDevGrade failed", "err", err)
	}

	// news guild grade
	if err := s.guild.GenNewsGrade("ad2cf585b08843fea7cf40a682bc4529", "d5f9fc70910b45d4ab8811f37716637d", startDateOfNews, end); err != nil {
		log.Error("genNewsGrade failed", "err", err)
	}

	log.Info("genGrade done")
}

func (s *Service) genPromotionsStat() {
	end := GetCurrentDate()
	// Automatic settlement of brand promotion points
	if err := s.guild.GenPromotionSettlement("14debb08a4e8416e9b0de7ce46821506", "2ea3ff42b3b84d5cbc9a575d4c436878", end); err != nil {
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
	success, paymentDateMap, err := s.guild.GenIncentiveStat("4c19704d927f4d52b2f030ebd1648ef3", now)
	if err != nil {
		log.Error("GenIncentiveStat failed", "err", err)
	}
	if success {
		err = s.guild.GenTotalIncentiveStat("04c301f8dc5448759c5919e618822854", paymentDateMap)
		if err != nil {
			log.Error("GenTotalIncentiveStat failed", "err", err)
		}
	}
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
