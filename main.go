package main

import (
	"log"
	"net/http"

	"github.com/PATILSHUBHAM69/HC_Patient_Data/database"
	"github.com/PATILSHUBHAM69/HC_Patient_Data/routes"
)

func main() {
	database.Connect()

	http.HandleFunc("/", routes.PatientRoutes)

	// Start the server
	log.Fatal(http.ListenAndServe(":8080", nil))
}
