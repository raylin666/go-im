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

func errFormat(format string, errMsg string) string {
	if len(errMsg) > 0 {
		format = fmt.Sprintf("%s: %s", format, errMsg)
	}

	return format
}

func (err *Error) Unknown(args ...interface{}) *errors.Error {
	return pb.ErrorUnknown(errFormat("未知错误", err.format), args...)
}

func (err *Error) Server(args ...interface{}) *errors.Error {
	return pb.ErrorServer(errFormat("服务异常", err.format), args...)
}

func (err *Error) ServerUpgrader(args ...interface{}) *errors.Error {
	return pb.ErrorServerUpgrader(errFormat("服务协议升级失败", err.format), args...)
}

func (err *Error) DataValidate(args ...interface{}) *errors.Error {
	return pb.ErrorDataValidate(errFormat("数据校验失败", err.format), args...)
}

func (err *Error) DataSelect(args ...interface{}) *errors.Error {
	return pb.ErrorDataSelect(errFormat("数据查询失败", err.format), args...)
}

func (err *Error) DataAlreadyExists(args ...interface{}) *errors.Error {
	return pb.ErrorDataAlreadyExists(errFormat("数据已存在", err.format), args...)
}

func (err *Error) DataNotFound(args ...interface{}) *errors.Error {
	return pb.ErrorDataNotFound(errFormat("数据不存在", err.format), args...)
}

func (err *Error) DataAdd(args ...interface{}) *errors.Error {
	return pb.ErrorDataAdd(errFormat("新增数据失败", err.format), args...)
}

func (err *Error) DataUpdate(args ...interface{}) *errors.Error {
	return pb.ErrorDataUpdate(errFormat("更新数据失败", err.format), args...)
}

func (err *Error) DataDelete(args ...interface{}) *errors.Error {
	return pb.ErrorDataDelete(errFormat("数据删除失败", err.format), args...)
}

func (err *Error) DataResourceNotFound(args ...interface{}) *errors.Error {
	return pb.ErrorDataResourceNotFound(errFormat("数据资源不存在", err.format), args...)
}

func (err *Error) DataUpdateField(args ...interface{}) *errors.Error {
	return pb.ErrorDataUpdateField(errFormat("数据属性更新失败", err.format), args...)
}

func (err *Error) IdInvalidValue(args ...interface{}) *errors.Error {
	return pb.ErrorIdInvalidValue(errFormat("无效ID值", err.format), args...)
}

func (err *Error) CommandInvalidNotFound(args ...interface{}) *errors.Error {
	return pb.ErrorCommandInvalidNotFound(errFormat("无效的执行指令", err.format), args...)
}

func (err *Error) RequestParams(args ...interface{}) *errors.Error {
	return pb.ErrorRequestParams(errFormat("请求参数错误", err.format), args...)
}

func (err *Error) NotLogin(args ...interface{}) *errors.Error {
	return pb.ErrorNotLogin(errFormat("未登录帐号", err.format), args...)
}

func (err *Error) NotVisitAuth(args ...interface{}) *errors.Error {
	return pb.ErrorNotVisitAuth(errFormat("没有访问权限", err.format), args...)
}

func (err *Error) GenerateToken(args ...interface{}) *errors.Error {
	return pb.ErrorGenerateToken(errFormat("生成TOKEN失败", err.format), args...)
}

func (err *Error) AccountNotFound(args ...interface{}) *errors.Error {
	return pb.ErrorAccountNotFound(errFormat("账号不存在", err.format), args...)
}

func (err *Error) AccountLogin(args ...interface{}) *errors.Error {
	return pb.ErrorAccountLogin(errFormat("账号登录错误", err.format), args...)
}

func (err *Error) AccountIsLogin(args ...interface{}) *errors.Error {
	return pb.ErrorAccountIsLogin(errFormat("帐号已登录", err.format), args...)
}

func (err *Error) ToAccountNotFound(args ...interface{}) *errors.Error {
	return pb.ErrorToAccountNotFound(errFormat("接收者账号不存在", err.format), args...)
}

func (err *Error) SendMessageTypeNotFound(args ...interface{}) *errors.Error {
	return pb.ErrorSendMessageTypeNotFound(errFormat("发送消息类型错误", err.format), args...)
}

func (err *Error) SendMessageContentRequired(args ...interface{}) *errors.Error {
	return pb.ErrorSendMessageContentRequired(errFormat("发送消息内容必填", err.format), args...)
}

func (err *Error) ToAccountAndFromAccountSame(args ...interface{}) *errors.Error {
	return pb.ErrorToAccountAndFromAccountSame(errFormat("接收者账号不能和发送者账号一致", err.format), args...)
}

func (err *Error) SendMessage(args ...interface{}) *errors.Error {
	return pb.ErrorSendMessage(errFormat("发送消息失败, 请重试", err.format), args...)
}
