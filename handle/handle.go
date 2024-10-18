package handle

import "spider/config"

type Handler interface {
	Init(config config.TaskConfig) error
	Name() string
	SendRequest(cookie string) <-chan []byte // 发送请求并返回数据通道
	ParseToChan(data <-chan []byte)          // 解析响应数据并返回解析结果通道
	Store()                                  // 存储解析后的数据
}

var Handlers []Handler

func init() {
	Handlers = append(Handlers, &ProjectSummaryHandler{})
}

func CreateTaskHandler(cfg config.TaskConfig) Handler {
	for _, h := range Handlers {
		if h.Name() == cfg.HandlerName {
			return h
		}
	}
	return nil
}
