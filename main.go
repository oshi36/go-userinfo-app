package main

import (
	"fmt"
	"log"
	"net/http"
	"htmlgo/server/router"
    "go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"os"

)

func main() {
   
   
	r := router.Router()

	r.Use(otelmux.Middleware(os.Getenv("OTEL_SERVICE_NAME")))

	fmt.Println("Starting server on the port 80...")
	log.Fatal(http.ListenAndServe(":80", r))



}