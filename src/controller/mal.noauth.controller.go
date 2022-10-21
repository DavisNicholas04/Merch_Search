package controller


import (
	"net/http"
	"time"
)

func SetUp() {
	http.HandleFunc("/ndavis20", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(time.Now().String()))
		if err != nil {
			return
		}
	})
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}

}
