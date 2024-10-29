package main

import (
	"fmt"
	"log"
	"net/http"
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

	// 初始化并添加cookie fetcher
	for _, i := range cfg.Tasks {
		if i.Cookie.Actions != nil {
			fetcher := cookie.CreateFetcher(i.Cookie.Actions, i.Cookie.FetcherName, i.Cookie.LoginURL)
			fetcher.Update()
			cookieManager.Register(i.Cookie.FetcherName, fetcher)
			sched.AddTask(i.Cookie.Schedule, fetcher.Update)

			if i.Cookie.HttpServerPath != "" {
				handler := CreateCookieHttpHandlerFunc(i.Cookie.FetcherName, &cookieManager)
				http.HandleFunc(i.Cookie.HttpServerPath, handler)
			}
		}
	}

	addTaskScheduler(cfg, sched, &cookieManager)

	// 启动调度器
	sched.Start()
	http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", cfg.HttpPort), nil)
}

func CreateCookieHttpHandlerFunc(fetcherName string, cm *cookie.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Host, r.RequestURI)
		res := cm.GetCookie(fetcherName)
		w.Write([]byte(res))
	}
}

func addTaskScheduler(cfg *config.Config, sched *scheduler.Scheduler, cm *cookie.Manager) {
	for _, taskCfg := range cfg.Tasks {
		if taskCfg.HandlerName == "" || taskCfg.Request == nil {
			continue
		}

		if cm.GetCookie(taskCfg.Cookie.FetcherName) != "" {
			task := tasks.NewTask(taskCfg, cm.GetCookie)
			if err := sched.AddTask(taskCfg.Schedule, task.Execute); err != nil {
				log.Printf("Failed to add task %s to scheduler: %v", taskCfg.Name, err)
			}
		}
	}
}
