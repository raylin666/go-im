package errors

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	pb "mt/api/v1"
	"regexp"
	"strconv"
	"strings"
)

// ErrorMessage 获取错误信息
func ErrorMessage(err *errors.Error) (code uint32, message string) {
	code = uint32(err.Code)
	message = err.Message
	return
}

type ErrorMessageDetails struct {
	Code     int
	Reason   string
	Message  string
	Metadata map[string]string
	Cause    string
}

// ParseErrorMessage 解析错误信息字符串, 例如前置将 errors.Error 转换成了原生的 error, 此时需要再转换回 errors.Error
func ParseErrorMessage(err error) (*ErrorMessageDetails, error) {
	re := regexp.MustCompile(`error: code = (\d+) reason = (\w+) message = (.*?) metadata = (\w+\[.*?\]) cause = (<nil>|\w+)`)
	matches := re.FindStringSubmatch(err.Error())
	if len(matches) < 6 {
		return nil, fmt.Errorf("invalid error string format")
	}

	// 提取 Code
	code, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil, fmt.Errorf("failed to parse code: %v", err)
	}

	// 提取 Reason
	reason := matches[2]

	// 提取 Message
	message := strings.TrimSpace(matches[3])

	// 提取 Metadata
	metadata := make(map[string]string)
	if matches[4] != "map[]" {
		// 有内容，进一步解析
		// 例如：metadataStr = "map[key1:value1 key2:value2]"
		metadataRe := regexp.MustCompile(`(\w+):(\w+)`)
		metadataMatches := metadataRe.FindAllStringSubmatch(matches[4], -1)
		for _, match := range metadataMatches {
			metadata[match[1]] = match[2]
		}
	}

	// 提取 Cause
	cause := matches[5]

	return &ErrorMessageDetails{
		Code:     code,
		Reason:   reason,
		Message:  message,
		Metadata: metadata,
		Cause:    cause,
	}, nil
}

type Option func(opt *option)

func WithMessage(format string) Option {
	return func(opt *option) {
		opt.format = format
	}
}

type option struct{ format string }

func Is(err, target error) bool { return errors.Is(err, target) }

type Error struct{ *option }

func New(opts ...Option) *Error {
	var err = new(Error)
	var o = new(option)
	for _, opt := range opts {
		opt(o)
	}
	err.option = o
	return err
}

func (err *Error) Unknown(args ...interface{}) *errors.Error {
	var format = "未知错误"
	if len(err.format) > 0 {
		err.format = fmt.Sprintf("%s: %s", format, err.format)
	}

	return pb.ErrorUnknown(err.format, args...)
}

func (err *Error) Server(args ...interface{}) *errors.Error {
	var format = "服务异常"
	if len(err.format) > 0 {
		err.format = fmt.Sprintf("%s: %s", format, err.format)
	}

	return pb.ErrorServer(err.format, args...)
}

func (err *Error) ServerUpgrader(args ...interface{}) *errors.Error {
	var format = "服务协议升级失败"
	if len(err.format) > 0 {
		err.format = fmt.Sprintf("%s: %s", format, err.format)
	}

	return pb.ErrorServerUpgrader(err.format, args...)
}

func (err *Error) DataValidate(args ...interface{}) *errors.Error {
	var format = "数据校验失败"
	if len(err.format) > 0 {
		err.format = fmt.Sprintf("%s: %s", format, err.format)
	}

	return pb.ErrorDataValidate(err.format, args...)
}

func (err *Error) DataSelect(args ...interface{}) *errors.Error {
	var format = "数据查询失败"
	if len(err.format) > 0 {
		err.format = fmt.Sprintf("%s: %s", format, err.format)
	}

	return pb.ErrorDataSelect(err.format, args...)
}

func (err *Error) DataAlreadyExists(args ...interface{}) *errors.Error {
	var format = "数据已存在"
	if len(err.format) > 0 {
		err.format = fmt.Sprintf("%s: %s", format, err.format)
	}

	return pb.ErrorDataAlreadyExists(err.format, args...)
}

func (err *Error) DataNotFound(args ...interface{}) *errors.Error {
	var format = "数据不存在"
	if len(err.format) > 0 {
		err.format = fmt.Sprintf("%s: %s", format, err.format)
	}

	return pb.ErrorDataNotFound(err.format, args...)
}

func (err *Error) DataAdd(args ...interface{}) *errors.Error {
	var format = "新增数据失败"
	if len(err.format) > 0 {
		err.format = fmt.Sprintf("%s: %s", format, err.format)
	}

	return pb.ErrorDataAdd(err.format, args...)
}

func (err *Error) DataUpdate(args ...interface{}) *errors.Error {
	var format = "更新数据失败"
	if len(err.format) > 0 {
		err.format = fmt.Sprintf("%s: %s", format, err.format)
	}

	return pb.ErrorDataUpdate(err.format, args...)
}

