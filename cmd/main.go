package main

import (
	"fmt"
	"log"
	"net/http"

	"Yandex_Calc/routes"
)

func main() {
	router := routes.SetupRoutes()

	port := ":8080"
	fmt.Printf("Server is running on port %s...\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}
