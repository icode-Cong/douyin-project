// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: favoriteService.proto

package favoriteService

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

import (
	context "context"
	api "github.com/micro/go-micro/v2/api"
	client "github.com/micro/go-micro/v2/client"
	server "github.com/micro/go-micro/v2/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ api.Endpoint
var _ context.Context
var _ client.Option
var _ server.Option

// Api Endpoints for FavoriteService service

func NewFavoriteServiceEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for FavoriteService service

type FavoriteService interface {
	FavoriteAction(ctx context.Context, in *DouyinFavoriteActionRequest, opts ...client.CallOption) (*DouyinFavoriteActionResponse, error)
	FavoriteList(ctx context.Context, in *DouyinFavoriteListRequest, opts ...client.CallOption) (*DouyinFavoriteListResponse, error)
}

type favoriteService struct {
	c    client.Client
	name string
}

func NewFavoriteService(name string, c client.Client) FavoriteService {
	return &favoriteService{
		c:    c,
		name: name,
	}
}

func (c *favoriteService) FavoriteAction(ctx context.Context, in *DouyinFavoriteActionRequest, opts ...client.CallOption) (*DouyinFavoriteActionResponse, error) {
	req := c.c.NewRequest(c.name, "FavoriteService.FavoriteAction", in)
	out := new(DouyinFavoriteActionResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *favoriteService) FavoriteList(ctx context.Context, in *DouyinFavoriteListRequest, opts ...client.CallOption) (*DouyinFavoriteListResponse, error) {
	req := c.c.NewRequest(c.name, "FavoriteService.FavoriteList", in)
	out := new(DouyinFavoriteListResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for FavoriteService service

type FavoriteServiceHandler interface {
	FavoriteAction(context.Context, *DouyinFavoriteActionRequest, *DouyinFavoriteActionResponse) error
	FavoriteList(context.Context, *DouyinFavoriteListRequest, *DouyinFavoriteListResponse) error
}

func RegisterFavoriteServiceHandler(s server.Server, hdlr FavoriteServiceHandler, opts ...server.HandlerOption) error {
	type favoriteService interface {
		FavoriteAction(ctx context.Context, in *DouyinFavoriteActionRequest, out *DouyinFavoriteActionResponse) error
		FavoriteList(ctx context.Context, in *DouyinFavoriteListRequest, out *DouyinFavoriteListResponse) error
	}
	type FavoriteService struct {
		favoriteService
	}
	h := &favoriteServiceHandler{hdlr}
	return s.Handle(s.NewHandler(&FavoriteService{h}, opts...))
}

type favoriteServiceHandler struct {
	FavoriteServiceHandler
}

func (h *favoriteServiceHandler) FavoriteAction(ctx context.Context, in *DouyinFavoriteActionRequest, out *DouyinFavoriteActionResponse) error {
	return h.FavoriteServiceHandler.FavoriteAction(ctx, in, out)
}

func (h *favoriteServiceHandler) FavoriteList(ctx context.Context, in *DouyinFavoriteListRequest, out *DouyinFavoriteListResponse) error {
	return h.FavoriteServiceHandler.FavoriteList(ctx, in, out)
}
