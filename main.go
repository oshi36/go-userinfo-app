package main

import (
	"fmt"
	"log"
	"net/http"

	"htmlgo/server/router"
	
	
)

func main() {

	r := router.Router()
	
	fmt.Println("Starting server on the port 80...")
	log.Fatal(http.ListenAndServe(":80", r))
}