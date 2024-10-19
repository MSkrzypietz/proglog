package main

import (
	"github.com/MSkrzypietz/proglog/internal/httpserver"
	"log"
)

func main() {
	srv := httpserver.NewHttpServer(":8080")
	log.Fatal(srv.ListenAndServe())
}
