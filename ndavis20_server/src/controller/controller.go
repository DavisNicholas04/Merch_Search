package controller

import (
	"ndavis20_server/model"
	"ndavis20_server/service"
	"ndavis20_server/utils"
	"net/http"
	"os"
	"time"
)

func GetStatus(writer http.ResponseWriter, request *http.Request) {
	status := model.Status{SystemTime: time.Now()}
	service.EncodeJson(status, writer, "GET_status")
}

func PutEbayDeletionNotification(writer http.ResponseWriter, request *http.Request) {
	// checks if the verification token sent in the header is the correct. Report StatusUnauthorized and returns if not
	if !(request.Header.Get("Authorization") == os.Getenv("VERIFICATION_TOKEN")) {
		http.Error(writer, "Invalid verification token", http.StatusUnauthorized)
		return
	}
}

func GetAll(writer http.ResponseWriter, request *http.Request) {
	table := request.FormValue("table")

	if !utils.CheckTableRegex(table, writer) {
		return
	}

	items := utils.GetAllItems(table)
	service.EncodeJson(items, writer, "GET_getAll")
}

func Search(writer http.ResponseWriter, request *http.Request) {
	table := request.FormValue("table")
	itemId := request.FormValue("item_id")
	userId := request.FormValue("user_id")

	if !utils.CheckAgainstRegSearch(table, itemId, userId, writer) {
		return
	}

	searchedItem := utils.FindItem(table, userId, itemId)

	service.EncodeJson(searchedItem, writer, "GET_search")
}
