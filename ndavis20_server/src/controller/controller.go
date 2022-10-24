package controller

import (
	"encoding/json"
	"fmt"
	"ndavis20_server/handler"
	"ndavis20_server/utils"
	"net/http"
	"os"
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

	switch request.Header.Get("Authorization") {
	case os.Getenv("VERIFICATION_TOKEN"):
		fmt.Println(fmt.Sprintf("Method:%v\nHeader:%v\nBody:%v", request.Method, request.Header, request.Body))
		return
	default:
		http.Error(writer, "Invalid verification token", http.StatusUnauthorized)
	}
}
