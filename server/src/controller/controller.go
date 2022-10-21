package controller

import (
	"encoding/json"
	"net/http"
	"server/utils"
	"time"
)

type Status struct {
	SystemTime time.Time `json:"system_time"`
}

func GetStatus(writer http.ResponseWriter, request *http.Request) {
	utils.InstantiateClient("GET_status")
	switch request.Method {
	case http.MethodGet:
		writer.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(writer).Encode(Status{SystemTime: time.Now()})
		if err != nil {
			return
		}
	default:
		http.Error(writer, "Method Not Allowed", http.StatusMethodNotAllowed)
	}

}
