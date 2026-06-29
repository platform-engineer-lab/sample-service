package main

import (
	"fmt"
	"net/http"
	"os"
)

var version = "dev"

func rootHandler(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w, "sample-service %s\n , Hello World Test 2", version)
}

func healthzHandler(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintln(w, "ok")
}

func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/healthz", healthzHandler)
	port := "8080"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
