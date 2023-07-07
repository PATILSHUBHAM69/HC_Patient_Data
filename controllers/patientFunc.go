package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/PATILSHUBHAM69/HC_Patient_Data/models"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("mysql", "root:india@123@tcp(localhost:3306)/hc_patient_data")
	if err != nil {
		panic(err)
	}
}

func CreatePatient(w http.ResponseWriter, r *http.Request) {
	var patient models.Patient
	err := json.NewDecoder(r.Body).Decode(&patient)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	stmt, err := db.Prepare("INSERT INTO PatientDetails (name, age, gender, contact, medical_history) VALUES (?, ?, ?, ?, ?)")
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

	patientID, err := result.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Patient created with ID: %d", patientID)
}

func GetPatient(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		PatientID int `json:"patient_id"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	patientID := requestBody.PatientID

	rows, err := db.Query("SELECT id, name, age, gender, contact, medical_history FROM PatientDetails WHERE id = ?", patientID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	if rows.Next() {
		var patient models.Patient
		err := rows.Scan(&patient.ID, &patient.Name, &patient.Age, &patient.Gender, &patient.Contact, &patient.MedicalHistory)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		jsonData, err := json.Marshal(patient)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(jsonData)
	} else {
		http.NotFound(w, r)
	}
}

func UpdatePatient(w http.ResponseWriter, r *http.Request) {

	var requestBody struct {
		PatientID      int    `json:"patient_id"`
		Name           string `json:"name"`
		Age            int    `json:"age"`
		Gender         string `json:"gender"`
		Contact        string `json:"contact"`
		MedicalHistory string `json:"medical_history"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	patientID := requestBody.PatientID

	patient := models.Patient{
		ID:             patientID,
		Name:           requestBody.Name,
		Age:            requestBody.Age,
		Gender:         requestBody.Gender,
		Contact:        requestBody.Contact,
		MedicalHistory: requestBody.MedicalHistory,
	}

	stmt, err := db.Prepare("UPDATE PatientDetails SET name=?, age=?, gender=?, contact=?, medical_history=? WHERE id=?")
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
	var requestBody struct {
		PatientID int `json:"patient_id"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	patientID := requestBody.PatientID
	patient := models.Patient{
		ID: patientID,
	}

	stmt, err := db.Prepare("DELETE FROM PatientDetails WHERE id=?")
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
