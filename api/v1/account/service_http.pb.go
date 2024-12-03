// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.8.2
// - protoc             v3.15.8
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
const OperationServiceUpdate = "/v1.account.Service/Update"
const OperationServiceUpdateLogin = "/v1.account.Service/UpdateLogin"

type ServiceHTTPServer interface {
	// Create 创建账号
	Create(context.Context, *CreateRequest) (*CreateResponse, error)
	// Delete 删除账号
	Delete(context.Context, *DeleteRequest) (*emptypb.Empty, error)
	// GenerateToken 生成TOKEN
	GenerateToken(context.Context, *GenerateTokenRequest) (*GenerateTokenResponse, error)
	// Update 更新账号
	Update(context.Context, *UpdateRequest) (*UpdateResponse, error)
	// UpdateLogin 更新帐号登录信息
	UpdateLogin(context.Context, *UpdateLoginRequest) (*UpdateLoginResponse, error)
}

func RegisterServiceHTTPServer(s *http.Server, srv ServiceHTTPServer) {
	r := s.Route("/")
	r.POST("/api/account/create", _Service_Create0_HTTP_Handler(srv))
	r.PUT("/api/account/update/{account_id}", _Service_Update0_HTTP_Handler(srv))
	r.DELETE("/api/account/delete/{account_id}", _Service_Delete0_HTTP_Handler(srv))
	r.PUT("/api/account/update_login/{account_id}", _Service_UpdateLogin0_HTTP_Handler(srv))
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

func _Service_UpdateLogin0_HTTP_Handler(srv ServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in UpdateLoginRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationServiceUpdateLogin)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.UpdateLogin(ctx, req.(*UpdateLoginRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*UpdateLoginResponse)
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
	Update(ctx context.Context, req *UpdateRequest, opts ...http.CallOption) (rsp *UpdateResponse, err error)
	UpdateLogin(ctx context.Context, req *UpdateLoginRequest, opts ...http.CallOption) (rsp *UpdateLoginResponse, err error)
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

func (c *ServiceHTTPClientImpl) UpdateLogin(ctx context.Context, in *UpdateLoginRequest, opts ...http.CallOption) (*UpdateLoginResponse, error) {
	var out UpdateLoginResponse
	pattern := "/api/account/update_login/{account_id}"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationServiceUpdateLogin))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "PUT", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
