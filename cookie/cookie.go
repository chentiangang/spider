package cookie

type Fetcher interface {
	Init(actions []Action, url string)
	String() string
	Update()
	Name() string
}

var Fetchers []Fetcher

func init() {
	Fetchers = append(Fetchers, &AsoBccFetcher{})
}

func CreateFetcher(actions []Action, name string, url string) Fetcher {
	for _, fetcher := range Fetchers {
		if name == fetcher.Name() {
			fetcher.Init(actions, url)
			return fetcher
		}
	}
	return nil
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
