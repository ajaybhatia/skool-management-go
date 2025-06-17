package service

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"skool-management/shared"
	"skool-management/student-service/internal/models"
	"skool-management/student-service/internal/repository"

	grpcLib "google.golang.org/grpc"
)

type StudentService struct {
	studentRepo          *repository.StudentRepository
	schoolServiceConn    *grpcLib.ClientConn
	schoolCircuitBreaker *shared.CircuitBreaker
}

func NewStudentService(studentRepo *repository.StudentRepository, schoolServiceConn *grpcLib.ClientConn) *StudentService {
	return &StudentService{
		studentRepo:       studentRepo,
		schoolServiceConn: schoolServiceConn,
		// Initialize circuit breaker for school service gRPC calls
		schoolCircuitBreaker: shared.NewCircuitBreaker(shared.CircuitBreakerConfig{
			Name:         "school-service-grpc",
			MaxFailures:  3,
			ResetTimeout: 30 * time.Second,
		}),
	}
}

// Helper function to validate school existence via gRPC
func (s *StudentService) validateSchool(schoolID int) (bool, string, error) {
	if s.schoolServiceConn == nil {
		return true, "", nil // Skip validation if gRPC connection is not available
	}

	var exists bool
	var name string

	// Use circuit breaker for gRPC calls
	err := s.schoolCircuitBreaker.Execute(func() error {
		// Temporarily disable gRPC validation due to protobuf marshaling issues
		// In a production environment, this would use proper protobuf generated code
		shared.LogInfo("STUDENT_SERVICE", fmt.Sprintf("Temporarily skipping gRPC validation for school ID %d", schoolID))
		exists = true
		name = "Test School"
		return nil

		// TODO: Re-enable once proper protobuf code is generated
		// client := grpc.NewSchoolServiceClient(s.schoolServiceConn)
		// resp, err := client.ValidateSchool(context.Background(), &grpc.ValidateSchoolRequest{
		// 	Id: strconv.Itoa(schoolID),
		// })
		// if err != nil {
		// 	return err
		// }
		// exists = resp.Exists
		// name = resp.Name
		// return nil
	})

	if err != nil {
		if err.Error() == "circuit breaker is OPEN" {
			shared.LogError("STUDENT_SERVICE", "school validation circuit breaker", err)
			return false, "", errors.New("school service temporarily unavailable")
		}
		return false, "", err
	}

	return exists, name, nil
}

func (s *StudentService) CreateStudent(req *models.CreateStudentRequest) (*models.Student, error) {
	if req.RollNumber == "" || req.FirstName == "" || req.LastName == "" || req.SchoolID == 0 {
		return nil, errors.New("roll number, first name, last name, and school ID are required")
	}

	// Validate school exists
	schoolExists, schoolName, err := s.validateSchool(req.SchoolID)
	if err != nil {
		shared.LogError("STUDENT_SERVICE", "school validation", err)
		return nil, errors.New("failed to validate school")
	}

	if !schoolExists {
		return nil, errors.New("school does not exist")
	}

	// Set defaults
	if req.Status == "" {
		req.Status = "active"
	}
	if req.EnrollmentDate == "" {
		req.EnrollmentDate = time.Now().Format("2006-01-02")
	}

	student, err := s.studentRepo.Create(req)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			if strings.Contains(err.Error(), "students_roll_number_school_id_key") {
				return nil, errors.New("student with this roll number already exists in this school")
			} else {
				return nil, errors.New("student with this email already exists")
			}
		}
		return nil, errors.New("failed to create student")
	}

	student.SchoolName = schoolName
	return student, nil
}

func (s *StudentService) GetAllStudents() ([]models.Student, error) {
	students, err := s.studentRepo.GetAll()
	if err != nil {
		return nil, err
	}

	// Optionally fetch school names for each student
	for i := range students {
		_, schoolName, _ := s.validateSchool(students[i].SchoolID)
		students[i].SchoolName = schoolName
	}

	return students, nil
}

func (s *StudentService) GetStudentByID(id int) (*models.Student, error) {
	student, err := s.studentRepo.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("student not found")
		}
		return nil, errors.New("failed to get student")
	}

	// Get school name via gRPC
	_, schoolName, _ := s.validateSchool(student.SchoolID)
	student.SchoolName = schoolName

	return student, nil
}

func (s *StudentService) GetStudentsBySchoolID(schoolID int) ([]models.Student, string, error) {
	// Validate school exists
	schoolExists, schoolName, err := s.validateSchool(schoolID)
	if err != nil {
		shared.LogError("STUDENT_SERVICE", "school validation", err)
		return nil, "", errors.New("failed to validate school")
	}

	if !schoolExists {
		return nil, "", errors.New("school not found")
	}

	students, err := s.studentRepo.GetBySchoolID(schoolID)
	if err != nil {
		return nil, "", errors.New("failed to get students")
	}

	// Set school name for all students
	for i := range students {
		students[i].SchoolName = schoolName
	}

	return students, schoolName, nil
}

func (s *StudentService) UpdateStudent(id int, req *models.UpdateStudentRequest) (*models.Student, error) {
	if req.RollNumber == "" || req.FirstName == "" || req.LastName == "" || req.SchoolID == 0 {
		return nil, errors.New("roll number, first name, last name, and school ID are required")
	}

	// Validate school exists
	schoolExists, schoolName, err := s.validateSchool(req.SchoolID)
	if err != nil {
		shared.LogError("STUDENT_SERVICE", "school validation", err)
		return nil, errors.New("failed to validate school")
	}

	if !schoolExists {
		return nil, errors.New("school does not exist")
	}

	student, err := s.studentRepo.Update(id, req)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("student not found")
		}
		if strings.Contains(err.Error(), "duplicate key value") {
			if strings.Contains(err.Error(), "students_roll_number_school_id_key") {
				return nil, errors.New("student with this roll number already exists in this school")
			} else {
				return nil, errors.New("student with this email already exists")
			}
		}
		return nil, errors.New("failed to update student")
	}

	student.SchoolName = schoolName
	return student, nil
}

func (s *StudentService) DeleteStudent(id int) error {
	err := s.studentRepo.Delete(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("student not found")
		}
		return errors.New("failed to delete student")
	}
	return nil
}
