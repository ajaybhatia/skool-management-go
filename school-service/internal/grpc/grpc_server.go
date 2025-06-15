package grpc

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	"skool-management/school-service/internal/service"
)

// Protobuf message types
type ProtoSchool struct {
	Id                 string
	RegistrationNumber string
	Name               string
	Address            string
	Phone              string
	Email              string
	CreatedAt          string
	UpdatedAt          string
}

type GetSchoolRequest struct {
	Id string
}

type GetSchoolResponse struct {
	School *ProtoSchool
	Found  bool
}

type ValidateSchoolRequest struct {
	Id string
}

type ValidateSchoolResponse struct {
	Exists bool
	Name   string
}

// GRPCSchoolServer implements the gRPC server
type GRPCSchoolServer struct {
	schoolService *service.SchoolService
}

func NewGRPCSchoolServer(schoolService *service.SchoolService) *GRPCSchoolServer {
	return &GRPCSchoolServer{
		schoolService: schoolService,
	}
}

func (g *GRPCSchoolServer) GetSchool(ctx context.Context, req *GetSchoolRequest) (*GetSchoolResponse, error) {
	id, err := strconv.Atoi(req.Id)
	if err != nil {
		return &GetSchoolResponse{Found: false}, nil
	}

	school, err := g.schoolService.GetSchoolByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return &GetSchoolResponse{Found: false}, nil
		}
		return nil, err
	}

	protoSchool := &ProtoSchool{
		Id:                 strconv.Itoa(school.ID),
		RegistrationNumber: school.RegistrationNumber,
		Name:               school.Name,
		Address:            school.Address,
		Phone:              school.Phone,
		Email:              school.Email,
		CreatedAt:          school.CreatedAt.Format(time.RFC3339),
		UpdatedAt:          school.UpdatedAt.Format(time.RFC3339),
	}

	return &GetSchoolResponse{
		School: protoSchool,
		Found:  true,
	}, nil
}

func (g *GRPCSchoolServer) ValidateSchool(ctx context.Context, req *ValidateSchoolRequest) (*ValidateSchoolResponse, error) {
	id, err := strconv.Atoi(req.Id)
	if err != nil {
		return &ValidateSchoolResponse{Exists: false}, nil
	}

	school, err := g.schoolService.GetSchoolByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return &ValidateSchoolResponse{Exists: false}, nil
		}
		return &ValidateSchoolResponse{Exists: false}, err
	}

	return &ValidateSchoolResponse{
		Exists: true,
		Name:   school.Name,
	}, nil
}
