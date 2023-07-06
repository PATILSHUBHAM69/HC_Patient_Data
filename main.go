package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

// Patient struct represents a patient record
type Patient struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Age            int    `json:"age"`
	Gender         string `json:"gender"`
	Contact        string `json:"contact"`
	MedicalHistory string `json:"medical_history"`
}

// Database connection parameters
const (
	dbDriver = "mysql"
	dbUser   = "root"
	dbPass   = "india@123"
	dbName   = "database_name"
)

var db *sql.DB

func main() {
	// Initialize the database connection
	initDB()

	// Register endpoints
	http.HandleFunc("/patients", createPatient)
	http.HandleFunc("/patients/", getPatient)
	http.HandleFunc("/patients/", updatePatient)
	http.HandleFunc("/patients/", deletePatient)

	// Start the server
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func initDB() {
	var err error
	db, err = sql.Open(dbDriver, fmt.Sprintf("%s:%s@/%s", dbUser, dbPass, dbName))
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to the database")
}

func createPatient(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
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

	// Insert the patient record into the database
	stmt, err := db.Prepare("INSERT INTO patients (name, age, gender, contact, medical_history) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	result, err := stmt.Exec(name, age, gender, contact, medicalHistory)
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

	// Store the patientID in the same table
	updateStmt, err := db.Prepare("UPDATE patients SET id=? WHERE id=?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer updateStmt.Close()

	_, err = updateStmt.Exec(patientID, patientID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the patient ID as the response
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Patient created with ID: %d", patientID)
}

func getPatient(w http.ResponseWriter, r *http.Request) {
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
		var patient Patient
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

func updatePatient(w http.ResponseWriter, r *http.Request) {
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

	// Update the patient record in the database
	stmt, err := db.Prepare("UPDATE patients SET name=?, age=?, gender=?, contact=?, medical_history=? WHERE id=?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(name, age, gender, contact, medicalHistory, patientID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Patient updated successfully")
}

func deletePatient(w http.ResponseWriter, r *http.Request) {
	// Get the patient ID from the request URL
	patientIDStr := r.URL.Path[len("/patients/"):]
	patientID, err := strconv.Atoi(patientIDStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Delete the patient record from the database
	stmt, err := db.Prepare("DELETE FROM patients WHERE id=?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(patientID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Patient deleted successfully")
}
