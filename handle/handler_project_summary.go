package handle

import (
	"database/sql"
	"encoding/json"
	"spider/config"

	"github.com/chentiangang/xlog"
)

type ProjectSummaryHandler struct {
	//Resp   ProjectSummaryResponse
	RespCh chan ProjectSummaryResponse
	db     *sql.DB
	req    *Request
}

type ProjectSummaryResponse struct {
	Code    string
	Message string
	Data    interface{}
}

func (h *ProjectSummaryHandler) Init(cfg config.TaskConfig) error {
	h.RespCh = make(chan ProjectSummaryResponse)
	var err error
	h.req, err = NewRequest(cfg.Request)
	h.db = NewConn(cfg.Storage)
	if err != nil {
		xlog.Error("Failed to init Handler %s, err: %s:", h.Name(), err)
		return err
	}

	return nil
}

func (h *ProjectSummaryHandler) Name() string {
	return "ProjectSummaryHandler"
}

func (h *ProjectSummaryHandler) SendRequest(cookie string) <-chan []byte {
	data := make(chan []byte)

	go func() {
		for {
			bs, err := h.req.SendRequest(cookie)
			if err != nil {
				xlog.Error("Failed to send request cookie: %s, err: %s", cookie, err)
				return
			}
			data <- bs
		}
	}()

	return data
}

func (h *ProjectSummaryHandler) ParseToChan(data <-chan []byte) {
	go func() {
		defer close(h.RespCh)
		for bs := range data {
			var resp ProjectSummaryResponse
			err := json.Unmarshal(bs, &resp)
			if err != nil {
				xlog.Error("%s", err)
			}
			h.RespCh <- resp
		}
	}()
}

func (h *ProjectSummaryHandler) Store() {
	go func() {
		defer h.db.Close()
		for item := range h.RespCh {
			_, err := h.db.Exec("insert into", item.Data)
			if err != nil {
				xlog.Error("%s", err)
			}
		}
	}()
}
