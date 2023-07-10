package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func Connect() {
	var err error

	// Make Database Connection
	db, err = sql.Open("mysql", "root:india@123@tcp(localhost:3306)/hc_patient_data")

	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MySQL database!")

	//Create Tabel
	Create_Patient, err := db.Query("CREATE TABLE IF NOT EXISTS patientdetails (id INT AUTO_INCREMENT PRIMARY KEY,name VARCHAR(80), age VARCHAR(10), gender VARCHAR(10), contact VARCHAR(255), medical_history TEXT );")
	if err != nil {
		panic(err.Error())
	}
	defer Create_Patient.Close()
	fmt.Println("Patient Details Table Successfully Created")
}
