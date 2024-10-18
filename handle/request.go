package handle

import (
	"crypto/tls"
	"io"
	"net/http"
	"net/url"
	"spider/config"
	"strings"
	"time"

	"github.com/chentiangang/xlog"
)

type Request struct {
	req *http.Request
}

func NewRequest(cfg config.RequestConfig) (*Request, error) {
	URL, err := buildURL(cfg)

	// 根据请求方法构造请求
	var req *http.Request
	if cfg.Method == http.MethodPost {
		req, err = http.NewRequest(cfg.Method, URL, strings.NewReader(cfg.Body))
	} else {
		req, err = http.NewRequest(cfg.Method, URL, nil) // GET 请求不需要 Body
	}
	if err != nil {
		return nil, err
	}

	// 设置请求头
	for key, value := range cfg.Headers {
		req.Header.Set(key, value)
	}

	return &Request{req: req}, nil
}

func buildURL(cfg config.RequestConfig) (string, error) {
	parsedURL, err := url.Parse(cfg.URL)
	if err != nil {
		return "", err
	}

	// 如果有 URL 查询参数，将它们添加到 URL 查询字符串中
	if len(cfg.Params) > 0 {
		query := parsedURL.Query() // 获取现有的查询参数
		for key, value := range cfg.Params {
			query.Set(key, value) // 设置新的查询参数
		}
		parsedURL.RawQuery = query.Encode() // 将查询参数附加到 URL
	}
	return parsedURL.String(), nil
}

func (r *Request) SendRequest(cookie string) (bs []byte, err error) {
	client := &http.Client{
		Timeout: time.Second * 10,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	r.req.Header.Set("Cookie", cookie)
	resp, err := client.Do(r.req)
	if err != nil {
		xlog.Error("%s", err)
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		xlog.Error("received non-200 response code")
		return
	}

	bs, err = io.ReadAll(resp.Body)
	if err != nil {
		xlog.Error("failed to read response body,%s", err)
	}
	return bs, err
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
