package main

import (
	"log"
	"net/http"
	"server/controller"
	"server/utils"
)

func main() {
	utils.LoadDotEnv("../../.env")

	http.HandleFunc("/ndavis20/status", controller.GetStatus)

	http.HandleFunc("/ndavis20/ebay/deletion_notification", controller.PutEbayDeletionNotification)
	log.Fatalln(http.ListenAndServe(":45273", nil))
}
