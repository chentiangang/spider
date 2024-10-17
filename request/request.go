package request

import (
	"io"
	"net/http"
	"net/url"
	"spider/config"
	"strings"

	"github.com/chentiangang/xlog"
)

type Requester interface {
	SendRequest(cookie string) (<-chan []byte, error)
}

type APIRequest struct {
	req *http.Request
}

func NewAPIRequest(cfg config.RequestConfig) *APIRequest {
	req, err := BuildRequest(cfg)
	if err != nil {
		xlog.Error("%s", err)
		return nil
	}
	return &APIRequest{
		req: req,
	}
}

// BuildRequest 构造 HTTP 请求
func BuildRequest(cfg config.RequestConfig) (*http.Request, error) {
	// 解析 URL
	parsedURL, err := url.Parse(cfg.URL)
	if err != nil {
		return nil, err
	}

	// 如果有 URL 查询参数，将它们添加到 URL 查询字符串中
	if len(cfg.Params) > 0 {
		query := parsedURL.Query() // 获取现有的查询参数
		for key, value := range cfg.Params {
			query.Set(key, value) // 设置新的查询参数
		}
		parsedURL.RawQuery = query.Encode() // 将查询参数附加到 URL
	}

	// 根据请求方法构造请求
	var req *http.Request
	if cfg.Method == http.MethodPost {
		req, err = http.NewRequest(cfg.Method, parsedURL.String(), strings.NewReader(cfg.Body))
	} else {
		req, err = http.NewRequest(cfg.Method, parsedURL.String(), nil) // GET 请求不需要 Body
	}
	if err != nil {
		return nil, err
	}

	// 设置请求头
	for key, value := range cfg.Headers {
		req.Header.Set(key, value)
	}

	//req.Header.Set("Cookie", )

	return req, nil
}

func (r *APIRequest) SendRequest(cookie string) <-chan []byte {
	client := newClient()

	res := make(chan []byte)

	go func() {
		defer close(res)
		r.req.Header.Set("Cookie", cookie)
		resp, err := client.Do(r.req)
		if err != nil {
			xlog.Error("%s", err)
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			xlog.Error("received non-200 response code")
			return
		}

		bs, err := io.ReadAll(resp.Body)
		if err != nil {
			xlog.Error("failed to read response body,%s", err)
		}
		res <- bs

	}()
	return res
}

//// 装饰器：增加重试机制
//func WithRetry(r Requester, retries int, delay time.Duration) Requester {
//	return &retryRequestManager{
//		Requester: r,
//		retries:   retries,
//		delay:     delay,
//	}
//}
//
//type retryRequestManager struct {
//	Requester
//	retries int
//	delay   time.Duration
//}
//
//func (r *retryRequestManager) SendRequest(cookie string) ([]byte, error) {
//	var err error
//	for i := 0; i <= r.retries; i++ {
//		var response []byte
//		response, err = r.Requester.SendRequest(cookie)
//		if err == nil {
//			return response, nil
//		}
//		log.Printf("Request failed, retrying in %v: %v", r.delay, err)
//		time.Sleep(r.delay)
//	}
//	return nil, err
//}
