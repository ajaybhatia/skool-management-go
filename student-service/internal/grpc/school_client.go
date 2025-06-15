package grpc

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

// Protobuf message types for School Service communication
type ValidateSchoolRequest struct {
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (m *ValidateSchoolRequest) Reset()         { *m = ValidateSchoolRequest{} }
func (m *ValidateSchoolRequest) String() string { return fmt.Sprintf("ValidateSchoolRequest{Id: %s}", m.Id) }
func (*ValidateSchoolRequest) ProtoMessage()    {}

type ValidateSchoolResponse struct {
	Exists bool   `protobuf:"varint,1,opt,name=exists,proto3" json:"exists,omitempty"`
	Name   string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
}

func (m *ValidateSchoolResponse) Reset()         { *m = ValidateSchoolResponse{} }
func (m *ValidateSchoolResponse) String() string { return fmt.Sprintf("ValidateSchoolResponse{Exists: %t, Name: %s}", m.Exists, m.Name) }
func (*ValidateSchoolResponse) ProtoMessage()    {}

type GetSchoolRequest struct {
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (m *GetSchoolRequest) Reset()         { *m = GetSchoolRequest{} }
func (m *GetSchoolRequest) String() string { return fmt.Sprintf("GetSchoolRequest{Id: %s}", m.Id) }
func (*GetSchoolRequest) ProtoMessage()    {}

type GetSchoolResponse struct {
	School *ProtoSchool `protobuf:"bytes,1,opt,name=school,proto3" json:"school,omitempty"`
	Found  bool         `protobuf:"varint,2,opt,name=found,proto3" json:"found,omitempty"`
}

func (m *GetSchoolResponse) Reset()         { *m = GetSchoolResponse{} }
func (m *GetSchoolResponse) String() string { return fmt.Sprintf("GetSchoolResponse{Found: %t}", m.Found) }
func (*GetSchoolResponse) ProtoMessage()    {}

type ProtoSchool struct {
	Id                 string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	RegistrationNumber string `protobuf:"bytes,2,opt,name=registration_number,json=registrationNumber,proto3" json:"registration_number,omitempty"`
	Name               string `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	Address            string `protobuf:"bytes,4,opt,name=address,proto3" json:"address,omitempty"`
	Phone              string `protobuf:"bytes,5,opt,name=phone,proto3" json:"phone,omitempty"`
	Email              string `protobuf:"bytes,6,opt,name=email,proto3" json:"email,omitempty"`
	CreatedAt          string `protobuf:"bytes,7,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt          string `protobuf:"bytes,8,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
}

func (m *ProtoSchool) Reset()         { *m = ProtoSchool{} }
func (m *ProtoSchool) String() string { return fmt.Sprintf("ProtoSchool{Name: %s}", m.Name) }
func (*ProtoSchool) ProtoMessage()    {}

// Implement proto.Message interface methods
func (m *ValidateSchoolRequest) XXX_Unmarshal(b []byte) error { return nil }
func (m *ValidateSchoolRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) { return nil, nil }
func (m *ValidateSchoolRequest) XXX_Merge(src proto.Message) {}
func (m *ValidateSchoolRequest) XXX_Size() int { return 0 }
func (m *ValidateSchoolRequest) XXX_DiscardUnknown() {}

func (m *ValidateSchoolResponse) XXX_Unmarshal(b []byte) error { return nil }
func (m *ValidateSchoolResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) { return nil, nil }
func (m *ValidateSchoolResponse) XXX_Merge(src proto.Message) {}
func (m *ValidateSchoolResponse) XXX_Size() int { return 0 }
func (m *ValidateSchoolResponse) XXX_DiscardUnknown() {}

func (m *GetSchoolRequest) XXX_Unmarshal(b []byte) error { return nil }
func (m *GetSchoolRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) { return nil, nil }
func (m *GetSchoolRequest) XXX_Merge(src proto.Message) {}
func (m *GetSchoolRequest) XXX_Size() int { return 0 }
func (m *GetSchoolRequest) XXX_DiscardUnknown() {}

func (m *GetSchoolResponse) XXX_Unmarshal(b []byte) error { return nil }
func (m *GetSchoolResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) { return nil, nil }
func (m *GetSchoolResponse) XXX_Merge(src proto.Message) {}
func (m *GetSchoolResponse) XXX_Size() int { return 0 }
func (m *GetSchoolResponse) XXX_DiscardUnknown() {}

func (m *ProtoSchool) XXX_Unmarshal(b []byte) error { return nil }
func (m *ProtoSchool) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) { return nil, nil }
func (m *ProtoSchool) XXX_Merge(src proto.Message) {}
func (m *ProtoSchool) XXX_Size() int { return 0 }
func (m *ProtoSchool) XXX_DiscardUnknown() {}

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
	err := c.cc.Invoke(ctx, "/SchoolService/GetSchool", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *schoolServiceClient) ValidateSchool(ctx context.Context, in *ValidateSchoolRequest, opts ...grpc.CallOption) (*ValidateSchoolResponse, error) {
	out := new(ValidateSchoolResponse)
	err := c.cc.Invoke(ctx, "/SchoolService/ValidateSchool", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}
