package guild

import (
	"fmt"
	"github.com/dstotijn/go-notion"
	"github.com/permadao/dnotion/config"
	"github.com/permadao/dnotion/db"
	"github.com/permadao/dnotion/utils"
	"sort"
	"sync"
	"testing"
	"time"
)

func TestGrade(t *testing.T) {
	// c := config.New("config_temp")
	// d := db.New(c)
	// g := New(c, d)
	// start := time.Now().AddDate(0, 0, -28).Format("2006-01-02")
	// end := time.Now().Format("2006-01-02")
	// err := g.GenGrade("e8d79c55c0394cba83664f3e5737b0bd", "d8c270f68a8f44aaa6b24e17c927df2b", start, end)
	// fmt.Println(err)
}

func TestDevGrade(t *testing.T) {
	c := config.New("config_temp")
	d := db.New(c)
	g := New(c, d)
	now := time.Now()
	startDate := now.AddDate(0, 0, -15*7)
	startDateOfNews := startDate.Format("2006-01-02")
	end := time.Now().Format("2006-01-02")
	err := g.GenNewsGrade("ad2cf585b08843fea7cf40a682bc4529", "d5f9fc70910b45d4ab8811f37716637d", startDateOfNews, end)
	fmt.Println(err)
}

func TestGuild_GenPromotionSettlement(t *testing.T) {
	c := config.New("config_temp")
	d := db.New(c)
	g := New(c, d)
	err := g.GenPromotionSettlement("14debb08a4e8416e9b0de7ce46821506", "2ea3ff42b3b84d5cbc9a575d4c436878", "2024-05-24")
	if err != nil {
		fmt.Println(err)
	}
}

func TestGuild_GenIncentiveStat(t *testing.T) {
	c := config.New("config_temp")
	d := db.New(c)
	g := New(c, d)
	startDate, _ := time.Parse("2006-01-02", "2024-07-06")
	startDateStr := startDate.Format("2006-01-02")
	success, paydateMap, _ := g.GenIncentiveStat(utils.CincentiveWeeklyGuildRs, startDateStr)
	fmt.Println(success)
	fmt.Println(paydateMap)
}

func TestGuild_GenTotalIncentiveStat(t *testing.T) {
	c := config.New("config_temp")
	d := db.New(c)
	g := New(c, d)
	startDate, acDate := "2024-06-28", "2024-06-28"
	start, _ := notion.ParseDateTime(startDate)
	end, _ := notion.ParseDateTime(acDate)
	records, _ := g.db.GetIncentiveData(&notion.DatabaseQueryFilter{
		And: []notion.DatabaseQueryFilter{
			{
				Property: "Payment Date",
				DatabaseQueryPropertyFilter: notion.DatabaseQueryPropertyFilter{
					Date: &notion.DatePropertyFilter{
						OnOrAfter: &start.Time,
					},
				},
			},
			{
				Property: "Payment Date",
				DatabaseQueryPropertyFilter: notion.DatabaseQueryPropertyFilter{
					Date: &notion.DatePropertyFilter{
						OnOrBefore: &end.Time,
					},
				},
			},
		},
	})
	payDateMap := map[string]int{}
	for _, d := range records {
		payDateMap[d.PaymentDate]++
	}
	payDateSlice := []string{}
	for k, _ := range payDateMap {
		payDateSlice = append(payDateSlice, k)
	}
	sort.Slice(payDateSlice, func(i, j int) bool {
		return payDateSlice[i] < payDateSlice[j]
	})

	handle := &Handle{
		WorkerPoolSize: 1,
		TaskQueue:      make([]chan string, 5),
		Func: func(pd string, wg *sync.WaitGroup) bool {
			defer wg.Done()
			defer func() {
				if err := recover(); err != nil {
					fmt.Println(err)
				}
			}()
			fmt.Println("准备开始", pd)
			pm := map[string]int{}
			pm[pd]++
			err := g.GenTotalIncentiveStat(utils.CincentiveWeeklyRs, pm)
			if err != nil {
				fmt.Println("异常->", err)
				return false
			}
			fmt.Println("完成", pd)
			return true
		},
		mux: sync.Mutex{},
		wg:  sync.WaitGroup{},
	}
	handle.StartWorkerPool()
	for i, pd := range payDateSlice {
		handle.wg.Add(1)
		handle.SendMsgToTaskQueue(i%int(handle.WorkerPoolSize), pd)
	}
	fmt.Println("等待完成----")
	handle.wg.Wait()
	result := handle.TaskResult
	sort.Strings(result)
	fmt.Println("失败的天->", result)
}

