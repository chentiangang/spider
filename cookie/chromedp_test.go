package cookie

import (
	"testing"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/stretchr/testify/assert"
)

func MockBrowserAction() chromedp.Tasks {
	// 模拟设置 cookie 的操作
	return chromedp.Tasks{
		chromedp.SendKeys(`#kw`, "为什么", chromedp.ByID),
		chromedp.Click(),
		chromedp.Submit(`#su`, chromedp.ByID),
		//chromedp.Sleep(1 * time.Second),
	}
}

func TestGetCookie(t *testing.T) {
	url := "http://baidu.com" // 测试用的 URL
	updateCycle := 4 * time.Second
	cm := NewChromedp(url, updateCycle, MockBrowserAction)

	// 设置初始状态，模拟第一次调用 GetCookie 时 cookies 为空
	cm.lastUpdated = time.Now() // 将上次更新时间设置为当前，确保不会触发更新
	cm.cookies = nil            // 确保初始 cookie 是空的

	// 获取初始 cookie，应为空
	err := cm.GetCookie()
	assert.NoError(t, err)    // 确保没有错误
	assert.Nil(t, cm.cookies) // 初始状态应为空

	// 打印初始 cookie 的值（应为空）
	t.Logf("Initial cookie: %v", cm.cookies)

	// 手动设置 cookie 并更新
	cm.lastUpdated = time.Now().Add(-3 * time.Second) // 设置 lastUpdated 使其过期
	err = cm.UpdateCookie()
	assert.NoError(t, err) // 确保没有错误

	// 打印更新后的 cookie 值
	t.Logf("Updated cookie after first update: %v", cm.cookies)

	// 再次调用 GetCookie，应该获取到设置的值
	err = cm.GetCookie()
	assert.NoError(t, err)
	assert.Equal(t, true, len(cm.cookies) > 0)

	// 打印最终获取到的 cookie 值
	for _, cookie := range cm.cookies {
		t.Logf("Final cookie: %v", cookie)
	}
}

func TestUpdateCookie(t *testing.T) {
	url := "http://baidu.com"
	updateCycle := 2 * time.Second
	cm := NewChromedp(url, updateCycle, MockBrowserAction)

	// 更新 cookie
	err := cm.UpdateCookie()
	assert.NoError(t, err)

	// 打印更新后的 cookie 值
	//t.Logf("Updated cookie: %v", cm.cookies)
	// 打印最终获取到的 cookie 值
	for _, cookie := range cm.cookies {
		t.Logf("Final cookie: %v", cookie)
	}

	assert.Equal(t, true, len(cm.cookies) > 0)                 // 确保 cookies 被正确设置
	assert.True(t, time.Since(cm.lastUpdated) < 1*time.Second) // lastUpdated 时间应被更新
}
