package tasks

import (
	"log"
	"spider/config"
	"spider/handle"
	"sync"

	"github.com/chentiangang/xlog"
)

type Task struct {
	config            config.TaskConfig
	cookieFetcherFunc func(name string) string
	handler           handle.Handler
}

// Execute 是该任务的执行函数
func (t *Task) Execute() {
	// 发送 HTTP 请求获取数据
	cookie := t.cookieFetcherFunc(t.config.Cookie.FetcherName)
	if cookie == "" {
		xlog.Error("cookie is empty")
		return
	}

	err := t.handler.Init(t.config)
	if err != nil {
		xlog.Error("%s", err)
		return
	}

	bsCh := t.handler.SendRequest(cookie)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		t.handler.ParseToChan(bsCh)
	}()

	go func() {
		defer wg.Done()
		t.handler.Store()
	}()

	wg.Wait()
	log.Printf("Task %s executed successfully", t.config.Name)
}

// NewTask 是初始化一个Task 实例
func NewTask(config config.TaskConfig, cookieFunc func(name string) string) *Task {
	return &Task{
		config:            config,
		cookieFetcherFunc: cookieFunc,
		handler:           handle.CreateTaskHandler(config),
	}
}
