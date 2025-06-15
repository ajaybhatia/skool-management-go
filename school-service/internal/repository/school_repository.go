package repository

import (
	"database/sql"
	"time"

	"skool-management/school-service/internal/models"
)

type SchoolRepository struct {
	db *sql.DB
}

func NewSchoolRepository(db *sql.DB) *SchoolRepository {
	return &SchoolRepository{db: db}
}

func (r *SchoolRepository) Create(school *models.CreateSchoolRequest) (*models.School, error) {
	query := `
		INSERT INTO schools (registration_number, name, address, phone, email, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, registration_number, name, address, phone, email, created_at, updated_at
	`

	now := time.Now()
	var result models.School
	err := r.db.QueryRow(query, school.RegistrationNumber, school.Name, school.Address, 
		school.Phone, school.Email, now, now).Scan(
		&result.ID, &result.RegistrationNumber, &result.Name, &result.Address, 
		&result.Phone, &result.Email, &result.CreatedAt, &result.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *SchoolRepository) GetAll() ([]models.School, error) {
	query := `
		SELECT id, registration_number, name, address, phone, email, created_at, updated_at
		FROM schools
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schools []models.School
	for rows.Next() {
		var school models.School
		err := rows.Scan(
			&school.ID, &school.RegistrationNumber, &school.Name, &school.Address, 
			&school.Phone, &school.Email, &school.CreatedAt, &school.UpdatedAt,
		)
		if err != nil {
			continue
		}
		schools = append(schools, school)
	}

	return schools, nil
}

func (r *SchoolRepository) GetByID(id int) (*models.School, error) {
	query := `
		SELECT id, registration_number, name, address, phone, email, created_at, updated_at
		FROM schools
		WHERE id = $1
	`

	var school models.School
	err := r.db.QueryRow(query, id).Scan(
		&school.ID, &school.RegistrationNumber, &school.Name, &school.Address, 
		&school.Phone, &school.Email, &school.CreatedAt, &school.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &school, nil
}

func (r *SchoolRepository) Update(id int, school *models.UpdateSchoolRequest) (*models.School, error) {
	query := `
		UPDATE schools
		SET registration_number = $1, name = $2, address = $3, phone = $4, email = $5, updated_at = $6
		WHERE id = $7
		RETURNING id, registration_number, name, address, phone, email, created_at, updated_at
	`

	var result models.School
	err := r.db.QueryRow(query, school.RegistrationNumber, school.Name, school.Address, 
		school.Phone, school.Email, time.Now(), id).Scan(
		&result.ID, &result.RegistrationNumber, &result.Name, &result.Address, 
		&result.Phone, &result.Email, &result.CreatedAt, &result.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *SchoolRepository) Delete(id int) error {
	query := `DELETE FROM schools WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
