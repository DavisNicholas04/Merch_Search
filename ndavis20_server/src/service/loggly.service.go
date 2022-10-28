package service

import (
	"encoding/json"
	"ndavis20_server/handler"
	"ndavis20_server/model"
	"ndavis20_server/utils"
	"net/http"
	"time"
)

func EncodeJson(writer http.ResponseWriter, logglyTag string) {
	status := model.Status{SystemTime: time.Now()}
	encoderErr := json.NewEncoder(writer).Encode(status)
	if encoderErr != nil {
		statusClient := utils.InstantiateClient(logglyTag)
		clientErr := statusClient.EchoSend("error", encoderErr.Error())
		handler.ClientErrorCheck(clientErr)
	}
}

func SendStatus(request *http.Request, status int) {
	statusClient := utils.InstantiateClient("controller")
	logglyModel := model.LogglyStatus{
		Method:          request.Method,
		SourceIpAddress: utils.GetOutboundIP().String(),
		RequestPath:     request.RequestURI,
		StatusCode:      status,
	}

	logglyMarshal, err := json.Marshal(logglyModel)
	handler.MarshalError(err, statusClient, "logglyStatusModel")

	clientErr := statusClient.EchoSend("info", string(logglyMarshal))
	handler.ClientErrorCheck(clientErr)
}
