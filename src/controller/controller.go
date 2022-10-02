package controller

import (
	"encoding/json"
	"fmt"
	"github.com/jamespearly/loggly"
	"io"
	"log"
	"merchSearch/model"
	"merchSearch/service"
	"strconv"
)

func Run() {
	tag := "merch_search"

	service.LoadDotEnv(".env")
	// Instantiate the client
	client := loggly.New(tag)

	response, httpErr := service.GetCompletedAnime("Curiossity")
	if httpErr != nil {
		clientErr := client.EchoSend("error", fmt.Sprintf("Was not able to pull users Completed list\nerror: %s", httpErr))
		if clientErr != nil {
			log.Fatalln(fmt.Sprintf("could not connect to client\nerror: %s", clientErr))
		}
	}
	clientErr := client.EchoSend("info", strconv.FormatInt(response.ContentLength, 10))
	if clientErr != nil {
		log.Fatalln(fmt.Sprintf("could not connect to client\nerror: %s", clientErr))
	}

	bodyBytes, readBytesErr := io.ReadAll(response.Body)
	defer func() {
		err := response.Body.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}()
	if readBytesErr != nil {
		log.Fatalln(readBytesErr)
	}
	var userCompletedResponse model.MalUserListResponse

	err := json.Unmarshal(bodyBytes, &userCompletedResponse)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(userCompletedResponse)
}
