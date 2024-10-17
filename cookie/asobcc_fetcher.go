package cookie

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/chentiangang/xlog"
)

type AsoBccFetcher struct {
	*Chromedp
}

func (a AsoBccFetcher) String() string {
	cookies := make(map[string]string)
	for _, i := range a.cookies {
		cookies[i.Name] = i.Value
	}

	var aso, bcc string
	aso = fmt.Sprintf("aliyun_territory=%s; aliyun_lang=%s; apsara_api_token=%s; x-as-console-token-enc=%s; x-as-console-token=%s; aliyun_asapi_key=%s; login_aliyunid=%s; hb-1_login_aliyunid=%s; login_aliyunid_ticket=%s; hb-1_login_aliyunid_ticket=%s; login_aliyunid_csrf=%s; hb-1_login_aliyunid_csrf=%s; aliyun_country=%s; aso_region=%s; deptId=%s; top_level_domain_biz=%s; top_level_domain_ops=%s; deptName=",
		cookies["aliyun_territory"], cookies["aliyun_lang"], cookies["apsara_api_token"], cookies["x-as-console-token-enc"], cookies["x-as-console-token"], cookies["aliyun_asapi_key"],
		cookies["login_aliyunid"], cookies["hb-1_login_aliyunid"], cookies["login_aliyunid_ticket"], cookies["hb-1_login_aliyunid_ticket"], cookies["login_aliyunid_csrf"], cookies["hb-1_login_aliyunid_csrf"], cookies["aliyun_country"],
		cookies["aso_region"], cookies["deptId"], cookies["top_level_domain_biz"], cookies["top_level_domain_ops"])
	aso = aso + cookies["deptName"] + ";"
	bcc = bccSSOToken(aso)
	return string(aso + "; " + bcc)
}

func bccSSOToken(cookie string) string {
	r, err := http.NewRequest("GET", "https://abm.bjdc-1.ops.sgmc.sgcc.com.cn/gateway/v2/common/authProxy/auth/user/info?appId=bcc&noCache=1", nil)
	if err != nil {
		xlog.Error("%s", err)
		return ""
	}

	r.Header.Add("cookie", cookie)
	r.Header.Add("authority", "abm.bjdc-1.ops.sgmc.sgcc.com.cn")
	//r.Header.Add("path", "/gateway/v2/common/authProxy/auth/user/info?appId=bcc&noCache=1")
	r.Header.Add("scheme", "https")
	r.Header.Add("accept", "application/json, text/plain, */*")
	r.Header.Add("accept-encoding", "gzip, deflate, br")
	r.Header.Add("accept-language", "zh-CN,zh;q=0.9")
	r.Header.Add("sec-ch-ua", `".Not/A)Brand";v="99", "Google Chrome";v="103", "Chromium";v="103"`)
	r.Header.Add("sec-ch-ua-mobile", "?0")
	r.Header.Add("sec-ch-ua-platform", "Windows")
	r.Header.Add("sec-fetch-dest", "empty")
	r.Header.Add("sec-fetch-mode", "cors")
	r.Header.Add("sec-fetch-site", "same-origin")
	r.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.5060.66 Safari/537.36")
	r.Header.Add("x-biz-app", "base")
	r.Header.Add("x-env", "prod")

	client := http.Client{
		Timeout: time.Second * 10,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	resp, err := client.Do(r)
	if err != nil {
		xlog.Error("%s", err)
		return "err"

	}
	return resp.Header.Get("set-cookie")
}
