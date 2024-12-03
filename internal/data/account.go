package data

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"mt/internal/app"
	"mt/internal/biz"
	"mt/internal/constant/defined"
	typeAccount "mt/internal/constant/types/account"
	"mt/internal/repositories/dbrepo"
	"mt/internal/repositories/dbrepo/model"
	"time"
)

type accountRepo struct {
	data  *Data
	tools *app.Tools
}

func NewAccountRepo(data *Data, tools *app.Tools) biz.AccountRepo {
	return &accountRepo{
		data:  data,
		tools: tools,
	}
}

// Create 创建账号
func (r *accountRepo) Create(ctx context.Context, data *typeAccount.CreateRequest) (*model.Account, error) {
	account := &model.Account{
		AccountId: data.AccountId,
		Nickname:  data.Nickname,
		Avatar:    data.Avatar,
		IsAdmin:   0,
	}

	if data.IsAdmin {
		account.IsAdmin = 1
	}

	account.CreatedAt = time.Now()

	q := dbrepo.NewDefaultDbQuery(r.data.DbRepo).Account
	if _, dataExistErr := q.WithContext(ctx).Where().FirstByAccountId(account.AccountId); !errors.Is(dataExistErr, gorm.ErrRecordNotFound) {
		return nil, defined.ErrorDataAlreadyExists
	}
	if createDataErr := q.WithContext(ctx).Create(account); createDataErr != nil {
		r.tools.Logger().UseSQL(ctx).Error("创建账号错误", zap.Any("account", account), zap.Error(createDataErr))
		return nil, defined.ErrorDataAddError
	}

	return account, nil
}

// Update 更新账号
func (r *accountRepo) Update(ctx context.Context, accountId string, data *typeAccount.UpdateRequest) (*model.Account, error) {
	q := dbrepo.NewDefaultDbQuery(r.data.DbRepo).Account
	account, dataExistErr := q.WithContext(ctx).FirstByAccountId(accountId)
	if errors.Is(dataExistErr, gorm.ErrRecordNotFound) {
		return nil, defined.ErrorDataNotFound
	}

	account.Nickname = data.Nickname
	account.Avatar = data.Avatar
	if data.IsAdmin {
		account.IsAdmin = 1
	} else {
		account.IsAdmin = 0
	}

	if updateDataErr := q.WithContext(ctx).Where(q.AccountId.Eq(accountId)).Save(&account); updateDataErr != nil {
		r.tools.Logger().UseSQL(ctx).Error("更新账号错误", zap.Any("account", account), zap.Error(updateDataErr))
		return nil, defined.ErrorDataUpdateError
	}

	return &account, nil
}

// Delete 删除账号
func (r *accountRepo) Delete(ctx context.Context, accountId string) (*model.Account, error) {
	q := dbrepo.NewDefaultDbQuery(r.data.DbRepo).Account
	account, dataExistErr := q.WithContext(ctx).FirstByAccountId(accountId)
	if errors.Is(dataExistErr, gorm.ErrRecordNotFound) {
		return nil, defined.ErrorDataNotFound
	}
	if result, deleteDataErr := q.WithContext(ctx).Where(q.AccountId.Eq(accountId)).Delete(&account); deleteDataErr != nil {
		r.tools.Logger().UseSQL(ctx).Error("删除账号错误", zap.Any("account_id", accountId), zap.Any("result", result), zap.Error(deleteDataErr))
		return nil, defined.ErrorDataDeleteError
	}

	return &account, nil
}
