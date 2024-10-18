package handle

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"spider/config"
	"spider/utils"

	"github.com/chentiangang/xlog"
)

type ProjectSummaryHandler struct {
	//Resp   ProjectSummaryResponse
	reqConfig config.RequestConfig
	RespCh    chan ProjectSummaryResponse
	db        *sql.DB
	req       *Request
}

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
	h.reqConfig = cfg.Request
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
		var err error
		h.req, err = NewRequest(h.reqConfig)
		bs, err := h.req.SendRequest(cookie)
		if err != nil {
			xlog.Error("Failed to send request cookie: %s, err: %s", cookie, err)
			return
		}

		res := h.parse(bs)

		h.reqConfig.Params["pageNum"] = "1"
		h.reqConfig.Params["pageSize"] = fmt.Sprintf("%d", res.Data.Total)
		h.req, err = NewRequest(h.reqConfig)
		if err != nil {
			xlog.Error("%s", err)
			return
		}

		bs, err = h.req.SendRequest(cookie)
		if err != nil {
			xlog.Error("%s")
			return
		}
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
		for {
			bs, ok := <-data
			if !ok {
				xlog.Debug("Channel is closed,no more data need parse.")
				break
			}
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
		for {
			res, ok := <-h.RespCh
			if !ok {
				xlog.Debug("Channel is closed,no more data need store.")
				break
			}
			for _, i := range res.Data.Items {
				_, err := h.db.Exec("insert into max_compute_project_summary(`project`,`region_cluster`,`quota`,`disk_use`,`disk_use_rate`,`disk_total`) values(?,?,?,?,?,?)",
					i.Project, "hb-1", i.Quota, utils.ConvertBytesToReadable(i.DiskUse), i.DiskUseRate, utils.ConvertBytesToReadable(i.DiskTotal))
				if err != nil {
					xlog.Error("%s", err)
				}
			}
		}
	}()
}
