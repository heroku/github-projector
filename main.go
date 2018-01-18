package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"os"
)

func handleHook(w http.ResponseWriter, r *http.Request) {
	b, err := httputil.DumpRequest(r, true)
	if err != nil {
		log.Println("error dumping request:", err)
		return
	}
	log.Printf("%s", b)
}

func main() {
	var port string
	if port = os.Getenv("PORT"); port == "" {
		port = "8080"
	}

	http.HandleFunc("/webhooks", handleHook)

	http.ListenAndServe(":"+port, nil)
}
