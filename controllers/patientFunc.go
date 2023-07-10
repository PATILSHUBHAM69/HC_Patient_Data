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
	patientID := r.URL.Query().Get("id")
	if patientID == "" {
		http.Error(w, "Missing patient ID", http.StatusBadRequest)
		return
	}

	patientIDInt, err := strconv.Atoi(patientID)
	if err != nil {
		http.Error(w, "Invalid patient ID", http.StatusBadRequest)
		return
	}

	rows, err := db.Query("SELECT id, name, age, gender, contact, medical_history FROM PatientDetails WHERE id = ?", patientIDInt)
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
	patientIDStr := r.URL.Query().Get("id")
	patientID, err := strconv.Atoi(patientIDStr)
	if err != nil {
		http.Error(w, "Invalid patient ID", http.StatusBadRequest)
		return
	}

	var patient models.Patient
	err = json.NewDecoder(r.Body).Decode(&patient)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	stmt, err := db.Prepare("UPDATE PatientDetails SET name=?, age=?, gender=?, contact=?, medical_history=? WHERE id=?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(patient.Name, patient.Age, patient.Gender, patient.Contact, patient.MedicalHistory, patientID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Patient updated successfully")
}

func DeletePatient(w http.ResponseWriter, r *http.Request) {
	patientIDStr := r.URL.Query().Get("id")
	patientID, err := strconv.Atoi(patientIDStr)
	if err != nil {
		http.Error(w, "Invalid patient ID", http.StatusBadRequest)
		return
	}

	stmt, err := db.Prepare("DELETE FROM PatientDetails WHERE id=?")
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
