package event

import (
	"context"
	"encoding/json"
	"go.uber.org/zap"
	"mt/internal/constant/defined"
	"mt/internal/repositories/dbrepo/model"
	"mt/internal/websocket"
	typesEvent "mt/internal/websocket/types/event"
	typesMessage "mt/internal/websocket/types/message"
	"mt/pkg/utils"
	"time"
)

// SendC2CMessage 发送C2C消息 [消息事件处理]
func (event *events) SendC2CMessage(ctx context.Context, client *websocket.Client, seq string, message []byte) (code uint32, msg string, data interface{}, send bool) {
	code, msg, _, send = defaultEventResponse()

	var loggerFields = event.loggerFields(websocket.EventSendC2CMessage, seq, message)

	// 解析数据包
	request := &typesEvent.SendC2CMessageRequest{}
	err := json.Unmarshal(message, request)
	if err != nil {
		loggerFields = append(loggerFields, zap.Error(err))
		event.logger(ctx).Error("消息发送事件-解析消息数据包错误 json.Marshal", loggerFields...)
		code, msg = utils.ErrorMessage(defined.ErrorRequestParamsError)
		return
	}
	if request.Type == "" || !typesMessage.HasType(request.Type) {
		code, msg = utils.ErrorMessage(defined.ErrorSendMessageTypeNotFound)
		return
	}
	if request.Content == "" {
		code, msg = utils.ErrorMessage(defined.IsSendMessageContentRequired)
		return
	}
	if request.ToAccount == "" {
		code, msg = utils.ErrorMessage(defined.ErrorToAccountNotFound)
		return
	}
	// 判断是否等于当前登录用户
	if request.ToAccount == client.Account.ID {
		code, msg = utils.ErrorMessage(defined.ErrorToAccountAndFromAccountSame)
		return
	}

	// 根据账号ID查询账号是否存在
	accountQuery := event.dbQuery().Account
	if accountExistsResult, err := accountQuery.WithContext(ctx).ExistsByAccountId(client.Account.ID); err == nil {
		if existsResult, existsResultOk := accountExistsResult["ok"]; existsResultOk {
			existsValue, existsValueOk := existsResult.(int64)
			if existsValueOk && existsValue == 0 {
				code, msg = utils.ErrorMessage(defined.ErrorToAccountNotFound)
				return
			}
		}
	}

	var (
		content typesMessage.BuilderType

		sendMessageTime = time.Now().Unix()
	)

	switch request.Type {
	case typesMessage.TypeText:
		content = typesMessage.BuilderTextType{Text: request.Content}
	}

	// 写入消息表数据
	c2cMessageData, err := json.Marshal(content)
	if err != nil {
		code, msg = utils.ErrorMessage(defined.ErrorSendMessageError)
		return
	}

	c2cMessageModel := &model.C2CMessage{
		FromAccountId: client.Account.ID,
		ToAccountId:   request.ToAccount,
		MsgType:       request.Type,
		Data:          string(c2cMessageData),
		SendTime:      time.Unix(sendMessageTime, 0),
	}

	queryC2cMessage := event.dbQuery().C2CMessage
	err = queryC2cMessage.WithContext(ctx).Create(c2cMessageModel)
	if err != nil {
		code, msg = utils.ErrorMessage(defined.ErrorSendMessageError)
		return
	}

	data = typesEvent.SendC2CMessageResponse{
		MessageId:       c2cMessageModel.ID,
		FromAccount:     c2cMessageModel.FromAccountId,
		ToAccount:       c2cMessageModel.ToAccountId,
		MessageSendTime: c2cMessageModel.SendTime.Unix(),
		Type:            c2cMessageModel.MsgType,
		Content:         content,
	}

	// TODO 判断对端是否在线。 在线则直接下发消息给对端, 如果离线则写入离线消息
	/*toUserOnline := websocket.IsAccountOnline(c2cMessageModel.ToAccountId)
	if !toUserOnline {
		// TODO 对端离线
		tableName := model.C2COfflineMessageTableName(client.AppKey)
		queryC2cOfflineMessage := newDbQuery.C2COfflineMessage.Table(tableName)
		c2cOfflineMessageModel, _ := queryC2cOfflineMessage.WithContext(ctx).FirstByUser(c2cMessageModel.FromUserId, c2cMessageModel.ToUserId)
		if c2cOfflineMessageModel.ID == 0 {
			queryC2cOfflineMessage.WithContext(ctx).Create(&model.C2COfflineMessage{
				MessageId:  c2cMessageModel.ID,
				FromUserId: c2cMessageModel.FromUserId,
				ToUserId:   c2cMessageModel.ToUserId,
			})
		} else if c2cOfflineMessageModel.MessageId == 0 {
			queryC2cOfflineMessage.WithContext(ctx).Save(&model.C2COfflineMessage{
				MessageId: c2cMessageModel.ID,
			})
		}
	} else {
		// TODO 对端在线, 将消息发送给对端
		clientMessage := NewMessage(ctx, client.AppKey)
		clientMessage.SendMessageToUser(seq, c2cMessageModel.ToUserId, code, msg, data)
	}*/

	//loggerFields = append(loggerFields, zap.Bool("to_account_online", toUserOnline), zap.Any("response", data))
	event.logger(ctx).Info("消息发送事件-发送成功", loggerFields...)

	return
}
