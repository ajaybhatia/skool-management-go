package service

import (
	"database/sql"
	"errors"
	"strings"

	"skool-management/school-service/internal/models"
	"skool-management/school-service/internal/repository"
)

type SchoolService struct {
	schoolRepo *repository.SchoolRepository
}

func NewSchoolService(schoolRepo *repository.SchoolRepository) *SchoolService {
	return &SchoolService{
		schoolRepo: schoolRepo,
	}
}

func (s *SchoolService) CreateSchool(req *models.CreateSchoolRequest) (*models.School, error) {
	if req.Name == "" {
		return nil, errors.New("school name is required")
	}

	if req.RegistrationNumber == "" {
		return nil, errors.New("school registration number is required")
	}

	school, err := s.schoolRepo.Create(req)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return nil, errors.New("school with this registration number already exists")
		}
		return nil, errors.New("failed to create school")
	}

	return school, nil
}

func (s *SchoolService) GetAllSchools() ([]models.School, error) {
	return s.schoolRepo.GetAll()
}

func (s *SchoolService) GetSchoolByID(id int) (*models.School, error) {
	return s.schoolRepo.GetByID(id)
}

func (s *SchoolService) UpdateSchool(id int, req *models.UpdateSchoolRequest) (*models.School, error) {
	if req.RegistrationNumber == "" || req.Name == "" {
		return nil, errors.New("registration number and name are required")
	}

	school, err := s.schoolRepo.Update(id, req)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("school not found")
		}
		if strings.Contains(err.Error(), "duplicate key value") {
			if strings.Contains(err.Error(), "schools_registration_number_key") {
				return nil, errors.New("school with this registration number already exists")
			} else {
				return nil, errors.New("school with this email already exists")
			}
		}
		return nil, errors.New("failed to update school")
	}

	return school, nil
}

func (s *SchoolService) DeleteSchool(id int) error {
	err := s.schoolRepo.Delete(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("school not found")
		}
		return errors.New("failed to delete school")
	}
	return nil
}
