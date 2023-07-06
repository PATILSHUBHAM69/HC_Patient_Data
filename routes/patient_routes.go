package routes

import (
	"net/http"
)

func patient_routes(http.ResponseWriter, *http.Request) {
	// Register endpoints
	http.HandleFunc("/create_patient", createPatient)
	http.HandleFunc("/get_patient/", getPatient)
	http.HandleFunc("/update_patient/", updatePatient)
	http.HandleFunc("/delete_patient/", deletePatient)
}
