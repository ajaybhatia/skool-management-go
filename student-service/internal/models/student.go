package models

import "time"

type Student struct {
	ID             int       `json:"id" db:"id"`
	RollNumber     string    `json:"roll_number" db:"roll_number"`
	FirstName      string    `json:"first_name" db:"first_name"`
	LastName       string    `json:"last_name" db:"last_name"`
	Email          string    `json:"email" db:"email"`
	Phone          string    `json:"phone" db:"phone"`
	DateOfBirth    string    `json:"date_of_birth" db:"date_of_birth"`
	Address        string    `json:"address" db:"address"`
	SchoolID       int       `json:"school_id" db:"school_id"`
	SchoolName     string    `json:"school_name,omitempty"`
	EnrollmentDate string    `json:"enrollment_date" db:"enrollment_date"`
	Status         string    `json:"status" db:"status"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

type CreateStudentRequest struct {
	RollNumber     string `json:"roll_number"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Email          string `json:"email"`
	Phone          string `json:"phone"`
	DateOfBirth    string `json:"date_of_birth"`
	Address        string `json:"address"`
	SchoolID       int    `json:"school_id"`
	EnrollmentDate string `json:"enrollment_date"`
	Status         string `json:"status"`
}

type UpdateStudentRequest struct {
	RollNumber     string `json:"roll_number"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Email          string `json:"email"`
	Phone          string `json:"phone"`
	DateOfBirth    string `json:"date_of_birth"`
	Address        string `json:"address"`
	SchoolID       int    `json:"school_id"`
	EnrollmentDate string `json:"enrollment_date"`
	Status         string `json:"status"`
}
