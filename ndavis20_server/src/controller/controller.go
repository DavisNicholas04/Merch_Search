package controller

import (
	"ndavis20_server/service"
	"net/http"
	"os"
)

func GetStatus(writer http.ResponseWriter, request *http.Request) {
	service.EncodeJson(writer, "GET_status")
}

func PutEbayDeletionNotification(writer http.ResponseWriter, request *http.Request) {
	// checks if the verification token sent in the header is the correct. Report StatusUnauthorized and returns if not
	if !(request.Header.Get("Authorization") == os.Getenv("VERIFICATION_TOKEN")) {
		http.Error(writer, "Invalid verification token", http.StatusUnauthorized)
		return
	}
}

func GetAll(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte(request.RequestURI + "\n" + request.RemoteAddr + "\n" + request.Host))
}
