// package main

// import (
// 	"log"
// 	"net/http"

// 	"github.com/PATILSHUBHAM69/HC_Patient_Data/controllers"
// 	"github.com/PATILSHUBHAM69/HC_Patient_Data/database"
// )

// func main() {
// 	database.Connect()

// 	// Serve static files (HTML, CSS, JS)
// 	fs := http.FileServer(http.Dir("static"))
// 	http.Handle("/", fs)

// 	// Register API endpoints
// 	http.HandleFunc("/create_patient", controllers.CreatePatient)
// 	http.HandleFunc("/get_patient", controllers.GetPatient)
// 	http.HandleFunc("/update_patient", controllers.UpdatePatient)
// 	http.HandleFunc("/delete_patient", controllers.DeletePatient)

// 	// Start the server
// 	log.Fatal(http.ListenAndServe(":8080", nil))
// }

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
