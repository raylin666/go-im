package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"mt/internal/constant/defined"
	"mt/internal/websocket/model"
)

type Process struct {
	Client *Client
}

func NewProcess(client *Client) *Process {
	return &Process{client}
}

func (p *Process) HandlerMessage(ctx context.Context, message []byte) {
	var (
		logAddr    = zap.String("address", p.Client.Addr)
		logMessage = zap.String("message", string(message))
	)

	ManagerInstance().Logger().UseWebSocket(ctx).Info("进入消息处理", logAddr, logMessage)

	defer func() {
		if r := recover(); r != nil {
			ManagerInstance().Logger().UseWebSocket(ctx).Error("消息处理异常", logAddr, logMessage, zap.Any("recover", r))
		}
	}()

	// TODO 数据包合法性校验/解析消息数据包
	request := &model.Request{}
	err := json.Unmarshal(message, request)
	if err != nil {
		ManagerInstance().Logger().UseWebSocket(ctx).Error("数据包合法性校验失败 json.Unmarshal", logAddr, logMessage, zap.Error(err))

		// 返回错误给客户端

		return
	}

	requestData, err := json.Marshal(request.Data)
	if err != nil {
		ManagerInstance().Logger().UseWebSocket(ctx).Error("解析消息数据包错误 json.Marshal", logAddr, logMessage, zap.Error(err))

		// 返回错误给客户端

		return
	}

	// TODO 将处理完成的数据包返回给客户端
	seq := request.Seq
	event := request.Event

	var (
		responseCode    uint32
		responseMessage string
		responseData    interface{}

		logSeq   = zap.String("seq", seq)
		logEvent = zap.String("event", event)
		logData  = zap.String("data", string(requestData))
	)

	ManagerInstance().Logger().UseWebSocket(ctx).Info("解析消息数据包完成", logAddr, logSeq, logEvent, logData)

	// 采用 MAP 处理事件
	if value, ok := getHandlerEvent(event); ok {
		responseCode, responseMessage, responseData = value(ctx, p.Client, seq, requestData)
	} else {
		e := defined.ErrorCommandInvalidNotFound
		responseCode = uint32(e.GetCode())
		responseMessage = e.GetMessage()
		ManagerInstance().Logger().UseWebSocket(ctx).Warn(fmt.Sprintf("处理事件 %s 不存在!", event), logAddr, logSeq, logEvent, logData)
	}

	responseHead := model.NewResponseHead(seq, event, responseCode, responseMessage, responseData)

	headByte, err := json.Marshal(responseHead)
	if err != nil {
		ManagerInstance().Logger().UseWebSocket(ctx).Error("处理响应数据错误 json.Marshal", logAddr, logSeq, logEvent, logData,
			zap.Uint32("response_code", responseCode),
			zap.String("response_message", responseMessage),
			zap.Any("response_data", responseData),
			zap.Error(err))

		return
	}

	var (
		logAppId           = zap.Uint32("app_id", p.Client.AppId)
		logUserId          = zap.String("user_id", p.Client.UserId)
		logResponseMessage = zap.String("message", string(headByte))
	)

	ok := p.Client.SendMessage(ctx, headByte)
	if ok {
		ManagerInstance().Logger().UseWebSocket(ctx).Info("发送消息成功", logAddr, logAppId, logUserId, logResponseMessage)
	} else {
		ManagerInstance().Logger().UseWebSocket(ctx).Error("发送消息失败", logAddr, logAppId, logUserId, logResponseMessage)
	}

	return
}
