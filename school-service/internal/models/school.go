package models

import "time"

type School struct {
	ID                 int       `json:"id" db:"id"`
	RegistrationNumber string    `json:"registration_number" db:"registration_number"`
	Name               string    `json:"name" db:"name"`
	Address            string    `json:"address" db:"address"`
	Phone              string    `json:"phone" db:"phone"`
	Email              string    `json:"email" db:"email"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
}

type CreateSchoolRequest struct {
	RegistrationNumber string `json:"registration_number"`
	Name               string `json:"name"`
	Address            string `json:"address"`
	Phone              string `json:"phone"`
	Email              string `json:"email"`
}

type UpdateSchoolRequest struct {
	RegistrationNumber string `json:"registration_number"`
	Name               string `json:"name"`
	Address            string `json:"address"`
	Phone              string `json:"phone"`
	Email              string `json:"email"`
}
