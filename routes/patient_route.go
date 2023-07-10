package routes

import (
	"net/http"

	"github.com/PATILSHUBHAM69/HC_Patient_Data/controllers"
)

func PatientRoutes() {

	http.HandleFunc("/create_patient", controllers.CreatePatient)
	http.HandleFunc("/get_patient/", controllers.GetPatient)
	http.HandleFunc("/update_patient/", controllers.UpdatePatient)
	http.HandleFunc("/delete_patient/", controllers.DeletePatient)
}
