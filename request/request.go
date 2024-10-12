package request

import (
	"errors"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type Request interface {
	SendRequest(cookie string) ([]byte, error)
}

type APIRequestManager struct {
	url     string
	method  string
	headers map[string]string
	body    string
}

func NewAPIRequestManager(url, method, body string, headers map[string]string) *APIRequestManager {
	return &APIRequestManager{
		url:     url,
		method:  method,
		headers: headers,
		body:    body,
	}
}

func (r *APIRequestManager) SendRequest(cookie string) ([]byte, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest(r.method, r.url, strings.NewReader(r.body))
	if err != nil {
		return nil, err
	}

	for key, value := range r.headers {
		req.Header.Add(key, value)
	}
	req.Header.Add("Cookie", cookie)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New("received non-200 response code")
	}

	return io.ReadAll(resp.Body)
}

// 装饰器：增加重试机制
func WithRetry(r Request, retries int, delay time.Duration) Request {
	return &retryRequestManager{
		Request: r,
		retries: retries,
		delay:   delay,
	}
}

type retryRequestManager struct {
	Request
	retries int
	delay   time.Duration
}

func (r *retryRequestManager) SendRequest(cookie string) ([]byte, error) {
	var err error
	for i := 0; i <= r.retries; i++ {
		var response []byte
		response, err = r.Request.SendRequest(cookie)
		if err == nil {
			return response, nil
		}
		log.Printf("Request failed, retrying in %v: %v", r.delay, err)
		time.Sleep(r.delay)
	}
	return nil, err
}
