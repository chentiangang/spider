package tasks

import (
	"log"
	"spider/config"
	"spider/handle"

	"github.com/chentiangang/xlog"
)

type Task struct {
	config     config.TaskConfig
	cookieFunc func(name string) string
	handler    handle.Handler
}

// Execute 是该任务的执行函数
func (t *Task) Execute() {
	// 发送 HTTP 请求获取数据
	cookie := t.cookieFunc(t.config.Cookie.FetcherName)
	if cookie == "" {
		xlog.Error("cookie is empty")
		return
	}

	bs, err := t.handler.SendRequest(cookie)
	if err != nil {
		return
	}

	t.handler.ParseToChan(bs)
	t.handler.Store()
	log.Printf("Task %s executed successfully", t.config.Name)
}

// NewTask 是初始化一个Task 实例
func NewTask(config config.TaskConfig, cookieFunc func(name string) string) *Task {
	return &Task{
		config:     config,
		cookieFunc: cookieFunc,
		handler:    handle.CreateTaskHandler(config),
	}
}
