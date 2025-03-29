package main

import (
	"log"
	"net/http"
)

const (
	httpAddr = ":1212"
)

func main() {
	// routes

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/customers/{customerID}/orders", handleCreateOrder)

	// start server
	log.Printf("Server started on port %s", httpAddr)

	if err := http.ListenAndServe(httpAddr, mux); err != nil {
		log.Fatal("Failed to start server")
	}
}

/* placed here for ease of use during learning  */
func handleCreateOrder(w http.ResponseWriter, r *http.Request) {

}
