package main

import (
	"log"
	"microservice-template/common"
	"net/http"

	_ "github.com/joho/godotenv/autoload" // package that loads env
)

var (
	httpAddr = common.EnvString("PORT", "2222")
)

func main() {

	// routes
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/customers/{customerID}/orders", handleCreateOrder)

	// start server
	log.Printf("Server started on port %s", httpAddr)

	if err := http.ListenAndServe(":"+httpAddr, mux); err != nil {
		log.Fatal("Failed to start server")
	}
}

/* placed here for ease of use during learning  */
func handleCreateOrder(w http.ResponseWriter, r *http.Request) {

}
