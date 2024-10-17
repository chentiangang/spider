package parser

import "github.com/chentiangang/xlog"

type Parser interface {
	Parse(response []byte) error
}

func CreateParserWithAPI(api string) Parser {
	switch api {
	case "/gateway/v2/apps/california/project/proprietary_cloud/get_project_info_summary":
		return ProjectSummaryResp{}
	default:
		xlog.Error("Unknown api %s", api)
		return nil
	}
}

type ProjectSummaryResp struct{}

func (p ProjectSummaryResp) Parse(data []byte) error {
	return nil
}
