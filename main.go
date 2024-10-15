package main

import (
	"log"
	"spider/config"
	"spider/cookie"
	"spider/scheduler"
	"spider/tasks"
)

var cookieManager cookie.Manager

func init() {
	//cookieManager.Register("")
}

func main() {
	// 加载配置
	cfg, err := config.LoadConfig("config.yml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 创建调度器
	sched := scheduler.NewScheduler()
	for _, i := range cfg.Tasks {
		tasks := cookie.GenerateTasks(i.Cookie.Actions)
		cookieServer := cookie.NewChromedp(i.Cookie.URL, tasks)
		cookieServer.Update()
		cookieManager.Register(i.Cookie.Name, cookieServer)
		sched.AddTask(i.Cookie.Schedule, cookieServer.Update)
	}

	// 初始化并添加任务
	for _, taskCfg := range cfg.Tasks {
		var task tasks.Task[float64]
		if err := task.Init(taskCfg, cookieManager.Get); err != nil {
			log.Printf("Failed to init task %s: %v", taskCfg.Name, err)
			continue
		}
		if err := sched.AddTask(taskCfg.Schedule, task.Execute); err != nil {
			log.Printf("Failed to add task %s to scheduler: %v", taskCfg.Name, err)
			continue
		}
	}

	// 启动调度器
	sched.Start()

	// 阻塞主线程
	select {}
}
