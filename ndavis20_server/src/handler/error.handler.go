// this package is to be removed upon making the https://github.com/DavisNicholas04/Merch_Search public
// in its place both: github.com/DavisNicholas04/Merch_Search/blob/main/polling_agent/src/handler/error.handler.go
// will be imported

package handler

import (
	"fmt"
	"github.com/jamespearly/loggly"
	"log"
	"runtime"
)

func ReadBytesErrorCheck(readBytesErr error, client *loggly.ClientType) {
	// Get the name of the file that called this function
	_, file, _, _ := runtime.Caller(1)

	// Handle error by sending a log to loggly or log.Fatalln if not possible
	if readBytesErr != nil {
		clientErr := client.EchoSend(
			"error", fmt.Sprintf("Location: %v\n io.Read all could not readBytes \nerror: %s", file, readBytesErr),
		)
		if clientErr != nil {
			log.Fatalln(fmt.Sprintf("could not connect to client\nerror: %s", clientErr))
		}
	}
}

func ClientErrorCheck(clientErr error) {
	if clientErr != nil {
		log.Fatalln(fmt.Sprintf("could not connect to client\nerror: %s", clientErr))
	}
}

func UnmarshalError(unmarshalErr error, client *loggly.ClientType, jsonAttempted string) {
	if unmarshalErr != nil {
		clientErr := client.EchoSend(
			"error", fmt.Sprintf(
				"Was not able to unmarshal the bytes of the response body\nerror: %s\ntitle attempted: %s\n",
				unmarshalErr,
				jsonAttempted),
		)
		if clientErr != nil {
			log.Fatalln(fmt.Sprintf("could not connect to client\nerror: %s", clientErr))
		}
	}
}

func HttpErrorCheck(httpErr error, client *loggly.ClientType) {
	if httpErr != nil {
		clientErr := client.EchoSend("error", fmt.Sprintf("Was not able to connect to endpoint list\nerror: %s", httpErr))
		if clientErr != nil {
			log.Fatalln(fmt.Sprintf("could not connect to client\nerror: %s", clientErr))
		}
	}
}

func MarshalError(unmarshalErr error, client *loggly.ClientType, modelName string) {
	if unmarshalErr != nil {
		clientErr := client.EchoSend(
			"error", fmt.Sprintf(
				"Was not able to marshal the %v: %s\n",
				modelName,
				unmarshalErr,
			),
		)
		if clientErr != nil {
			log.Fatalln(fmt.Sprintf("could not connect to client\nerror: %s", clientErr))
		}
	}
}
