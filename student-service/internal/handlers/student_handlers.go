package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"skool-management/student-service/internal/models"
	"skool-management/student-service/internal/service"
	"skool-management/shared"
)

type StudentHandlers struct {
	studentService *service.StudentService
}

func NewStudentHandlers(studentService *service.StudentService) *StudentHandlers {
	return &StudentHandlers{
		studentService: studentService,
	}
}

func (h *StudentHandlers) CreateStudent(w http.ResponseWriter, r *http.Request) {
	var req models.CreateStudentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.WriteErrorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	student, err := h.studentService.CreateStudent(&req)
	if err != nil {
		switch err.Error() {
		case "roll number, first name, last name, and school ID are required":
			shared.WriteErrorResponse(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		case "school does not exist":
			shared.WriteErrorResponse(w, http.StatusBadRequest, "INVALID_SCHOOL", err.Error())
		case "student with this roll number already exists in this school":
			shared.WriteErrorResponse(w, http.StatusConflict, "ROLL_NUMBER_EXISTS", err.Error())
		case "student with this email already exists":
			shared.WriteErrorResponse(w, http.StatusConflict, "EMAIL_EXISTS", err.Error())
		case "failed to validate school":
			shared.WriteErrorResponse(w, http.StatusInternalServerError, "VALIDATION_ERROR", err.Error())
		default:
			shared.LogError("STUDENT_SERVICE", "create student", err)
			shared.WriteErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create student")
		}
		return
	}

	shared.WriteSuccessResponse(w, http.StatusCreated, "Student created successfully", student)
}

func (h *StudentHandlers) GetStudents(w http.ResponseWriter, r *http.Request) {
	students, err := h.studentService.GetAllStudents()
	if err != nil {
		shared.LogError("STUDENT_SERVICE", "get students", err)
		shared.WriteErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get students")
		return
	}

	shared.WriteSuccessResponse(w, http.StatusOK, "Students retrieved successfully", students)
}

func (h *StudentHandlers) GetStudent(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/students/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		shared.WriteErrorResponse(w, http.StatusBadRequest, "INVALID_ID", "Invalid student ID")
		return
	}

	student, err := h.studentService.GetStudentByID(id)
	if err != nil {
		if err.Error() == "student not found" {
			shared.WriteErrorResponse(w, http.StatusNotFound, "STUDENT_NOT_FOUND", err.Error())
			return
		}
		shared.LogError("STUDENT_SERVICE", "get student", err)
		shared.WriteErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get student")
		return
	}

	shared.WriteSuccessResponse(w, http.StatusOK, "Student retrieved successfully", student)
}

func (h *StudentHandlers) GetStudentsBySchool(w http.ResponseWriter, r *http.Request) {
	schoolIDStr := strings.TrimPrefix(r.URL.Path, "/students/school/")
	schoolID, err := strconv.Atoi(schoolIDStr)
	if err != nil {
		shared.WriteErrorResponse(w, http.StatusBadRequest, "INVALID_SCHOOL_ID", "Invalid school ID")
		return
	}

	students, _, err := h.studentService.GetStudentsBySchoolID(schoolID)
	if err != nil {
		switch err.Error() {
		case "school not found":
			shared.WriteErrorResponse(w, http.StatusNotFound, "SCHOOL_NOT_FOUND", err.Error())
		case "failed to validate school":
			shared.WriteErrorResponse(w, http.StatusInternalServerError, "VALIDATION_ERROR", err.Error())
		default:
			shared.LogError("STUDENT_SERVICE", "get students by school", err)
			shared.WriteErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get students")
		}
		return
	}

	shared.WriteSuccessResponse(w, http.StatusOK, "Students retrieved successfully", students)
}

func (h *StudentHandlers) UpdateStudent(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/students/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		shared.WriteErrorResponse(w, http.StatusBadRequest, "INVALID_ID", "Invalid student ID")
		return
	}

	var req models.UpdateStudentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.WriteErrorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	student, err := h.studentService.UpdateStudent(id, &req)
	if err != nil {
		switch err.Error() {
		case "roll number, first name, last name, and school ID are required":
			shared.WriteErrorResponse(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		case "student not found":
			shared.WriteErrorResponse(w, http.StatusNotFound, "STUDENT_NOT_FOUND", err.Error())
		case "school does not exist":
			shared.WriteErrorResponse(w, http.StatusBadRequest, "INVALID_SCHOOL", err.Error())
		case "student with this roll number already exists in this school":
			shared.WriteErrorResponse(w, http.StatusConflict, "ROLL_NUMBER_EXISTS", err.Error())
		case "student with this email already exists":
			shared.WriteErrorResponse(w, http.StatusConflict, "EMAIL_EXISTS", err.Error())
		case "failed to validate school":
			shared.WriteErrorResponse(w, http.StatusInternalServerError, "VALIDATION_ERROR", err.Error())
		default:
			shared.LogError("STUDENT_SERVICE", "update student", err)
			shared.WriteErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update student")
		}
		return
	}

	shared.WriteSuccessResponse(w, http.StatusOK, "Student updated successfully", student)
}

func (h *StudentHandlers) DeleteStudent(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/students/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		shared.WriteErrorResponse(w, http.StatusBadRequest, "INVALID_ID", "Invalid student ID")
		return
	}

	err = h.studentService.DeleteStudent(id)
	if err != nil {
		if err.Error() == "student not found" {
			shared.WriteErrorResponse(w, http.StatusNotFound, "STUDENT_NOT_FOUND", err.Error())
			return
		}
		shared.LogError("STUDENT_SERVICE", "delete student", err)
		shared.WriteErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to delete student")
		return
	}

	shared.WriteSuccessResponse(w, http.StatusOK, "Student deleted successfully", nil)
}

func (h *StudentHandlers) Health(w http.ResponseWriter, r *http.Request) {
	shared.WriteSuccessResponse(w, http.StatusOK, "Student service is healthy", nil)
}
