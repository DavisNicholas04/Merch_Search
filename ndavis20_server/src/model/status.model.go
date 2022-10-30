package model

type Status struct {
	TableName   string `json:"table_name"`
	RecordCount int64  `json:"record_count"`
}
