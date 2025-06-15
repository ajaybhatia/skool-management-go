// Code generated for gRPC communication with School Service
package main

import (
	"context"

	"google.golang.org/grpc"
)

// Protobuf message types
type ValidateSchoolRequest struct {
	Id string
}

type ValidateSchoolResponse struct {
	Exists bool
	Name   string
}

type GetSchoolRequest struct {
	Id string
}

type GetSchoolResponse struct {
	School *ProtoSchool
	Found  bool
}

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

// School Service gRPC client interface
type SchoolServiceClient interface {
	GetSchool(ctx context.Context, in *GetSchoolRequest, opts ...grpc.CallOption) (*GetSchoolResponse, error)
	ValidateSchool(ctx context.Context, in *ValidateSchoolRequest, opts ...grpc.CallOption) (*ValidateSchoolResponse, error)
}

type schoolServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewSchoolServiceClient(cc grpc.ClientConnInterface) SchoolServiceClient {
	return &schoolServiceClient{cc}
}

func (c *schoolServiceClient) GetSchool(ctx context.Context, in *GetSchoolRequest, opts ...grpc.CallOption) (*GetSchoolResponse, error) {
	out := new(GetSchoolResponse)
	// In a real implementation, this would make the actual gRPC call
	// For this example, we'll return a default response
	return out, nil
}

func (c *schoolServiceClient) ValidateSchool(ctx context.Context, in *ValidateSchoolRequest, opts ...grpc.CallOption) (*ValidateSchoolResponse, error) {
	out := new(ValidateSchoolResponse)
	// In a real implementation, this would make the actual gRPC call
	// For this example, we'll return a default response indicating the school exists
	out.Exists = true
	out.Name = "Default School"
	return out, nil
}
