package service

import (
	"errors"
	"fmt"
	"github.com/jamespearly/loggly"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"os"
)

func first(key string, isExist bool) bool {
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
	ReadBytesErrorCheck(readBytesErr, client)

	deferResponseBodyClose(response)

	// send loggly msg
	clientErr := client.EchoSend(
		"info", fmt.Sprintf("statusCode: %v\nresponseSize: %v", response.StatusCode, len(bodyBytes)),
	)
	ClientErrorCheck(clientErr)

	return bodyBytes
}
