package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	ok := server.ListenAndServe()
	if ok != nil {
		fmt.Println(ok.Error())
	}
}
