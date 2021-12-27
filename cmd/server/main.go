package main

import (
	"log"

	"github.com/joshuaejs/godcls/internal/server"
)

func main() {
	srv := server.NewHTTPServer(":8080")
	log.Fatal(srv.ListenAndServe())
}
