package cookie

import (
	"fmt"
)

type BaiduFetcher struct {
	*Chromedp
}

func (f BaiduFetcher) String() string {
	ck := make(map[string]string)
	for _, i := range f.cookies {
		ck[i.Name] = i.Value
	}
	return fmt.Sprintf("%s,%s", ck["a"], ck["b"])
}
