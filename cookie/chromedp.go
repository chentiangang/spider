package cookie

import (
	"context"
	"log"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

// ActionType 描述可以执行的操作类型
type ActionType string

const (
	ActionClick    ActionType = "click"    // 点击操作
	ActionSendKeys ActionType = "sendKeys" // 输入文本操作
	ActionSubmit   ActionType = "submit"   // 提交表单操作
	// 你可以根据需要扩展更多类型
)

// Action 定义了一个浏览器动作
type Action struct {
	Type     ActionType `yaml:"type"`     // 动作类型 (click, sendKeys, submit)
	Selector string     `yaml:"selector"` // 元素选择器
	Value    string     `yaml:"value"`    // 对于sendKeys，可以有输入的值
	By       string     `yaml:"by"`       // 选择器的方式 (ByID, ByQuery, ByXPath 等)
}

// GenerateTasks 根据配置文件生成浏览器操作任务
func GenerateTasks(actions []Action) chromedp.Tasks {
	var tasks chromedp.Tasks

	for _, action := range actions {
		switch action.Type {
		case ActionSendKeys:
			// 根据配置中的 "by" 动态选择选择器类型
			tasks = append(tasks, chromedp.SendKeys(action.Selector, action.Value, getBy(action.By)))
		case ActionClick:
			tasks = append(tasks, chromedp.Click(action.Selector, getBy(action.By)))
		case ActionSubmit:
			tasks = append(tasks, chromedp.Submit(action.Selector, getBy(action.By)))
		default:
			log.Printf("Unknown action type: %s", action.Type)
		}
	}
	return tasks
}

// getBy 根据配置的 "by" 字段返回 chromedp 的选择器类型
func getBy(by string) chromedp.QueryOption {
	switch by {
	case "ByID":
		return chromedp.ByID
	case "ByQuery":
		return chromedp.ByQuery
	case "ByXPath":
		return chromedp.BySearch
	default:
		return chromedp.BySearch // 默认为 BySearch 或你可以使用其他默认值
	}
}

type Chromedp struct {
	url         string
	cookies     []*network.Cookie
	updateCycle time.Duration // 更新周期
	lastUpdated time.Time
	actions     chromedp.Tasks // 项目特定的操作链
}

// NewChromedp 允许传入不同的浏览器操作
func NewChromedp(url string, action chromedp.Tasks) *Chromedp {
	return &Chromedp{
		url:     url,
		actions: action,
	}
}

func (cm *Chromedp) Update() {
	opts := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", false),
		chromedp.Flag("test-type", true),
		chromedp.Flag("ignore-certificate-errors", true),
		chromedp.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36 Edg/128.0.0.0"),
	}
	opts = append(chromedp.DefaultExecAllocatorOptions[:], opts...)

	ctx, _ := chromedp.NewExecAllocator(context.Background(), opts...)
	newctx, cancel := chromedp.NewContext(ctx)
	defer chromedp.Cancel(newctx)

	chromedp.Run(newctx, make([]chromedp.Action, 0, 1)...)

	timeoutCtx, cancel := context.WithTimeout(newctx, 180*time.Second)
	defer cancel()

	var cookies []*network.Cookie
	var err error
	err = chromedp.Run(timeoutCtx, chromedp.Navigate(cm.url),
		chromedp.WaitVisible("#kw", chromedp.ByID),
		cm.actions,
		chromedp.ActionFunc(func(ctx context.Context) error {
			// 使用 chromedp.Cookies 获取所有 cookies
			cookies, err = network.GetCookies().Do(ctx)
			return err
		}))
	if err != nil {
		log.Println(err)
		return
	}
	cm.cookies = cookies
	return
}

func (cm *Chromedp) String() string {
	return ""
}
