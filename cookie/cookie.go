package cookie

import (
	"fmt"
	"net/http"
)

type Fetcher interface {
	Get() ([]*http.Cookie, error)
	String() string
	Update() error
}

type Manager struct {
	fetchers map[string]Fetcher
}

// 注册获取器到主获取器
func (m *Manager) Register(name string, fetcher Fetcher) {
	m.fetchers[name] = fetcher
}

// 获取指定名称的 Cookie
func (m *Manager) GetCookies(name string) ([]*http.Cookie, error) {
	if fetcher, ok := m.fetchers[name]; ok {
		return fetcher.Get()
	}
	return nil, fmt.Errorf("fetcher not found: %s", name)
}

// 定时更新所有 Cookie
func (m *Manager) UpdateAllCookies() {
	for _, fetcher := range m.fetchers {
		go fetcher.Update() // 使用 goroutine 并发更新
	}
}
