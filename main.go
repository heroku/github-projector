package main

import (
	"io"
	"net/http"
	"os"
)

func main() {
	var port string
	if port = os.Getenv("PORT"); port == "" {
		port = "8080"
	}

	http.HandleFunc("/webhooks", handleHook)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Work In Progress\n")
		io.WriteString(w, "\n")
		io.WriteString(w, "See github.com/freeformz/github-projector\n")
	})

	http.ListenAndServe(":"+port, nil)
}
