package controller

import (
	"fmt"
	"net/http"
	"os"
)

const (
	MalRoot                    = "https://api.myanimelist.net/v2/"
	AnimeListCompletedEndpoint = "users/%s/animelist?status=completed&limit=5"
	//MangaListCompletedEndpoint = "users/%s/mangalist?status=completed"
)

func GetCompleted(endPoint string, user string) (response *http.Response, err error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", MalRoot+fmt.Sprintf(endPoint, user), nil)
	req.Header.Set("X-MAL-CLIENT-ID", os.Getenv("X-MAL-CLIENT-ID"))
	response, err = client.Do(req)
	return
}
