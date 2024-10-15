package cookie

import (
	"context"
	"log"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

type Chromedp struct {
	url           string
	cookies       []*network.Cookie
	updateCycle   time.Duration // 更新周期
	lastUpdated   time.Time
	browserAction BrowserAction // 项目特定的操作链
}

// BrowserAction 定义一个浏览器操作的函数类型
type BrowserAction func() chromedp.Tasks

// NewChromedpCookieManager 允许传入不同的浏览器操作
func NewChromedp(url string, updateCycle time.Duration, action BrowserAction) *Chromedp {
	return &Chromedp{
		url:           url,
		updateCycle:   updateCycle,
		lastUpdated:   time.Time{},
		browserAction: action,
	}
}

// GetCookie 获取并返回 Cookie
func (cm *Chromedp) GetCookie() error {
	if time.Since(cm.lastUpdated) > cm.updateCycle {
		if err := cm.UpdateCookie(); err != nil {
			return err
		}
	}
	return nil
}

func (cm *Chromedp) UpdateCookie() error {
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
		cm.browserAction(),
		chromedp.ActionFunc(func(ctx context.Context) error {
			// 使用 chromedp.Cookies 获取所有 cookies
			cookies, err = network.GetCookies().Do(ctx)
			return err
		}))
	if err != nil {
		log.Println(err)
		return err
	}
	cm.cookies = cookies
	cm.lastUpdated = time.Now()
	return nil
}
