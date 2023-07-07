package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/PATILSHUBHAM69/HC_Patient_Data/models"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {
	// Initialize the database connection
	var err error
	db, err = sql.Open("mysql", "root:india@123@tcp(localhost:3306)/hc_patient_data")
	if err != nil {
		panic(err)
	}
}

// Rest of the controller code...

func CreatePatient(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	fmt.Println("hello")
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Extract patient details from the request form
	name := r.FormValue("name")
	age, _ := strconv.Atoi(r.FormValue("age"))
	gender := r.FormValue("gender")
	contact := r.FormValue("contact")
	medicalHistory := r.FormValue("medical_history")

	// Perform data validation
	patient := models.Patient{
		Name:           name,
		Age:            age,
		Gender:         gender,
		Contact:        contact,
		MedicalHistory: medicalHistory,
	}

	// Insert the patient record into the database
	stmt, err := db.Prepare("INSERT INTO PatientDetails (name, age, gender, contact_no, medical_history) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	result, err := stmt.Exec(patient.Name, patient.Age, patient.Gender, patient.Contact, patient.MedicalHistory)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the ID of the newly created patient
	patientID, err := result.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the patient ID as the response
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Patient created with ID: %d", patientID)
}

func GetPatient(w http.ResponseWriter, r *http.Request) {
	// Get the patient ID from the request URL
	patientIDStr := r.URL.Path[len("/patients/"):]
	patientID, err := strconv.Atoi(patientIDStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Query the database to retrieve the patient record
	rows, err := db.Query("SELECT id, name, age, gender, contact, medical_history FROM patients WHERE id = ?", patientID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	if rows.Next() {
		// Extract patient details from the query result
		var patient models.Patient
		err := rows.Scan(&patient.ID, &patient.Name, &patient.Age, &patient.Gender, &patient.Contact, &patient.MedicalHistory)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Return the patient details as the response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(patient)
	} else {
		http.NotFound(w, r)
	}
}

func UpdatePatient(w http.ResponseWriter, r *http.Request) {
	// Get the patient ID from the request URL
	patientIDStr := r.URL.Path[len("/patients/"):]
	patientID, err := strconv.Atoi(patientIDStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Parse the request body
	err = r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Extract updated patient details from the request form
	name := r.FormValue("name")
	age, _ := strconv.Atoi(r.FormValue("age"))
	gender := r.FormValue("gender")
	contact := r.FormValue("contact")
	medicalHistory := r.FormValue("medical_history")

	// Create a new patient instance
	patient := models.Patient{
		ID:             patientID,
		Name:           name,
		Age:            age,
		Gender:         gender,
		Contact:        contact,
		MedicalHistory: medicalHistory,
	}

	// Update the patient record in the database
	stmt, err := db.Prepare("UPDATE patients SET name=?, age=?, gender=?, contact=?, medical_history=? WHERE id=?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(patient.Name, patient.Age, patient.Gender, patient.Contact, patient.MedicalHistory, patient.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Patient updated successfully")
}

func DeletePatient(w http.ResponseWriter, r *http.Request) {
	// Get the patient ID from the request URL
	patientIDStr := r.URL.Path[len("/patients/"):]
	patientID, err := strconv.Atoi(patientIDStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create a new patient instance with only the ID field set
	patient := models.Patient{
		ID: patientID,
	}

	// Delete the patient record from the database
	stmt, err := db.Prepare("DELETE FROM patients WHERE id=?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(patient.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Patient deleted successfully")
}
