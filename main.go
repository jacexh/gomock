package main

import (
	"flag"
	"net/http"
	"strconv"

	gomock "github.com/jacexh/gomock/lib"
)

func main() {
	port := flag.Int("port", 8080, "the default listen port of gomock")
	flag.Parse()

	http.HandleFunc("/create", gomock.HandleCreate)
	http.HandleFunc("/import", gomock.HandleImport)
	http.HandleFunc("/export", gomock.HandleExport)
	http.HandleFunc("/", gomock.HandleMock)

	addr := strconv.FormatInt(int64(*port), 10)
	gomock.Logger.Info("start gomock on port " + addr)
	http.ListenAndServe(":"+addr, nil)
}
