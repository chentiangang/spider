package cookie

type Cookie interface {
	GetCookie() (string, error)
	UpdateCookie() error
}

func NewCookieWithMethod(method string) {
	switch method {
	case "chromed":
		return
	}
}
