package websocket

import (
	"context"
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"mt/internal/constant/defined"
	"mt/internal/lib"
	"mt/internal/repositories/dbrepo"
	"mt/internal/repositories/dbrepo/model"
	typesEvent "mt/internal/websocket/types/event"
	"mt/pkg/utils"
	"time"
)

type EventDisposeFunc func(ctx context.Context, client *Client, seq string, message []byte) (code uint32, msg string, data interface{})

type Events struct {
	Registers map[string]EventDisposeFunc
}

func NewEvents() (events *Events) {
	events = &Events{}
	events.Registers = make(map[string]EventDisposeFunc)

	// 注册处理事件
	events.Registers["ping"] = events.Ping
	events.Registers["login"] = events.Login
	events.Registers["logout"] = events.Logout
	events.Registers["loginStatus"] = events.LoginStatus

	return
}

// Ping 心跳检测
func (event *Events) Ping(ctx context.Context, client *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {
	code = codeStatusOk
	msg = codeMessageOk
	data = "pong"

	client.Heartbeat(uint64(time.Now().Unix()))

	return
}

// Login 账号登录 (必须登录完成后才能进行用户事件)
func (event *Events) Login(ctx context.Context, client *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {
	code = codeStatusOk
	msg = codeMessageOk
	data = nil

	// 设置账号数据表
	tableName := model.AccountTableName(client.AppKey)
	q := dbrepo.NewDefaultDbQuery(DbRepo()).Account.Table(tableName)

	// 获取账号信息
	accountFunc := func(userId string) (account model.Account, err error) {
		account, err = q.WithContext(ctx).FirstByUserId(userId)
		return
	}

	// 获取响应协议
	responseFunc := func(account model.Account, repeatLogin bool) typesEvent.LoginResponse {
		return typesEvent.LoginResponse{
			UserId:         account.UserId,
			Username:       account.Username,
			Avatar:         account.Avatar,
			IsAdmin:        account.IsAdmin == 1,
			Status:         model.AccountConvertStatus(account.Status),
			FirstLoginTime: *account.FirstLoginTime,
			LastLoginTime:  *account.LastLoginTime,
			LastLoginIp:    account.LastLoginIp,
			RepeatLogin:    repeatLogin,
		}
	}

	// 判断是否已登录
	if client.AccountOnline() {
		account, err := accountFunc(client.UserId)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			code, msg = utils.ErrorMessage(defined.ErrorAccountNotFound)
			return
		}

		data = responseFunc(account, true)
		return
	}

	// 解析数据包
	request := &typesEvent.LoginRequest{}
	err := json.Unmarshal(message, request)
	if err != nil || request.UserId == "" || request.Usersig == "" {
		Logger(ctx).Error("账号登录事件-解析消息数据包错误 json.Marshal", zap.Error(err))
		code, msg = utils.ErrorMessage(defined.ErrorRequestParamsError)
		return
	}

	// 查询用户信息
	account, err := accountFunc(request.UserId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		code, msg = utils.ErrorMessage(defined.ErrorAccountNotFound)
		return
	}

	currentTime := time.Now()
	account.Status = model.AccountStatusOnline
	account.LastLoginTime = &currentTime
	if account.FirstLoginTime == nil {
		account.FirstLoginTime = &currentTime
	}

	httpRequest := lib.GetContextHttpRequest(ctx)
	if httpRequest != nil {
		account.LastLoginIp = utils.ClientIP(httpRequest)
	}

	if err = q.WithContext(ctx).Save(&account); err != nil {
		Logger(ctx).Error("账号登录事件-数据写入失败", zap.String("table_name", tableName), zap.Any("data", account), zap.Error(err))
		code, msg = utils.ErrorMessage(defined.ErrorAccountLoginError)
		return
	}

	// TODO 账号登录成功, 更新连接账号数据
	data = responseFunc(account, false)
	client.AccountLogin(account.UserId, uint64(account.LastLoginTime.Unix()))

	Logger(ctx).Info("账号登录事件-登录成功",
		zap.String("user_id", account.UserId),
		zap.Time("first_time", *account.LastLoginTime),
		zap.String("login_ip", account.LastLoginIp),
		zap.Any("account", account),
		zap.Any("response", data))

	return
}

// Logout 账号登出
func (event *Events) Logout(ctx context.Context, client *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {
	code = codeStatusOk
	msg = codeMessageOk
	data = nil

	// 判断是否未登录
	if !client.AccountOnline() {
		code, msg = utils.ErrorMessage(defined.ErrorNotLoginError)
		return
	}

	userId := client.UserId
	currentTime := time.Now()
	data = typesEvent.LogoutResponse{UserId: userId, LogoutTime: currentTime}

	// TODO 账号登出成功, 更新连接账号数据
	client.AccountLogout()

	// 设置账号数据表
	tableName := model.AccountTableName(client.AppKey)
	q := dbrepo.NewDefaultDbQuery(DbRepo()).Account.Table(tableName)
	// 查询用户信息
	account, err := q.WithContext(ctx).FirstByUserId(userId)
	// 账号不存在也直接返回登出成功
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return
	}

	account.LastLogoutTime = &currentTime
	account.Status = model.AccountStatusOffline

	httpRequest := lib.GetContextHttpRequest(ctx)
	if httpRequest != nil {
		account.LastLogoutIp = utils.ClientIP(httpRequest)
	}

	if err = q.WithContext(ctx).Save(&account); err != nil {
		Logger(ctx).Error("账号登出事件-数据写入失败", zap.String("table_name", tableName), zap.Any("data", account), zap.Error(err))
		return
	}

	Logger(ctx).Info("账号登出事件-登出成功",
		zap.String("user_id", userId),
		zap.Any("account", account),
		zap.Any("response", data))

	return
}

// LoginStatus 获取登录状态
func (event *Events) LoginStatus(ctx context.Context, client *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {
	code = codeStatusOk
	msg = codeMessageOk

	var loginUserId string
	var loginStatus = "Logout"
	// 判断是否已登录
	if client.AccountOnline() {
		loginUserId = client.UserId
		loginStatus = "Login"
	}

	data = typesEvent.LoginStatusResponse{UserId: loginUserId, Status: loginStatus}

	return
}
