package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"server/handler"
	"server/utils"
	"time"
)

type Status struct {
	SystemTime time.Time `json:"system_time"`
}

func GetStatus(writer http.ResponseWriter, request *http.Request) {
	statusClient := utils.InstantiateClient("GET_status")
	switch request.Method {
	case http.MethodGet:
		writer.Header().Set("Content-Type", "application/json")
		encoderErr := json.NewEncoder(writer).Encode(Status{SystemTime: time.Now()})
		if encoderErr != nil {
			clientErr := statusClient.EchoSend("error", fmt.Sprintf("%v", encoderErr))
			handler.ClientErrorCheck(clientErr)
		}
	default:
		http.Error(writer, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
	clientErr := statusClient.EchoSend("info", request.Method)
	handler.ClientErrorCheck(clientErr)
}

func PutEbayDeletionNotification(writer http.ResponseWriter, request *http.Request) {
	statusClient := utils.InstantiateClient("GET_status")
	clientErr := statusClient.EchoSend("info", request.Method)
	handler.ClientErrorCheck(clientErr)
	fmt.Println(fmt.Sprintf("Method:%v\nHeader:%v\nBody:%v", request.Method, request.Header, request.Body))
}
