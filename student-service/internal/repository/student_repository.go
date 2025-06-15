package repository

import (
	"database/sql"
	"time"

	"skool-management/student-service/internal/models"
)

type StudentRepository struct {
	db *sql.DB
}

func NewStudentRepository(db *sql.DB) *StudentRepository {
	return &StudentRepository{db: db}
}

func (r *StudentRepository) Create(student *models.CreateStudentRequest) (*models.Student, error) {
	query := `
		INSERT INTO students (roll_number, first_name, last_name, email, phone, date_of_birth, address, school_id, enrollment_date, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, roll_number, first_name, last_name, email, phone, date_of_birth, address, school_id, enrollment_date, status, created_at, updated_at
	`

	now := time.Now()
	var result models.Student
	err := r.db.QueryRow(query, student.RollNumber, student.FirstName, student.LastName, 
		student.Email, student.Phone, student.DateOfBirth, student.Address, student.SchoolID, 
		student.EnrollmentDate, student.Status, now, now).Scan(
		&result.ID, &result.RollNumber, &result.FirstName, &result.LastName, &result.Email, 
		&result.Phone, &result.DateOfBirth, &result.Address, &result.SchoolID, 
		&result.EnrollmentDate, &result.Status, &result.CreatedAt, &result.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *StudentRepository) GetAll() ([]models.Student, error) {
	query := `
		SELECT id, roll_number, first_name, last_name, email, phone, date_of_birth, address, school_id, enrollment_date, status, created_at, updated_at
		FROM students
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []models.Student
	for rows.Next() {
		var student models.Student
		err := rows.Scan(
			&student.ID, &student.RollNumber, &student.FirstName, &student.LastName, &student.Email, 
			&student.Phone, &student.DateOfBirth, &student.Address, &student.SchoolID, 
			&student.EnrollmentDate, &student.Status, &student.CreatedAt, &student.UpdatedAt,
		)
		if err != nil {
			continue
		}
		students = append(students, student)
	}

	return students, nil
}

func (r *StudentRepository) GetByID(id int) (*models.Student, error) {
	query := `
		SELECT id, roll_number, first_name, last_name, email, phone, date_of_birth, address, school_id, enrollment_date, status, created_at, updated_at
		FROM students
		WHERE id = $1
	`

	var student models.Student
	err := r.db.QueryRow(query, id).Scan(
		&student.ID, &student.RollNumber, &student.FirstName, &student.LastName, &student.Email, 
		&student.Phone, &student.DateOfBirth, &student.Address, &student.SchoolID, 
		&student.EnrollmentDate, &student.Status, &student.CreatedAt, &student.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &student, nil
}

func (r *StudentRepository) GetBySchoolID(schoolID int) ([]models.Student, error) {
	query := `
		SELECT id, roll_number, first_name, last_name, email, phone, date_of_birth, address, school_id, enrollment_date, status, created_at, updated_at
		FROM students
		WHERE school_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, schoolID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []models.Student
	for rows.Next() {
		var student models.Student
		err := rows.Scan(
			&student.ID, &student.RollNumber, &student.FirstName, &student.LastName, &student.Email, 
			&student.Phone, &student.DateOfBirth, &student.Address, &student.SchoolID, 
			&student.EnrollmentDate, &student.Status, &student.CreatedAt, &student.UpdatedAt,
		)
		if err != nil {
			continue
		}
		students = append(students, student)
	}

	return students, nil
}

func (r *StudentRepository) Update(id int, student *models.UpdateStudentRequest) (*models.Student, error) {
	query := `
		UPDATE students
		SET roll_number = $1, first_name = $2, last_name = $3, email = $4, phone = $5, date_of_birth = $6,
		    address = $7, school_id = $8, enrollment_date = $9, status = $10, updated_at = $11
		WHERE id = $12
		RETURNING id, roll_number, first_name, last_name, email, phone, date_of_birth, address, school_id, enrollment_date, status, created_at, updated_at
	`

	var result models.Student
	err := r.db.QueryRow(query, student.RollNumber, student.FirstName, student.LastName, 
		student.Email, student.Phone, student.DateOfBirth, student.Address, student.SchoolID, 
		student.EnrollmentDate, student.Status, time.Now(), id).Scan(
		&result.ID, &result.RollNumber, &result.FirstName, &result.LastName, &result.Email, 
		&result.Phone, &result.DateOfBirth, &result.Address, &result.SchoolID, 
		&result.EnrollmentDate, &result.Status, &result.CreatedAt, &result.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *StudentRepository) Delete(id int) error {
	query := `DELETE FROM students WHERE id = $1`
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
