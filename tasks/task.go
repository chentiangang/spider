package tasks

import (
	"log"
	"spider/config"
	"spider/parser"
	"spider/request"

	"github.com/chentiangang/xlog"
)

type Task struct {
	config  config.TaskConfig
	request request.Requester
	cookie  func(name string) string
	parser  parser.Parser
	//storage storage.Storage[
}

// Init 包含了这个任务实例的初始化操作
func (t *Task) Init(config config.TaskConfig, cookieFunc func(name string) string) error {
	t.config = config
	t.cookie = cookieFunc
	t.request = request.NewAPIRequest(config.Request)
	//t.request.SendRequest()
	//cookie.NewChromedp(t.config.Cookie.Method)
	//t.parser =
	//t.processor = &processor.DBProcessor{} // 根据配置选择处理器
	return nil
}

// Execute 是该任务的执行函数
func (t *Task) Execute() {
	// 模拟浏览器获取 cookie
	// 这里根据 t.config.CookieConfig 进行具体实现
	// ...
	xlog.Debug("%+v", t)
	return

	// 发送 HTTP 请求获取数据
	cookie := t.cookie(t.config.Cookie.Name)
	if cookie == "" {
		xlog.Error("cookie is empty")
		return
	}

	bs, err := t.request.SendRequest(cookie)
	if err != nil {
		return
	}
	xlog.Debug("%s", bs)

	// 处理数据
	//res, err := t.parser.Parse(bs)
	//if err != nil {
	//	return
	//}
	//
	//err = t.storage.Save(res)
	//if err != nil {
	//	return
	//}
	//
	log.Printf("Task %s executed successfully", t.config.Name)
}

// NewTask 是新建一个任务实例
//func NewTask(config config.TaskConfig) Task[T] {
//	return Task{}
//}
