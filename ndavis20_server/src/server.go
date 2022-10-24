package main

import (
	"log"
	"ndavis20_server/controller"
	"ndavis20_server/utils"
	"net/http"
)

func main() {
	utils.LoadDotEnv("../.env")

	http.HandleFunc("/ndavis20/status", controller.GetStatus)

	http.HandleFunc("/ndavis20/ebay/deletion_notification", controller.PutEbayDeletionNotification)
	log.Fatalln(http.ListenAndServeTLS(":45273", "../etc/ssl/certs/34.207.90.86.crt", "../etc/ssl/certs/34.207.90.86.key", nil))
}
