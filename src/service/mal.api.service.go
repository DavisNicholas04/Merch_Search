package service

import (
	"github.com/jamespearly/loggly"
	"merchSearch/controller"
)

func RequestMalBytes(malClient *loggly.ClientType, endpoint string) []byte {
	response, httpErr := controller.GetCompleted(endpoint, "Curiossity")
	HttpErrorCheck(httpErr, malClient)

	return GetBytes(response, malClient)
}
