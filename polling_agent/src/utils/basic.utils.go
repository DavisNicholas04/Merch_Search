package utils

import (
	"errors"
	"fmt"
	"github.com/jamespearly/loggly"
	"github.com/joho/godotenv"
	"io"
	"log"
	"merchSearch/handler"
	"net/http"
	"os"
	"regexp"
)

func RemoveNonAlphaNums(str string) string {
	var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9 ]+`)
	return nonAlphanumericRegex.ReplaceAllString(str, "")
}

func Second(key string, isExist bool) bool {
	return isExist
}

func FileExist(file string) bool {
	_, err := os.Stat(file)
	if errors.Is(err, os.ErrNotExist) {
		return false
	} else {
		return true
	}
}

func LoadDotEnv(filenames ...string) {
	for _, name := range filenames {
		if FileExist(name) {
			err1 := godotenv.Load(filenames...)
			if err1 != nil {
				log.Fatalln("Error loading .env file")
			}
		}
	}
}

func deferResponseBodyClose(response *http.Response) {
	defer func() {
		err := response.Body.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}()
}

func GetBytes(response *http.Response, client *loggly.ClientType) []byte {
	bodyBytes, readBytesErr := io.ReadAll(response.Body)
	handler.ReadBytesErrorCheck(readBytesErr, client)

	deferResponseBodyClose(response)

	if response.StatusCode == 200 {
		// send loggly msg
		clientErr := client.EchoSend(
			"info", fmt.Sprintf("statusCode: %v\nresponseSize: %v", response.StatusCode, len(bodyBytes)),
		)
		handler.ClientErrorCheck(clientErr)
	} else {
		clientErr := client.EchoSend(
			"info", fmt.Sprintf("statusCode: %v\nresponseSize: %v\nresponse: ", response.StatusCode, len(bodyBytes), string(bodyBytes)),
		)
		handler.ClientErrorCheck(clientErr)
	}

	return bodyBytes
}
