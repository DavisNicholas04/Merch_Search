package utils

import (
	"github.com/jamespearly/loggly"
	"merchSearch/service"
	"net/http"
	"os"
)

func HttpGet(endPoint string, client *loggly.ClientType) *http.Response {
	hhtpClient := &http.Client{}
	req, _ := http.NewRequest("GET", endPoint, nil)
	req.Header.Set("X-MAL-CLIENT-ID", os.Getenv("X_MAL_CLIENT_ID"))
	response, httpErr := hhtpClient.Do(req)
	service.HttpErrorCheck(httpErr, client)
	return response
}
