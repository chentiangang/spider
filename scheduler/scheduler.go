package scheduler

import (
	"log"

	cron "github.com/robfig/cron/v3"
)

type Scheduler struct {
	cron *cron.Cron
}

// NewScheduler 初始化调度器
func NewScheduler() *Scheduler {
	return &Scheduler{
		cron: cron.New(),
	}
}

// AddTask 添加定时任务
func (s *Scheduler) AddTask(schedule string, task func()) error {
	_, err := s.cron.AddFunc(schedule, task)
	if err != nil {
		return err
	}
	return nil
}

//func (s *Scheduler) AddCookieUpdate(schedule string, task func()) error {
//	_, err := s.cron.AddFunc(schedule, task)
//	if err != nil {
//		return err
//	}
//	return nil
//}

// Start 开始调度
func (s *Scheduler) Start() {
	s.cron.Start()
	log.Println("Scheduler started")
}

// Stop 停止调度
func (s *Scheduler) Stop() {
	s.cron.Stop()
}
