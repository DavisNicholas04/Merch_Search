package utils

import (
	"github.com/jamespearly/loggly"
	"net/http"
	"os"
	"polling_agent/handler"
)

func HttpGet(endPoint string, client *loggly.ClientType) *http.Response {
	httpClient := &http.Client{}
	req, _ := http.NewRequest("GET", endPoint, nil)
	req.Header.Set("X-MAL-CLIENT-ID", os.Getenv("X_MAL_CLIENT_ID"))
	response, httpErr := httpClient.Do(req)
	handler.HttpErrorCheck(httpErr, client)
	return response
}
