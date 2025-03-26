package data

import (
	"context"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"mt/errors"
	"mt/internal/app"
	"mt/internal/constant/types"
	"mt/internal/repositories"
	"mt/internal/repositories/dbrepo/model"
	"time"
)

type MessageRepo interface {
	SendC2CMessage(ctx context.Context, data *types.MessageSendC2CMessageRequest) (*model.C2CMessage, *model.C2COfflineMessage, error)
}

type messageRepo struct {
	data  repositories.DataRepo
	tools *app.Tools
}

func NewMessageRepo(repo repositories.DataRepo, tools *app.Tools) MessageRepo {
	return &messageRepo{
		data:  repo,
		tools: tools,
	}
}

// SendC2CMessage 发送 C2C 消息
func (r *messageRepo) SendC2CMessage(ctx context.Context, data *types.MessageSendC2CMessageRequest) (*model.C2CMessage, *model.C2COfflineMessage, error) {
	if data.Message == "" {
		return nil, nil, errors.New().SendMessageContentRequired()
	}

	if data.FromAccount == "" {
		return nil, nil, errors.New().FromAccountNotFound()
	}

	if data.ToAccount == "" {
		return nil, nil, errors.New().ToAccountNotFound()
	}

	if data.FromAccount == data.ToAccount {
		return nil, nil, errors.New().ToAccountAndFromAccountSame()
	}

	dbQuery := r.data.DefaultDbQuery()
	formAccountResult, dataExistErr := dbQuery.Account.WithContext(ctx).ExistsByAccountId(data.FromAccount)
	if dataExistErr != nil || formAccountResult["ok"] == 0 {
		return nil, nil, errors.New().FromAccountNotFound()
	}
	toAccountResult, dataExistErr := dbQuery.Account.WithContext(ctx).ExistsByAccountId(data.ToAccount)
	if dataExistErr != nil || toAccountResult["ok"] == 0 {
		return nil, nil, errors.New().ToAccountNotFound()
	}

	var c2cMessage = &model.C2CMessage{
		FromAccount: data.FromAccount,
		ToAccount:   data.ToAccount,
		Data:        data.Message,
		Status:      model.C2CMessageStatusOn,
		IsRevoke:    model.C2CMessageRevokeNo,
		SendAt:      time.Now(),
	}

	// TODO 保存消息记录
	if createDataErr := dbQuery.C2CMessage.WithContext(ctx).Create(c2cMessage); createDataErr != nil {
		r.tools.Logger().UseSQL(ctx).Error("C2C消息记录错误", zap.Any("c2c_message", c2cMessage), zap.Error(createDataErr))
		return nil, nil, errors.New(errors.WithMessage(createDataErr.Error())).DataAdd()
	}

	// TODO 离线消息记录
	c2cOfflineMessage, err := dbQuery.C2COfflineMessage.WithContext(ctx).FirstByAccount(data.FromAccount, data.ToAccount)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		r.tools.Logger().UseSQL(ctx).Error("C2C离线消息记录查询错误", zap.String("from_account", data.FromAccount), zap.String("to_account", data.ToAccount), zap.Error(err))
		return nil, nil, errors.New().Server()
	}

	if (err != nil && errors.Is(err, gorm.ErrRecordNotFound)) || c2cOfflineMessage.ID <= 0 {
		// 数据不存在, 创建数据
		c2cOfflineMessage = model.C2COfflineMessage{FromAccount: data.FromAccount, ToAccount: data.ToAccount}
		if err = dbQuery.C2COfflineMessage.WithContext(ctx).Create(&c2cOfflineMessage); err != nil {
			r.tools.Logger().UseSQL(ctx).Error("C2C离线消息记录错误", zap.Any("c2c_offline_message", c2cOfflineMessage), zap.Error(err))
			return nil, nil, errors.New().Server()
		}
	}

	accountOnlineResult, err := dbQuery.AccountOnline.WithContext(ctx).IsOnline(data.ToAccount)
	if err == nil && accountOnlineResult["ok"].(int64) == 0 {
		// 对端离线
		dbQuery.C2COfflineMessage.WithContext(ctx).Where(
			dbQuery.C2COfflineMessage.FromAccount.Eq(data.FromAccount),
			dbQuery.C2COfflineMessage.ToAccount.Eq(data.ToAccount),
		).UpdateSimple(
			dbQuery.C2COfflineMessage.MessageId.Value(c2cMessage.ID),
			dbQuery.C2COfflineMessage.UnreadNum.Value(c2cOfflineMessage.UnreadNum+1),
		)
	}

	return c2cMessage, &c2cOfflineMessage, nil
}
