package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"skool-management/school-service/internal/models"
	"skool-management/school-service/internal/service"
	"skool-management/shared"
)

type SchoolHandlers struct {
	schoolService *service.SchoolService
}

func NewSchoolHandlers(schoolService *service.SchoolService) *SchoolHandlers {
	return &SchoolHandlers{
		schoolService: schoolService,
	}
}

func (h *SchoolHandlers) CreateSchool(w http.ResponseWriter, r *http.Request) {
	var req models.CreateSchoolRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.WriteErrorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	school, err := h.schoolService.CreateSchool(&req)
	if err != nil {
		switch err.Error() {
		case "school name is required":
			shared.WriteErrorResponse(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		case "school registration number is required":
			shared.WriteErrorResponse(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		case "school with this registration number already exists":
			shared.WriteErrorResponse(w, http.StatusConflict, "REGISTRATION_EXISTS", err.Error())
		default:
			shared.LogError("SCHOOL_SERVICE", "create school", err)
			shared.WriteErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create school")
		}
		return
	}

	shared.WriteSuccessResponse(w, http.StatusCreated, "School created successfully", school)
}

func (h *SchoolHandlers) GetSchools(w http.ResponseWriter, r *http.Request) {
	schools, err := h.schoolService.GetAllSchools()
	if err != nil {
		shared.LogError("SCHOOL_SERVICE", "get schools", err)
		shared.WriteErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get schools")
		return
	}

	shared.WriteSuccessResponse(w, http.StatusOK, "Schools retrieved successfully", schools)
}

func (h *SchoolHandlers) GetSchool(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/schools/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		shared.WriteErrorResponse(w, http.StatusBadRequest, "INVALID_ID", "Invalid school ID")
		return
	}

	school, err := h.schoolService.GetSchoolByID(id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			shared.WriteErrorResponse(w, http.StatusNotFound, "SCHOOL_NOT_FOUND", "School not found")
			return
		}
		shared.LogError("SCHOOL_SERVICE", "get school", err)
		shared.WriteErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get school")
		return
	}

	shared.WriteSuccessResponse(w, http.StatusOK, "School retrieved successfully", school)
}

func (h *SchoolHandlers) UpdateSchool(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/schools/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		shared.WriteErrorResponse(w, http.StatusBadRequest, "INVALID_ID", "Invalid school ID")
		return
	}

	var req models.UpdateSchoolRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.WriteErrorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	school, err := h.schoolService.UpdateSchool(id, &req)
	if err != nil {
		switch err.Error() {
		case "registration number and name are required":
			shared.WriteErrorResponse(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		case "school not found":
			shared.WriteErrorResponse(w, http.StatusNotFound, "SCHOOL_NOT_FOUND", err.Error())
		case "school with this registration number already exists":
			shared.WriteErrorResponse(w, http.StatusConflict, "REGISTRATION_NUMBER_EXISTS", err.Error())
		case "school with this email already exists":
			shared.WriteErrorResponse(w, http.StatusConflict, "EMAIL_EXISTS", err.Error())
		default:
			shared.LogError("SCHOOL_SERVICE", "update school", err)
			shared.WriteErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update school")
		}
		return
	}

	shared.WriteSuccessResponse(w, http.StatusOK, "School updated successfully", school)
}

func (h *SchoolHandlers) DeleteSchool(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/schools/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		shared.WriteErrorResponse(w, http.StatusBadRequest, "INVALID_ID", "Invalid school ID")
		return
	}

	err = h.schoolService.DeleteSchool(id)
	if err != nil {
		if err.Error() == "school not found" {
			shared.WriteErrorResponse(w, http.StatusNotFound, "SCHOOL_NOT_FOUND", err.Error())
			return
		}
		shared.LogError("SCHOOL_SERVICE", "delete school", err)
		shared.WriteErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to delete school")
		return
	}

	shared.WriteSuccessResponse(w, http.StatusOK, "School deleted successfully", nil)
}

func (h *SchoolHandlers) Health(w http.ResponseWriter, r *http.Request) {
	shared.WriteSuccessResponse(w, http.StatusOK, "School service is healthy", nil)
}
