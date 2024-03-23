package defined

import api "mt/api/v1"

var (
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
	ErrorDataHandlerError       = api.ErrorDataHandlerError("数据处理失败")
	ErrorDataTableAlreadyExists = api.ErrorDataTableAlreadyExists("数据表已存在")
	ErrorDataTableCreateError   = api.ErrorDataTableCreateError("创建数据表失败")
	ErrorDataTableRenameError   = api.ErrorDataTableRenameError("重命名数据表失败")

	ErrorIdInvalidValueError    = api.ErrorIdInvalidValueError("无效ID值")
	ErrorCommandInvalidNotFound = api.ErrorCommandInvalidNotFound("无效的执行事件")
	ErrorRequestParamsError     = api.ErrorRequestParamsError("请求参数错误")

	ErrorNotVisitAuth   = api.ErrorNotVisitAuth("没有访问权限, 请联系管理员")
	ErrorAppAuthClose   = api.ErrorAppAuthClose("应用权限已被关闭, 请联系管理员")
	ErrorAppAuthExpired = api.ErrorAppAuthExpired("应用已过期, 请联系管理员")

	ErrorWebsocketUpgraderError = api.ErrorWebsocketUpgraderError("WebSocket 协议升级失败")

	ErrorNotLoginError     = api.ErrorNotLoginError("请先登录后再操作")
	ErrorAccountNotFound   = api.ErrorAccountNotFound("账号不存在")
	ErrorAccountLoginError = api.ErrorAccountLoginError("账号登录错误")
)
