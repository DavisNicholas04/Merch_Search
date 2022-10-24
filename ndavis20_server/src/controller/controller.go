package controller

import (
	"encoding/json"
	"fmt"
	"ndavis20_server/handler"
	"ndavis20_server/model"
	"ndavis20_server/service"
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
	logglyModel := model.LogglyStatus{
		Method:          request.Method,
		SourceIpAddress: utils.GetOutboundIP().String(),
		RequestPath:     request.RequestURI,
		StatusCode:      -999,
	}

	switch request.Method {
	case http.MethodGet:
		service.EncodeJson(writer, statusClient)
		logglyModel.StatusCode = http.StatusOK
	default:
		http.Error(writer, "Method Not Allowed", http.StatusMethodNotAllowed)
		logglyModel.StatusCode = http.StatusMethodNotAllowed
	}

	logglyMarshal, err := json.Marshal(logglyModel)
	handler.MarshalError(err, statusClient, "logglyStatusModel")

	clientErr := statusClient.EchoSend("info", string(logglyMarshal))
	handler.ClientErrorCheck(clientErr)
}

func PutEbayDeletionNotification(writer http.ResponseWriter, request *http.Request) {
	statusClient := utils.InstantiateClient("GET_status")
	clientErr := statusClient.EchoSend("info", request.Method)
	handler.ClientErrorCheck(clientErr)

	// check if the ebay deletion notification is a put request and continue if so.
	// else exit the function early and report StatusMethodNotAllowed.
	switch request.Method {
	case http.MethodPut:
		break
	default:
		http.Error(writer, "Invalid verification token", http.StatusMethodNotAllowed)
		return
	}

	// check if the verification token sent in the header is the correct token. If so notify ebay with a 200 OK
	// else report StatusUnauthorized
	switch request.Header.Get("Authorization") {
	case os.Getenv("VERIFICATION_TOKEN"):
		fmt.Println(fmt.Sprintf("Method:%v\nHeader:%v\nBody:%v", request.Method, request.Header, request.Body))
		return
	default:
		http.Error(writer, "Invalid verification token", http.StatusUnauthorized)
	}
}
