package service

import (
	"encoding/json"
	"fmt"
	"github.com/jamespearly/loggly"
	"ndavis20_server/handler"
	"ndavis20_server/model"
	"net/http"
	"time"
)

func EncodeJson(writer http.ResponseWriter, statusClient *loggly.ClientType) {
	status := model.Status{SystemTime: time.Now()}
	encoderErr := json.NewEncoder(writer).Encode(status)
	if encoderErr != nil {
		clientErr := statusClient.EchoSend("error", fmt.Sprintf("%v", encoderErr))
		handler.ClientErrorCheck(clientErr)
	}
}
