package odscraper

import (
	"net/http"
	"time"
)

// Timeout is the timeout on getting player info from Overbuff
var Timeout = 5

// Get a client that won't sit and wait for the server forever
func saneClient() *http.Client {
	return &http.Client{Timeout: time.Duration(Timeout) * time.Second}
}
