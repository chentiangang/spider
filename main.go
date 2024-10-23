package main

import (
	"log"
	"spider/config"
	"spider/cookie"
	"spider/scheduler"
	"spider/tasks"
)

var cookieManager cookie.Manager

//func init() {
//cookieManager.Register("")
//}

func main() {
	// 加载配置
	cfg, err := config.LoadConfig("config.yml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 创建调度器
	sched := scheduler.NewScheduler()
	for _, i := range cfg.Tasks {
		if i.Cookie.Actions != nil {
			cookieFetcher := cookie.CreateFetcher(i.Cookie.Actions, i.Cookie.FetcherName, i.Cookie.LoginURL)
			cookieFetcher.Update()
			cookieManager.Register(i.Cookie.FetcherName, cookieFetcher)
			sched.AddTask(i.Cookie.Schedule, cookieFetcher.Update)
		}
	}

	// 初始化并添加任务
	for _, taskCfg := range cfg.Tasks {
		task := tasks.NewTask(taskCfg, cookieManager.GetCookie)
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
