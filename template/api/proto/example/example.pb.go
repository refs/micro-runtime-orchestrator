// Code generated by protoc-gen-go.
// source: github.com/micro/examples/template/api/proto/example/example.proto
// DO NOT EDIT!

/*
Package go_micro_api_template is a generated protocol buffer package.

It is generated from these files:
	github.com/micro/examples/template/api/proto/example/example.proto

It has these top-level messages:
*/
package go_micro_api_template

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import go_micro_api "github.com/micro/micro/api/proto"

import (
	client "github.com/micro/go-micro/client"
	server "github.com/micro/go-micro/server"
	context "golang.org/x/net/context"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for Example service

type ExampleClient interface {
	Call(ctx context.Context, in *go_micro_api.Request, opts ...client.CallOption) (*go_micro_api.Response, error)
}

type exampleClient struct {
	c           client.Client
	serviceName string
}

func NewExampleClient(serviceName string, c client.Client) ExampleClient {
	if c == nil {
		c = client.NewClient()
	}
	if len(serviceName) == 0 {
		serviceName = "go.micro.api.template"
	}
	return &exampleClient{
		c:           c,
		serviceName: serviceName,
	}
}

func (c *exampleClient) Call(ctx context.Context, in *go_micro_api.Request, opts ...client.CallOption) (*go_micro_api.Response, error) {
	req := c.c.NewRequest(c.serviceName, "Example.Call", in)
	out := new(go_micro_api.Response)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Example service

type ExampleHandler interface {
	Call(context.Context, *go_micro_api.Request, *go_micro_api.Response) error
}

func RegisterExampleHandler(s server.Server, hdlr ExampleHandler, opts ...server.HandlerOption) {
	s.Handle(s.NewHandler(&Example{hdlr}, opts...))
}

type Example struct {
	ExampleHandler
}

func (h *Example) Call(ctx context.Context, in *go_micro_api.Request, out *go_micro_api.Response) error {
	return h.ExampleHandler.Call(ctx, in, out)
}

var fileDescriptor0 = []byte{
	// 142 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0xf2, 0x48, 0xcf, 0x2c, 0xc9,
	0x28, 0x4d, 0xd2, 0x4b, 0xce, 0xcf, 0xd5, 0xcf, 0xcd, 0x4c, 0x2e, 0xca, 0x87, 0x92, 0xa9, 0x15,
	0x89, 0xb9, 0x05, 0x39, 0xa9, 0xc5, 0xfa, 0x25, 0xa9, 0x40, 0x3a, 0xb1, 0x24, 0x55, 0x3f, 0xb1,
	0x20, 0x53, 0xbf, 0xa0, 0x28, 0xbf, 0x04, 0x2e, 0x07, 0xa3, 0xf5, 0xc0, 0xa2, 0x42, 0xa2, 0xe9,
	0xf9, 0x7a, 0x60, 0xbd, 0x7a, 0x40, 0x95, 0x7a, 0x30, 0x6d, 0x52, 0x5a, 0x38, 0x2c, 0x40, 0x18,
	0x07, 0x52, 0x0e, 0x66, 0x19, 0x39, 0x71, 0xb1, 0xbb, 0x42, 0xcc, 0x14, 0x32, 0xe7, 0x62, 0x71,
	0x4e, 0xcc, 0xc9, 0x11, 0x12, 0xd5, 0x43, 0x31, 0x36, 0x28, 0xb5, 0xb0, 0x34, 0xb5, 0xb8, 0x44,
	0x4a, 0x0c, 0x5d, 0xb8, 0xb8, 0x20, 0x3f, 0xaf, 0x38, 0x55, 0x89, 0x21, 0x89, 0x0d, 0x6c, 0x94,
	0x31, 0x20, 0x00, 0x00, 0xff, 0xff, 0xf9, 0xcb, 0xce, 0xa3, 0xd9, 0x00, 0x00, 0x00,
}