// 修复没有Token的数据
func TestUpdateTotalIncentiveStatData(t *testing.T) {
	c := config.New("config_temp")
	d := db.New(c)
	g := New(c, d)
	startDate, acDate := "2022-12-31", "2024-06-20"
	start, _ := notion.ParseDateTime(startDate)
	end, _ := notion.ParseDateTime(acDate)
	records, _ := g.db.GetTotalIncentiveData(&notion.DatabaseQueryFilter{
		And: []notion.DatabaseQueryFilter{
			{
				Property: "Payment Date",
				DatabaseQueryPropertyFilter: notion.DatabaseQueryPropertyFilter{
					Date: &notion.DatePropertyFilter{
						OnOrAfter: &start.Time,
					},
				},
			},
			{
				Property: "Payment Date",
				DatabaseQueryPropertyFilter: notion.DatabaseQueryPropertyFilter{
					Date: &notion.DatePropertyFilter{
						OnOrBefore: &end.Time,
					},
				},
			},
			{
				Property: "Token",
				DatabaseQueryPropertyFilter: notion.DatabaseQueryPropertyFilter{
					Select: &notion.SelectDatabaseQueryFilter{
						IsEmpty: true,
					},
				},
			},
		},
	})
	payDateMap := map[string]int{}
	for _, d := range records {
		payDateMap[d.PaymentDate]++
	}
	payDateSlice := []string{}
	for k, _ := range payDateMap {
		payDateSlice = append(payDateSlice, k)
	}
	sort.Slice(payDateSlice, func(i, j int) bool {
		return payDateSlice[i] < payDateSlice[j]
	})

	handle := &Handle{
		WorkerPoolSize: 5,
		TaskQueue:      make([]chan string, 5),
		Func: func(pd string, wg *sync.WaitGroup) bool {
			defer wg.Done()
			defer func() {
				if err := recover(); err != nil {
					fmt.Println(err)
				}
			}()
			fmt.Println("准备开始", pd)
			paymentDate, _ := notion.ParseDateTime(pd)
			d, _ := g.db.GetTotalIncentiveData(&notion.DatabaseQueryFilter{
				And: []notion.DatabaseQueryFilter{
					{
						Property: "Payment Date",
						DatabaseQueryPropertyFilter: notion.DatabaseQueryPropertyFilter{
							Date: &notion.DatePropertyFilter{
								Equals: &paymentDate.Time,
							},
						},
					},
				},
			})
			for _, ti := range d {
				if ti.Token == "" {
					ti.Token = "USD"
					err := g.db.UpdatePage(&ti)
					if err != nil {
						fmt.Println("异常->", err)
						return false
					}
				}
			}
			fmt.Println("完成", pd)
			return true
		},
		mux: sync.Mutex{},
		wg:  sync.WaitGroup{},
	}
	handle.StartWorkerPool()
	for i, pd := range payDateSlice {
		handle.wg.Add(1)
		handle.SendMsgToTaskQueue(i%int(handle.WorkerPoolSize), pd)
	}
	fmt.Println("等待完成----")
	handle.wg.Wait()
	result := handle.TaskResult
	sort.Strings(result)
	fmt.Println("失败的天->", result)
}

