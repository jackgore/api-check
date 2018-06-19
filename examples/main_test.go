package main

import (
	"net/http"
	"testing"

	"github.com/JonathonGore/api-check/suite"
)

func setup() {
	http.HandleFunc("/foo", fooHandler)

	go http.ListenAndServe(":8080", nil)
}

func cleanup() {
	// Clean up test server here
}

func TestAPICheck(t *testing.T) {
	setup()
	suite.Run(t)
	cleanup()
}
