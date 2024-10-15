package cookie

type Fetcher interface {
	String() string
	Update()
}

type Manager struct {
	fetchers map[string]Fetcher
}

// 注册获取器到主获取器
func (m *Manager) Register(name string, fetcher Fetcher) {
	m.fetchers[name] = fetcher
}

func (m *Manager) Get(name string) string {
	return m.fetchers[name].String()
}

//// 获取指定名称的 Cookie
//func (m *Manager) GetCookies(name string) (string, error) {
//	if fetcher, ok := m.fetchers[name]; ok {
//		err := fetcher.Fetch()
//		if err != nil {
//			return "", err
//		}
//		return fetcher.String(), nil
//	}
//	return "", fmt.Errorf("fetcher not found: %s", name)
//}

//// 定时更新所有 Cookie
//func (m *Manager) UpdateAllCookies() {
//	for _, fetcher := range m.fetchers {
//		go fetcher.Update() // 使用 goroutine 并发更新
//	}
//}
