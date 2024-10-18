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

//type ProjectSummaryResponse struct {
//	Code    string
//	Message string
//	Data    interface{}
//}

type ProjectSummaryResponse struct {
	Code int `json:"code"`
	Data struct {
		Items []Item `json:"items"`
		Total int    `json:"total"`
	}
	Message string `json:"message"`
}

type Item struct {
	Bu                    string      `json:"bu"`
	Cluster               string      `json:"cluster"`
	CpuMax                int         `json:"cpu_max"`
	CpuMin                int         `json:"cpu_min"`
	CpuUse                int         `json:"cpu_use"`
	CpuUseRate            string      `json:"cpu_use_rate"`
	CpuUseRateYesterday   string      `json:"cpu_use_rate_yesterday"`
	CreateTime            string      `json:"create_time"`
	DiskTotal             int64       `json:"disk_total"`
	DiskTotalUnitTB       int         `json:"disk_total_unit_tb"`
	DiskUse               int64       `json:"disk_use"`
	DiskUseRate           string      `json:"disk_use_rate"`
	DiskUseUnitTB         int         `json:"disk_use_unit_tb"`
	FileNumberUseRate     string      `json:"file_number_use_rate"`
	GpuUseRateYesterday   interface{} `json:"gpu_use_rate_yesterday"`
	MemoryMax             int64       `json:"memory_max"`
	MemoryMin             int64       `json:"memory_min"`
	MemoryUse             int         `json:"memory_use"`
	MemoryUseRate         string      `json:"memory_use_rate"`
	Project               string      `json:"project"`
	ProjectDescription    interface{} `json:"project_description"`
	ProjectMaxFileNumber  int64       `json:"project_max_file_number"`
	ProjectNameCn         string      `json:"project_name_cn"`
	ProjectOwner          string      `json:"project_owner"`
	ProjectUsedFileNumber int64       `json:"project_used_file_number"`
	Quota                 string      `json:"quota"`
	QuotaID               int         `json:"quota_id"`
	Region                string      `json:"region"`
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
		defer close(data)

		bs, err := h.req.SendRequest(cookie)
		if err != nil {
			xlog.Error("Failed to send request cookie: %s, err: %s", cookie, err)
			return
		}

		res := h.parse(bs)
		//res.Data.Total

		data <- bs
	}()

	return data
}

func (h *ProjectSummaryHandler) parse(data []byte) ProjectSummaryResponse {
	var resp ProjectSummaryResponse
	err := json.Unmarshal(data, &resp)
	if err != nil {
		xlog.Error("%s", err)
	}
	return resp
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
