// this package is to be removed upon making the https://github.com/DavisNicholas04/Merch_Search public
// in its place both: github.com/DavisNicholas04/Merch_Search/blob/main/polling_agent/src/handler/error.handler.go
// will be imported

package handler

import (
	"fmt"
	"github.com/jamespearly/loggly"
	"log"
)

func ClientErrorCheck(clientErr error) {
	if clientErr != nil {
		log.Fatalln(fmt.Sprintf("could not connect to client\nerror: %s", clientErr))
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
