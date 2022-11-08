package controller

import (
	"ndavis20_server/service"
	"ndavis20_server/utils"
	"net/http"
	"os"
	"strings"
)

/*
GetStatus : Returns the name of the dynamodb table being searched and the number of items present

Takes in a query_param "liveCount" which can be set to true or false (If liveCount is not set,
it will default to false).

liveCount=false: When liveCount equals false, dynamodb's DescribeTable() is used which returns metrics that
are updated every 6 hours for the specified table. This means that the count returned may not be accurate. Use liveCount=true
if you need live metrics

liveCount=true: When liveCount equals false, dynamodb's Scan() is used which return live metrics for
the specified table.

WARNING:

liveCount=true uses dynamodb's Scan() operation which reads every item in the table. It is expensive and inefficient
and can exhaust a table's read capacity units and throttle user requests. If you do not need a live count DO NOT set
liveCount to true.
*/
func GetStatus(writer http.ResponseWriter, request *http.Request) {
	liveCount := strings.ToLower(request.FormValue("liveCount"))

	if !service.CheckLiveCountRegex(liveCount, writer) {
		return
	}

	status := service.GetStatus(liveCount)
	service.EncodeJson(status, writer, "GET_status")
}

/*
ReceiveEbayDeleteNotif : Receives the Ebay deletion notification, checks the verification token and sends back a
status 200 if they match.

This Function does not use the notification does not use the information sent in the request to delete data
from the dynamodb table. This function is meant to bypass the notification requirement if exemption is not possible.
Requires https with a trusted ssl certificate.
*/
func ReceiveEbayDeleteNotif(writer http.ResponseWriter, request *http.Request) {
	// checks if the verification token sent in the header is the correct. Report StatusUnauthorized and returns if not
	if !(request.Header.Get("Authorization") == os.Getenv("VERIFICATION_TOKEN")) {
		http.Error(writer, "Invalid verification token", http.StatusUnauthorized)
		return
	}
}

// GetAll : Returns all items from the dynamodb table in a json file following the type struct model.UserEntry's format
func GetAll(writer http.ResponseWriter, request *http.Request) {
	table := request.FormValue("table")

	/* COMMENTED OUT TO MEET A10 REQUIREMENTS

	if !service.CheckTableRegex(table, writer) {
		return
	}
	*/

	items := utils.GetAllItems(table)
	service.EncodeJson(items, writer, "GET_getAll")
}

// Search : returns an item from the specified table, given userId and itemId
func Search(writer http.ResponseWriter, request *http.Request) {
	table := request.FormValue("table")
	userId := request.FormValue("user_id")
	itemId := request.FormValue("item_id")

	if !(service.CheckTableRegex(table, writer) && service.CheckUserIdReg(userId, writer) && service.CheckItemIdReg(itemId, writer)) {
		return
	}

	searchedItem := utils.FindItem(table, userId, itemId)

	service.EncodeJson(searchedItem, writer, "GET_search")
}
