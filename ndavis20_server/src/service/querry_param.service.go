package service

import (
	"ndavis20_server/model"
	"ndavis20_server/utils"
	"net/http"
	"os"
	"regexp"
)

func CheckItemIdReg(itemId string, writer http.ResponseWriter) bool {
	itemIdReg, _ := regexp.Compile("^((v[0-9])\\|([0-9]{1,12})\\|([0-9]{1,12}))$")
	if !itemIdReg.MatchString(itemId) {
		http.Error(
			writer, "Malformed request, check your itemId", http.StatusBadRequest)
		return false
	}
	return true
}

func CheckUserIdReg(userId string, writer http.ResponseWriter) bool {
	userIdReg, _ := regexp.Compile("^[a-zA-z0-9_-]{2,16}$")

	if !userIdReg.MatchString(userId) {
		http.Error(
			writer, "Malformed request, check your userId", http.StatusBadRequest)
		return false
	}
	return true
}

func CheckTableRegex(table string, writer http.ResponseWriter) bool {
	tableReg, _ := regexp.Compile("^" + os.Getenv("DYNAMO_DB_TABLE_NAME") + "$")

	if !tableReg.MatchString(table) {
		http.Error(
			writer,
			"You either do not have permission to access this table "+
				"or the table does not exist",
			http.StatusUnauthorized,
		)
		return false
	}
	return true
}

func CheckLiveCountRegex(liveCount string, writer http.ResponseWriter) bool {
	liveCountReg, _ := regexp.Compile("^(true|false)$")
	if !liveCountReg.MatchString(liveCount) && liveCount != "" {
		http.Error(writer, "Malformed request, liveCount can only equal true or false", http.StatusBadRequest)
		return false
	}
	return true
}

func getItemCount(tableName string, liveCount string) int64 {
	var itemCount int64
	if liveCount == "true" {
		itemCount = *utils.GetLiveItemCount(tableName)
	} else {
		itemCount = *utils.GetItemCount(tableName)
	}
	return itemCount
}

func GetStatus(liveCount string) model.Status {
	tableName := os.Getenv("DYNAMO_DB_TABLE_NAME")

	var status model.Status

	itemCount := getItemCount(tableName, liveCount)

	status = model.Status{
		TableName:   tableName,
		RecordCount: itemCount,
	}

	return status
}
