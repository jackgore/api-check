package main

import (
	"fmt"
	"log"
	"net/http"
)

func fooHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "bar")
}

func main() {
	http.HandleFunc("/foo", fooHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
