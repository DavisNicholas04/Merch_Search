package service

import (
	"net/http"
	"os"
	"regexp"
)

func CheckItemIdReg(itemId string, writer http.ResponseWriter) bool {
	itemIdReg, _ := regexp.Compile("^((v[0-9])\\|([0-9]{1,12})\\|([0-9]{1,12}))$")
	if !itemIdReg.MatchString(itemId) {
		http.Error(
			writer,
			"Malformed request, check your itemId",
			http.StatusBadRequest,
		)
		return false
	}
	return true
}

func CheckUserIdReg(userId string, writer http.ResponseWriter) bool {
	userIdReg, _ := regexp.Compile("[a-zA-z0-9_-]{2,16}")

	if !userIdReg.MatchString(userId) {
		http.Error(
			writer,
			"Malformed request, check your userId",
			http.StatusBadRequest,
		)
		return false
	}
	return true
}

func CheckTableRegex(table string, writer http.ResponseWriter) bool {
	tableReg, _ := regexp.Compile(os.Getenv("DYNAMO_DB_TABLE_NAME"))

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
