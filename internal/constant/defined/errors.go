package defined

import api "mt/api/v1"

var (
	/* 系统相关 */
	ErrorUnknown              = api.ErrorUnknown("未知错误")
	ErrorServer               = api.ErrorServer("服务异常")
	ErrorServerUpgrader       = api.ErrorServerUpgrader("服务协议升级失败")
	ErrorDataValidate         = api.ErrorDataValidate("数据校验失败")
	ErrorDataSelect           = api.ErrorDataSelect("数据查询失败")
	ErrorDataAlreadyExists    = api.ErrorDataAlreadyExists("数据已存在")
	ErrorDataNotFound         = api.ErrorDataNotFound("数据不存在")
	ErrorDataAdd              = api.ErrorDataAdd("新增数据失败")
	ErrorDataUpdate           = api.ErrorDataUpdate("更新数据失败")
	ErrorDataDelete           = api.ErrorDataDelete("数据删除失败")
	ErrorDataResourceNotFound = api.ErrorDataResourceNotFound("数据资源不存在")
	ErrorDataUpdateField      = api.ErrorDataUpdateField("数据属性更新失败")

	ErrorIdInvalidValue         = api.ErrorIdInvalidValue("无效ID值")
	ErrorCommandInvalidNotFound = api.ErrorCommandInvalidNotFound("无效的执行指令")
	ErrorRequestParams          = api.ErrorRequestParams("请求参数错误")

	ErrorNotLogin        = api.ErrorNotLogin("未登录帐号")
	ErrorNotVisitAuth    = api.ErrorNotVisitAuth("没有访问权限")
	ErrorGenerateToken   = api.ErrorGenerateToken("生成TOKEN失败")
	ErrorAccountNotFound = api.ErrorAccountNotFound("账号不存在")
	ErrorAccountLogin    = api.ErrorAccountLogin("账号登录错误")
	ErrorAccountIsLogin  = api.ErrorAccountIsLogin("帐号已登录")

	ErrorToAccountNotFound           = api.ErrorToAccountNotFound("接收者账号不存在")
	ErrorSendMessageTypeNotFound     = api.ErrorSendMessageTypeNotFound("发送消息类型错误")
	IsSendMessageContentRequired     = api.ErrorSendMessageContentRequired("发送消息内容必填")
	ErrorToAccountAndFromAccountSame = api.ErrorToAccountAndFromAccountSame("接收者账号不能和发送者账号一致")
	ErrorSendMessage                 = api.ErrorSendMessage("发送消息失败, 请重试")
)
