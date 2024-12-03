package defined

import api "mt/api/v1"

var (
	/* 系统相关 */
	ErrorUnknownError           = api.ErrorUnknownError("未知错误")
	ErrorServerError            = api.ErrorServerError("服务异常")
	ErrorDataValidateError      = api.ErrorDataValidateError("数据校验失败")
	ErrorDataSelectError        = api.ErrorDataSelectError("数据查询失败")
	ErrorDataAlreadyExists      = api.ErrorDataAlreadyExists("数据已存在")
	ErrorDataNotFound           = api.ErrorDataNotFound("数据不存在")
	ErrorDataAddError           = api.ErrorDataAddError("新增数据失败")
	ErrorDataUpdateError        = api.ErrorDataUpdateError("更新数据失败")
	ErrorDataDeleteError        = api.ErrorDataDeleteError("数据删除失败")
	ErrorDataResourceNotFound   = api.ErrorDataResourceNotFound("数据资源不存在")
	ErrorDataUpdateFieldError   = api.ErrorDataUpdateFieldError("数据属性更新失败")
	ErrorIdInvalidValueError    = api.ErrorIdInvalidValueError("无效ID值")
	ErrorCommandInvalidNotFound = api.ErrorCommandInvalidNotFound("无效的执行指令")
	ErrorNotLoginError          = api.ErrorNotLoginError("请先登录后再操作")
	ErrorNotVisitAuth           = api.ErrorNotVisitAuth("没有访问权限")
	ErrorRequestParamsError     = api.ErrorRequestParamsError("请求参数错误")

	ErrorWebsocketUpgraderError = api.ErrorWebsocketUpgraderError("WebSocket 协议升级失败")

	ErrorGenerateTokenError = api.ErrorGenerateTokenError("生成TOKEN失败")
	ErrorAccountNotFound    = api.ErrorAccountNotFound("账号不存在")
	ErrorAccountLoginError  = api.ErrorAccountLoginError("账号登录错误")

	ErrorToAccountNotFound           = api.ErrorToAccountNotFound("接收者账号不存在")
	ErrorSendMessageTypeNotFound     = api.ErrorSendMessageTypeNotFound("发送消息类型错误")
	IsSendMessageContentRequired     = api.ErrorSendMessageContentRequired("发送消息内容必填")
	ErrorToAccountAndFromAccountSame = api.ErrorToAccountAndFromAccountSame("接收者账号不能和发送者账号一致")
	ErrorSendMessageError            = api.ErrorSendMessageError("发送消息失败, 请重试")
)
