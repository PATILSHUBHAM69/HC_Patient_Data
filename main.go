package main

import (
	"log"
	"net/http"

	"github.com/PATILSHUBHAM69/HC_Patient_Data/database"
	"github.com/PATILSHUBHAM69/HC_Patient_Data/routes"
)

func main() {
	database.Connect()

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	routes.PatientRoutes()

	log.Fatal(http.ListenAndServe(":8080", nil))
}
