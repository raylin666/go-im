// Code generated by protoc-gen-go-errors. DO NOT EDIT.

package v1

import (
	fmt "fmt"
	errors "github.com/go-kratos/kratos/v2/errors"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
const _ = errors.SupportPackageIsVersion1

// 未知错误
func IsUnknownError(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == ErrorReason_UNKNOWN_ERROR.String() && e.Code == 500
}

// 未知错误
func ErrorUnknownError(format string, args ...interface{}) *errors.Error {
	return errors.New(500, ErrorReason_UNKNOWN_ERROR.String(), fmt.Sprintf(format, args...))
}

// 服务异常
func IsServerError(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == ErrorReason_SERVER_ERROR.String() && e.Code == 500
}

// 服务异常
func ErrorServerError(format string, args ...interface{}) *errors.Error {
	return errors.New(500, ErrorReason_SERVER_ERROR.String(), fmt.Sprintf(format, args...))
}

// 数据校验失败
func IsDataValidateError(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == ErrorReason_DATA_VALIDATE_ERROR.String() && e.Code == 422
}

// 数据校验失败
func ErrorDataValidateError(format string, args ...interface{}) *errors.Error {
	return errors.New(422, ErrorReason_DATA_VALIDATE_ERROR.String(), fmt.Sprintf(format, args...))
}

// 数据查询失败
func IsDataSelectError(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == ErrorReason_DATA_SELECT_ERROR.String() && e.Code == 400
}

// 数据查询失败
func ErrorDataSelectError(format string, args ...interface{}) *errors.Error {
	return errors.New(400, ErrorReason_DATA_SELECT_ERROR.String(), fmt.Sprintf(format, args...))
}

// 数据已存在
func IsDataAlreadyExists(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == ErrorReason_DATA_ALREADY_EXISTS.String() && e.Code == 400
}

// 数据已存在
func ErrorDataAlreadyExists(format string, args ...interface{}) *errors.Error {
	return errors.New(400, ErrorReason_DATA_ALREADY_EXISTS.String(), fmt.Sprintf(format, args...))
}

// 数据不存在
func IsDataNotFound(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == ErrorReason_DATA_NOT_FOUND.String() && e.Code == 400
}

// 数据不存在
func ErrorDataNotFound(format string, args ...interface{}) *errors.Error {
	return errors.New(400, ErrorReason_DATA_NOT_FOUND.String(), fmt.Sprintf(format, args...))
}

// 新增数据失败
func IsDataAddError(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == ErrorReason_DATA_ADD_ERROR.String() && e.Code == 400
}

// 新增数据失败
func ErrorDataAddError(format string, args ...interface{}) *errors.Error {
	return errors.New(400, ErrorReason_DATA_ADD_ERROR.String(), fmt.Sprintf(format, args...))
}

// 更新数据失败
func IsDataUpdateError(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == ErrorReason_DATA_UPDATE_ERROR.String() && e.Code == 400
}

// 更新数据失败
func ErrorDataUpdateError(format string, args ...interface{}) *errors.Error {
	return errors.New(400, ErrorReason_DATA_UPDATE_ERROR.String(), fmt.Sprintf(format, args...))
}

// 数据删除失败
func IsDataDeleteError(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == ErrorReason_DATA_DELETE_ERROR.String() && e.Code == 400
}

// 数据删除失败
func ErrorDataDeleteError(format string, args ...interface{}) *errors.Error {
	return errors.New(400, ErrorReason_DATA_DELETE_ERROR.String(), fmt.Sprintf(format, args...))
}

// 数据资源不存在
func IsDataResourceNotFound(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == ErrorReason_DATA_RESOURCE_NOT_FOUND.String() && e.Code == 400
}

// 数据资源不存在
func ErrorDataResourceNotFound(format string, args ...interface{}) *errors.Error {
	return errors.New(400, ErrorReason_DATA_RESOURCE_NOT_FOUND.String(), fmt.Sprintf(format, args...))
}

// 数据属性更新失败
func IsDataUpdateFieldError(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == ErrorReason_DATA_UPDATE_FIELD_ERROR.String() && e.Code == 400
}

// 数据属性更新失败
func ErrorDataUpdateFieldError(format string, args ...interface{}) *errors.Error {
	return errors.New(400, ErrorReason_DATA_UPDATE_FIELD_ERROR.String(), fmt.Sprintf(format, args...))
}

// 无效ID值
func IsIdInvalidValueError(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == ErrorReason_ID_INVALID_VALUE_ERROR.String() && e.Code == 400
}

// 无效ID值
func ErrorIdInvalidValueError(format string, args ...interface{}) *errors.Error {
	return errors.New(400, ErrorReason_ID_INVALID_VALUE_ERROR.String(), fmt.Sprintf(format, args...))
}

// 无效的执行指令
func IsCommandInvalidNotFound(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == ErrorReason_COMMAND_INVALID_NOT_FOUND.String() && e.Code == 400
}

// 无效的执行指令
func ErrorCommandInvalidNotFound(format string, args ...interface{}) *errors.Error {
	return errors.New(400, ErrorReason_COMMAND_INVALID_NOT_FOUND.String(), fmt.Sprintf(format, args...))
}

// 请先登录后再操作
func IsNotLoginError(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == ErrorReason_NOT_LOGIN_ERROR.String() && e.Code == 401
}

// 请先登录后再操作
func ErrorNotLoginError(format string, args ...interface{}) *errors.Error {
	return errors.New(401, ErrorReason_NOT_LOGIN_ERROR.String(), fmt.Sprintf(format, args...))
}

// 没有访问权限
func IsNotVisitAuth(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == ErrorReason_NOT_VISIT_AUTH.String() && e.Code == 401
}

// 没有访问权限
func ErrorNotVisitAuth(format string, args ...interface{}) *errors.Error {
	return errors.New(401, ErrorReason_NOT_VISIT_AUTH.String(), fmt.Sprintf(format, args...))
}

// WebSocket 协议升级失败
func IsWebsocketUpgraderError(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == ErrorReason_WEBSOCKET_UPGRADER_ERROR.String() && e.Code == 400
}

// WebSocket 协议升级失败
func ErrorWebsocketUpgraderError(format string, args ...interface{}) *errors.Error {
	return errors.New(400, ErrorReason_WEBSOCKET_UPGRADER_ERROR.String(), fmt.Sprintf(format, args...))
}

// 生成TOKEN失败
func IsGenerateTokenError(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == ErrorReason_GENERATE_TOKEN_ERROR.String() && e.Code == 500
}

// 生成TOKEN失败
func ErrorGenerateTokenError(format string, args ...interface{}) *errors.Error {
	return errors.New(500, ErrorReason_GENERATE_TOKEN_ERROR.String(), fmt.Sprintf(format, args...))
}

// 账号不存在
func IsAccountNotFound(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == ErrorReason_ACCOUNT_NOT_FOUND.String() && e.Code == 400
}

// 账号不存在
func ErrorAccountNotFound(format string, args ...interface{}) *errors.Error {
	return errors.New(400, ErrorReason_ACCOUNT_NOT_FOUND.String(), fmt.Sprintf(format, args...))
}

// 账号登录错误
func IsAccountLoginError(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == ErrorReason_ACCOUNT_LOGIN_ERROR.String() && e.Code == 400
}

// 账号登录错误
func ErrorAccountLoginError(format string, args ...interface{}) *errors.Error {
	return errors.New(400, ErrorReason_ACCOUNT_LOGIN_ERROR.String(), fmt.Sprintf(format, args...))
}

// 接收者账号不存在
func IsToAccountNotFound(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == ErrorReason_TO_ACCOUNT_NOT_FOUND.String() && e.Code == 400
}

// 接收者账号不存在
func ErrorToAccountNotFound(format string, args ...interface{}) *errors.Error {
	return errors.New(400, ErrorReason_TO_ACCOUNT_NOT_FOUND.String(), fmt.Sprintf(format, args...))
}

// 发送消息类型错误
func IsSendMessageTypeNotFound(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == ErrorReason_SEND_MESSAGE_TYPE_NOT_FOUND.String() && e.Code == 400
}

// 发送消息类型错误
func ErrorSendMessageTypeNotFound(format string, args ...interface{}) *errors.Error {
	return errors.New(400, ErrorReason_SEND_MESSAGE_TYPE_NOT_FOUND.String(), fmt.Sprintf(format, args...))
}

// 发送消息内容必填
func IsSendMessageContentRequired(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == ErrorReason_SEND_MESSAGE_CONTENT_REQUIRED.String() && e.Code == 400
}

// 发送消息内容必填
func ErrorSendMessageContentRequired(format string, args ...interface{}) *errors.Error {
	return errors.New(400, ErrorReason_SEND_MESSAGE_CONTENT_REQUIRED.String(), fmt.Sprintf(format, args...))
}

// 接收者账号不能和发送者账号一致
func IsToAccountAndFromAccountSame(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == ErrorReason_TO_ACCOUNT_AND_FROM_ACCOUNT_SAME.String() && e.Code == 400
}

// 接收者账号不能和发送者账号一致
func ErrorToAccountAndFromAccountSame(format string, args ...interface{}) *errors.Error {
	return errors.New(400, ErrorReason_TO_ACCOUNT_AND_FROM_ACCOUNT_SAME.String(), fmt.Sprintf(format, args...))
}

// 发送消息失败, 请重试
func IsSendMessageError(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == ErrorReason_SEND_MESSAGE_ERROR.String() && e.Code == 400
}

// 发送消息失败, 请重试
func ErrorSendMessageError(format string, args ...interface{}) *errors.Error {
	return errors.New(400, ErrorReason_SEND_MESSAGE_ERROR.String(), fmt.Sprintf(format, args...))
}
