package main

import (
	"net/http"
)

type Mux struct {
	DDay, API *http.ServeMux
}
