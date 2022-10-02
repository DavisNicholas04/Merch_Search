package service

import (
	"fmt"
	"net/http"
	"os"
)

const root string = "https://api.myanimelist.net/v2/"

func GetCompletedAnime(user string) (response *http.Response, err error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", root+fmt.Sprintf("users/%s/animelist?status=completed", user), nil)
	req.Header.Set("X-MAL-CLIENT-ID", os.Getenv("X-MAL-CLIENT-ID"))

	response, err = client.Do(req)
	return
}
