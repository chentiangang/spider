package cookie

import (
	"context"
	"log"
	"time"

	"github.com/chromedp/chromedp"
)

type Chromedp struct {
	url           string
	cookies       string
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
func (cm *Chromedp) GetCookie() (string, error) {
	if time.Since(cm.lastUpdated) > cm.updateCycle {
		if err := cm.UpdateCookie(); err != nil {
			return "", err
		}
	}
	return cm.cookies, nil
}

// UpdateCookie 根据传入的操作链更新 Cookie
func (cm *Chromedp) UpdateCookie() error {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var cookies string
	// 使用传入的 BrowserAction 执行特定的浏览器操作
	err := chromedp.Run(ctx, chromedp.Navigate(cm.url), cm.browserAction(), chromedp.Evaluate(`document.cookie`, &cookies))
	if err != nil {
		return err
	}

	cm.cookies = cookies
	cm.lastUpdated = time.Now()
	log.Printf("Cookie updated: %s", cookies)
	return nil
}
