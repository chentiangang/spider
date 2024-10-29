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
			cookieFetcher := cookie.CreateFetcher(i.Cookie.Actions, i.Cookie.FetcherName, i.Cookie.LoginURL)
			cookieFetcher.Update()
			cookieManager.Register(i.Cookie.FetcherName, cookieFetcher)
			sched.AddTask(i.Cookie.Schedule, cookieFetcher.Update)

			if i.Cookie.HttpServerPath != "" {
				// 创建一个局部变量用于捕获当前的路径和 fetcher 名称
				path := i.Cookie.HttpServerPath
				fetcherName := i.Cookie.FetcherName
				http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
					log.Println(r.Host, r.RequestURI)
					res := cookieManager.GetCookie(fetcherName)
					w.Write([]byte(res))
				})
			}
		}
	}

	// 初始化并添加任务
	for _, taskCfg := range cfg.Tasks {
		if taskCfg.HandlerName == "" {
			continue
		}
		if taskCfg.Request == nil {
			continue
		}

		if cookieManager.GetCookie(taskCfg.Cookie.FetcherName) != "" {
			task := tasks.NewTask(taskCfg, cookieManager.GetCookie)
			if err := sched.AddTask(taskCfg.Schedule, task.Execute); err != nil {
				log.Printf("Failed to add task %s to scheduler: %v", taskCfg.Name, err)
				continue
			}
		}
	}

	// 启动调度器
	sched.Start()

	http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", cfg.HttpPort), nil)

	// 阻塞主线程
	//select {}
}