func (err *Error) DataDelete(args ...interface{}) *errors.Error {
	var format = "数据删除失败"
	if len(err.format) > 0 {
		err.format = fmt.Sprintf("%s: %s", format, err.format)
	}

	return pb.ErrorDataDelete(err.format, args...)
}

func (err *Error) DataResourceNotFound(args ...interface{}) *errors.Error {
	var format = "数据资源不存在"
	if len(err.format) > 0 {
		err.format = fmt.Sprintf("%s: %s", format, err.format)
	}

	return pb.ErrorDataResourceNotFound(err.format, args...)
}

func (err *Error) DataUpdateField(args ...interface{}) *errors.Error {
	var format = "数据属性更新失败"
	if len(err.format) > 0 {
		err.format = fmt.Sprintf("%s: %s", format, err.format)
	}

	return pb.ErrorDataUpdateField(err.format, args...)
}

func (err *Error) IdInvalidValue(args ...interface{}) *errors.Error {
	var format = "无效ID值"
	if len(err.format) > 0 {
		err.format = fmt.Sprintf("%s: %s", format, err.format)
	}

	return pb.ErrorIdInvalidValue(err.format, args...)
}

func (err *Error) CommandInvalidNotFound(args ...interface{}) *errors.Error {
	var format = "无效的执行指令"
	if len(err.format) > 0 {
		err.format = fmt.Sprintf("%s: %s", format, err.format)
	}

	return pb.ErrorCommandInvalidNotFound(err.format, args...)
}

func (err *Error) RequestParams(args ...interface{}) *errors.Error {
	var format = "请求参数错误"
	if len(err.format) > 0 {
		err.format = fmt.Sprintf("%s: %s", format, err.format)
	}

	return pb.ErrorRequestParams(err.format, args...)
}

func (err *Error) NotLogin(args ...interface{}) *errors.Error {
	var format = "未登录帐号"
	if len(err.format) > 0 {
		err.format = fmt.Sprintf("%s: %s", format, err.format)
	}

	return pb.ErrorNotLogin(err.format, args...)
}

func (err *Error) NotVisitAuth(args ...interface{}) *errors.Error {
	var format = "没有访问权限"
	if len(err.format) > 0 {
		err.format = fmt.Sprintf("%s: %s", format, err.format)
	}

	return pb.ErrorNotVisitAuth(err.format, args...)
}

func (err *Error) GenerateToken(args ...interface{}) *errors.Error {
	var format = "生成TOKEN失败"
	if len(err.format) > 0 {
		err.format = fmt.Sprintf("%s: %s", format, err.format)
	}

	return pb.ErrorGenerateToken(err.format, args...)
}

func (err *Error) AccountNotFound(args ...interface{}) *errors.Error {
	var format = "账号不存在"
	if len(err.format) > 0 {
		err.format = fmt.Sprintf("%s: %s", format, err.format)
	}

	return pb.ErrorAccountNotFound(err.format, args...)
}

func (err *Error) AccountLogin(args ...interface{}) *errors.Error {
	var format = "账号登录错误"
	if len(err.format) > 0 {
		err.format = fmt.Sprintf("%s: %s", format, err.format)
	}

	return pb.ErrorAccountLogin(err.format, args...)
}

func (err *Error) AccountIsLogin(args ...interface{}) *errors.Error {
	var format = "帐号已登录"
	if len(err.format) > 0 {
		err.format = fmt.Sprintf("%s: %s", format, err.format)
	}

	return pb.ErrorAccountIsLogin(err.format, args...)
}

func (err *Error) ToAccountNotFound(args ...interface{}) *errors.Error {
	var format = "接收者账号不存在"
	if len(err.format) > 0 {
		err.format = fmt.Sprintf("%s: %s", format, err.format)
	}

	return pb.ErrorToAccountNotFound(err.format, args...)
}

func (err *Error) SendMessageTypeNotFound(args ...interface{}) *errors.Error {
	var format = "发送消息类型错误"
	if len(err.format) > 0 {
		err.format = fmt.Sprintf("%s: %s", format, err.format)
	}

	return pb.ErrorSendMessageTypeNotFound(err.format, args...)
}

func (err *Error) SendMessageContentRequired(args ...interface{}) *errors.Error {
	var format = "发送消息内容必填"
	if len(err.format) > 0 {
		err.format = fmt.Sprintf("%s: %s", format, err.format)
	}

	return pb.ErrorSendMessageContentRequired(err.format, args...)
}

func (err *Error) ToAccountAndFromAccountSame(args ...interface{}) *errors.Error {
	var format = "接收者账号不能和发送者账号一致"
	if len(err.format) > 0 {
		err.format = fmt.Sprintf("%s: %s", format, err.format)
	}

	return pb.ErrorToAccountAndFromAccountSame(err.format, args...)
}

func (err *Error) SendMessage(args ...interface{}) *errors.Error {
	var format = "发送消息失败, 请重试"
	if len(err.format) > 0 {
		err.format = fmt.Sprintf("%s: %s", format, err.format)
	}

	return pb.ErrorSendMessage(err.format, args...)
}
