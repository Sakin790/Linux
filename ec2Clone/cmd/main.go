package main

import (
	"ece2clone/handlers"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/instances", handlers.ListInstancesHandler)











	
	addr := ":8080"
	log.Printf("Incus-api server listening on %s", addr)
	log.Printf("Try http://localhost:8080/api/v1/instances")
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
