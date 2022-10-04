package controller

import (
	"fmt"
	"net/http"
	"os"
)

func GetCompleted(endPoint string, user string) (response *http.Response, err error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", MalRoot+fmt.Sprintf(endPoint, user), nil)
	req.Header.Set("X-MAL-CLIENT-ID", os.Getenv("X-MAL-CLIENT-ID"))
	response, err = client.Do(req)
	return
}
