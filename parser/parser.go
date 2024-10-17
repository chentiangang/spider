package parser

import (
	"encoding/json"

	"github.com/chentiangang/xlog"
)

type Parser interface {
	Parse(response []byte) error
}

func CreateParserWithAPI(api string) Parser {
	switch api {
	case "/gateway/v2/apps/california/project/proprietary_cloud/get_project_info_summary":
		return &ProjectSummaryRespParser{}
	default:
		xlog.Error("Unknown api %s", api)
		return nil
	}
}

type ProjectSummaryRespParser struct{}

func (p *ProjectSummaryRespParser) Parse(data []byte) error {
	err := json.Unmarshal(data, &p)
	if err != nil {
		xlog.Error("Unmarshal err %v", err)
		return err
	}
	return nil
}