// 修复没有Token的数据
func TestUpdateIncentiveStatData(t *testing.T) {
	c := config.New("config_temp")
	d := db.New(c)
	g := New(c, d)
	startDate, acDate := "2022-12-31", "2024-06-20"
	start, _ := notion.ParseDateTime(startDate)
	end, _ := notion.ParseDateTime(acDate)
	records, _ := g.db.GetIncentiveData(&notion.DatabaseQueryFilter{
		And: []notion.DatabaseQueryFilter{
			{
				Property: "Payment Date",
				DatabaseQueryPropertyFilter: notion.DatabaseQueryPropertyFilter{
					Date: &notion.DatePropertyFilter{
						OnOrAfter: &start.Time,
					},
				},
			},
			{
				Property: "Payment Date",
				DatabaseQueryPropertyFilter: notion.DatabaseQueryPropertyFilter{
					Date: &notion.DatePropertyFilter{
						OnOrBefore: &end.Time,
					},
				},
			},
			{
				Property: "Token",
				DatabaseQueryPropertyFilter: notion.DatabaseQueryPropertyFilter{
					Select: &notion.SelectDatabaseQueryFilter{
						IsEmpty: true,
					},
				},
			},
		},
	})
	payDateMap := map[string]int{}
	for _, d := range records {
		payDateMap[d.PaymentDate]++
	}
	payDateSlice := []string{}
	for k, _ := range payDateMap {
		payDateSlice = append(payDateSlice, k)
	}
	sort.Slice(payDateSlice, func(i, j int) bool {
		return payDateSlice[i] < payDateSlice[j]
	})

	handle := &Handle{
		WorkerPoolSize: 5,
		TaskQueue:      make([]chan string, 5),
		Func: func(pd string, wg *sync.WaitGroup) bool {
			defer wg.Done()
			defer func() {
				if err := recover(); err != nil {
					fmt.Println(err)
				}
			}()
			fmt.Println("准备开始", pd)
			paymentDate, _ := notion.ParseDateTime(pd)
			incentiveData, _ := g.db.GetIncentiveData(&notion.DatabaseQueryFilter{
				And: []notion.DatabaseQueryFilter{
					{
						Property: "Payment Date",
						DatabaseQueryPropertyFilter: notion.DatabaseQueryPropertyFilter{
							Date: &notion.DatePropertyFilter{
								Equals: &paymentDate.Time,
							},
						},
					},
				},
			})
			for _, incentive := range incentiveData {
				if incentive.Token == "" {
					incentive.Token = "USD"
					err := g.db.UpdatePage(&incentive)
					if err != nil {
						fmt.Println("异常->", err)
						return false
					}
				}
			}
			fmt.Println("完成", pd)
			return true
		},
		mux: sync.Mutex{},
		wg:  sync.WaitGroup{},
	}
	handle.StartWorkerPool()
	for i, pd := range payDateSlice {
		handle.wg.Add(1)
		handle.SendMsgToTaskQueue(i%int(handle.WorkerPoolSize), pd)
	}
	fmt.Println("等待完成----")
	handle.wg.Wait()
	result := handle.TaskResult
	sort.Strings(result)
	fmt.Println("失败的天->", result)
}

type Handle struct {
	WorkerPoolSize uint32
	TaskQueue      []chan string
	Func           func(string, *sync.WaitGroup) bool
	TaskResult     []string
	mux            sync.Mutex
	wg             sync.WaitGroup
}

func (h *Handle) StartWorkerPool() {
	for i := 0; i < int(h.WorkerPoolSize); i++ {
		h.TaskQueue[i] = make(chan string, 5)
		go h.StartOneWorker(i, h.TaskQueue[i])
	}
}

func (h *Handle) StartOneWorker(workerID int, taskQueue chan string) {
	for {
		select {
		case request := <-taskQueue:
			fmt.Println(workerID, "正在执行", request)
			result := h.Func(request, &h.wg)
			if !result {
				h.mux.Lock()
				h.TaskResult = append(h.TaskResult, request)
				h.mux.Unlock()
			}
		}
	}
}

func (h *Handle) SendMsgToTaskQueue(workerID int, request string) {
	h.TaskQueue[workerID] <- request
}
