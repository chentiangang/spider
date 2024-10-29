package handle

type Disk struct {
	Code int      `json:"code"`
	D    DiskData `json:"data"`
}

type DiskData struct {
	TotalDiskSize     [][]interface{} `json:"total_disk_size"`
	UsedRate          [][]interface{} `json:"used_rate"`
	TotalUsedDiskSize [][]interface{} `json:"total_used_disk_size"`
}
