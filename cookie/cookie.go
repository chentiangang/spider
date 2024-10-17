package cookie

import (
	"log"
)

type Fetcher interface {
	String() string
	Update()
}

func CreateFetcher(actions []Action, name string, url string) Fetcher {
	tasks := GenerateTasks(actions)
	switch name {
	case "asobcc":
		fetcher := NewChromedp(url, tasks)
		return AsoBccFetcher{fetcher}
	default:
		log.Fatalf("Unknown project name: %s", name)
		return nil
	}
}

type Manager struct {
	fetchers map[string]Fetcher
}

// 注册获取器到主获取器
func (m *Manager) Register(name string, fetcher Fetcher) {
	if m.fetchers == nil {
		m.fetchers = make(map[string]Fetcher)
	}
	m.fetchers[name] = fetcher
}

func (m *Manager) GetCookie(name string) string {
	return m.fetchers[name].String()
}
