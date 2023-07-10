package models

type Patient struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Age            string `json:"age"`
	Gender         string `json:"gender"`
	Contact        string `json:"contact"`
	MedicalHistory string `json:"medical_history"`
}
