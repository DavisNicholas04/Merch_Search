package model

type LogglyStatus struct {
	Method          string `json:"method"`
	SourceIpAddress string `json:"source_ip_address"`
	RequestPath     string `json:"request_path"`
	StatusCode      int    `json:"status_code"`
}
