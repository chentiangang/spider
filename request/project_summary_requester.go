package request

import "spider/config"

type ProjectSummaryRequester struct {
	req *APIRequest
}

func NewProjectSummaryRequest(cfg config.RequestConfig) ProjectSummaryRequester {
	return ProjectSummaryRequester{
		req: NewAPIRequest(cfg),
	}
}

func (p ProjectSummaryRequester) SendRequest(cookie string) <-chan []byte {
	data := make(chan []byte)
	go func() {
		p.req.SendRequest(cookie)

	}()
	return nil
}
