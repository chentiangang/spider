package cookie

import (
	"encoding/json"
	"fmt"

	"github.com/chentiangang/xlog"
)

type SjztFetcher struct {
	*Chromedp
}

func (a *SjztFetcher) Init(actions []Action, url string) {
	tasks := GenerateTasks(actions)
	a.Chromedp = NewChromedp(url, tasks)
}

func (a *SjztFetcher) String() string {
	cookie := make(map[string]string)
	for _, i := range a.cookies {
		cookie[i.Name] = i.Value
	}
	var res = make(map[string]string)
	res["cookie"] = fmt.Sprintf("Admin-Token=%s; Admin-Expires-In=%s; loginType=%s; sn_token=%s", cookie["Admin-Token"], cookie["Admin-Expires-In"], cookie["loginType"], cookie["sn_token"])
	res["token"] = fmt.Sprintf("Bearer %s", cookie["Admin-Token"])

	bs, err := json.Marshal(res)
	if err != nil {
		xlog.Error("%s", err)
		return ""
	}
	return string(bs)
}

func (a *SjztFetcher) Name() string {
	return "sjzt_fetcher"
}
