// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.8.2
// - protoc             v5.29.1
// source: v1/account/service.proto

package account

import (
	context "context"
	http "github.com/go-kratos/kratos/v2/transport/http"
	binding "github.com/go-kratos/kratos/v2/transport/http/binding"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
var _ = new(context.Context)
var _ = binding.EncodeURL

const _ = http.SupportPackageIsVersion1

const OperationServiceCreate = "/v1.account.Service/Create"
const OperationServiceDelete = "/v1.account.Service/Delete"
const OperationServiceGenerateToken = "/v1.account.Service/GenerateToken"
const OperationServiceGetInfo = "/v1.account.Service/GetInfo"
const OperationServiceLogin = "/v1.account.Service/Login"
const OperationServiceUpdate = "/v1.account.Service/Update"

type ServiceHTTPServer interface {
	// Create 创建账号
	Create(context.Context, *CreateRequest) (*CreateResponse, error)
	// Delete 删除账号
	Delete(context.Context, *DeleteRequest) (*emptypb.Empty, error)
	// GenerateToken 生成TOKEN
	GenerateToken(context.Context, *GenerateTokenRequest) (*GenerateTokenResponse, error)
	// GetInfo 获取账号信息
	GetInfo(context.Context, *GetInfoRequest) (*GetInfoResponse, error)
	// Login 登录帐号
	Login(context.Context, *LoginRequest) (*LoginResponse, error)
	// Update 更新账号
	Update(context.Context, *UpdateRequest) (*UpdateResponse, error)
}

func RegisterServiceHTTPServer(s *http.Server, srv ServiceHTTPServer) {
	r := s.Route("/")
	r.POST("/api/account/create", _Service_Create0_HTTP_Handler(srv))
	r.PUT("/api/account/update/{account_id}", _Service_Update0_HTTP_Handler(srv))
	r.DELETE("/api/account/delete/{account_id}", _Service_Delete0_HTTP_Handler(srv))
	r.GET("/api/account/info/{account_id}", _Service_GetInfo0_HTTP_Handler(srv))
	r.PUT("/api/account/login/{account_id}", _Service_Login0_HTTP_Handler(srv))
	r.POST("/api/account/token/{account_id}", _Service_GenerateToken0_HTTP_Handler(srv))
}

func _Service_Create0_HTTP_Handler(srv ServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in CreateRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationServiceCreate)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.Create(ctx, req.(*CreateRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*CreateResponse)
		return ctx.Result(200, reply)
	}
}

func _Service_Update0_HTTP_Handler(srv ServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in UpdateRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationServiceUpdate)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.Update(ctx, req.(*UpdateRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*UpdateResponse)
		return ctx.Result(200, reply)
	}
}

func _Service_Delete0_HTTP_Handler(srv ServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in DeleteRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationServiceDelete)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.Delete(ctx, req.(*DeleteRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*emptypb.Empty)
		return ctx.Result(200, reply)
	}
}

func _Service_GetInfo0_HTTP_Handler(srv ServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in GetInfoRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationServiceGetInfo)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetInfo(ctx, req.(*GetInfoRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*GetInfoResponse)
		return ctx.Result(200, reply)
	}
}

func _Service_Login0_HTTP_Handler(srv ServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in LoginRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationServiceLogin)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.Login(ctx, req.(*LoginRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*LoginResponse)
		return ctx.Result(200, reply)
	}
}

func _Service_GenerateToken0_HTTP_Handler(srv ServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in GenerateTokenRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationServiceGenerateToken)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GenerateToken(ctx, req.(*GenerateTokenRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*GenerateTokenResponse)
		return ctx.Result(200, reply)
	}
}

type ServiceHTTPClient interface {
	Create(ctx context.Context, req *CreateRequest, opts ...http.CallOption) (rsp *CreateResponse, err error)
	Delete(ctx context.Context, req *DeleteRequest, opts ...http.CallOption) (rsp *emptypb.Empty, err error)
	GenerateToken(ctx context.Context, req *GenerateTokenRequest, opts ...http.CallOption) (rsp *GenerateTokenResponse, err error)
	GetInfo(ctx context.Context, req *GetInfoRequest, opts ...http.CallOption) (rsp *GetInfoResponse, err error)
	Login(ctx context.Context, req *LoginRequest, opts ...http.CallOption) (rsp *LoginResponse, err error)
	Update(ctx context.Context, req *UpdateRequest, opts ...http.CallOption) (rsp *UpdateResponse, err error)
}

type ServiceHTTPClientImpl struct {
	cc *http.Client
}

func NewServiceHTTPClient(client *http.Client) ServiceHTTPClient {
	return &ServiceHTTPClientImpl{client}
}

func (c *ServiceHTTPClientImpl) Create(ctx context.Context, in *CreateRequest, opts ...http.CallOption) (*CreateResponse, error) {
	var out CreateResponse
	pattern := "/api/account/create"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationServiceCreate))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *ServiceHTTPClientImpl) Delete(ctx context.Context, in *DeleteRequest, opts ...http.CallOption) (*emptypb.Empty, error) {
	var out emptypb.Empty
	pattern := "/api/account/delete/{account_id}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationServiceDelete))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "DELETE", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *ServiceHTTPClientImpl) GenerateToken(ctx context.Context, in *GenerateTokenRequest, opts ...http.CallOption) (*GenerateTokenResponse, error) {
	var out GenerateTokenResponse
	pattern := "/api/account/token/{account_id}"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationServiceGenerateToken))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *ServiceHTTPClientImpl) GetInfo(ctx context.Context, in *GetInfoRequest, opts ...http.CallOption) (*GetInfoResponse, error) {
	var out GetInfoResponse
	pattern := "/api/account/info/{account_id}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationServiceGetInfo))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *ServiceHTTPClientImpl) Login(ctx context.Context, in *LoginRequest, opts ...http.CallOption) (*LoginResponse, error) {
	var out LoginResponse
	pattern := "/api/account/login/{account_id}"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationServiceLogin))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "PUT", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *ServiceHTTPClientImpl) Update(ctx context.Context, in *UpdateRequest, opts ...http.CallOption) (*UpdateResponse, error) {
	var out UpdateResponse
	pattern := "/api/account/update/{account_id}"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationServiceUpdate))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "PUT", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
