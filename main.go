package main

import (
	"net/http"

	gomock "github.com/jacexh/gomock/lib"
)

func main() {
	http.HandleFunc("/", gomock.HandleMock)
	http.HandleFunc("/create", gomock.HandleCreate)
	http.HandleFunc("/import", gomock.HandleImport)
	http.HandleFunc("/export", gomock.HandleExport)

	http.ListenAndServe(":8080", nil)
}
