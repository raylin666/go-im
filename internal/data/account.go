package data

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gen/field"
	"gorm.io/gorm"
	"mt/internal/app"
	"mt/internal/biz"
	"mt/internal/constant/defined"
	typeAccount "mt/internal/constant/types/account"
	"mt/internal/repositories/dbrepo"
	"mt/internal/repositories/dbrepo/model"
	"mt/pkg/repositories"
	"time"
)

type accountRepo struct {
	data  repositories.DataRepo
	tools *app.Tools
}

func NewAccountRepo(repo repositories.DataRepo, tools *app.Tools) biz.AccountRepo {
	return &accountRepo{
		data:  repo,
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

	q := dbrepo.NewDefaultDbQuery(r.data.DbRepo()).Account
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
	q := dbrepo.NewDefaultDbQuery(r.data.DbRepo()).Account
	account, dataExistErr := q.WithContext(ctx).FirstByAccountId(accountId)
	if dataExistErr != nil {
		if errors.Is(dataExistErr, gorm.ErrRecordNotFound) {
			return nil, defined.ErrorDataNotFound
		}

		return nil, defined.ErrorDataSelectError
	}

	originAccount := account
	account.Nickname = data.Nickname
	account.Avatar = data.Avatar
	if data.IsAdmin {
		account.IsAdmin = 1
	} else {
		account.IsAdmin = 0
	}

	if updateDataErr := q.WithContext(ctx).Where(q.AccountId.Eq(accountId)).Save(&account); updateDataErr != nil {
		r.tools.Logger().UseSQL(ctx).Error("更新账号错误", zap.Any("origin_account", originAccount), zap.Any("account", account), zap.Error(updateDataErr))
		return nil, defined.ErrorDataUpdateError
	}

	return &account, nil
}

// Delete 删除账号
func (r *accountRepo) Delete(ctx context.Context, accountId string) (*model.Account, error) {
	q := dbrepo.NewDefaultDbQuery(r.data.DbRepo()).Account
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

// GetInfo 获取账号信息
func (r *accountRepo) GetInfo(ctx context.Context, accountId string) (*model.Account, error) {
	q := dbrepo.NewDefaultDbQuery(r.data.DbRepo()).Account
	account, dataExistErr := q.WithContext(ctx).FirstByAccountId(accountId)
	if dataExistErr != nil {
		if errors.Is(dataExistErr, gorm.ErrRecordNotFound) {
			return nil, defined.ErrorDataNotFound
		}

		return nil, defined.ErrorDataSelectError
	}

	return &account, nil
}

// Login 登录帐号
func (r *accountRepo) Login(ctx context.Context, accountId string, data *typeAccount.LoginRequest) (*model.Account, error) {
	q := dbrepo.NewDefaultDbQuery(r.data.DbRepo()).Account
	originAccount, dataExistErr := q.WithContext(ctx).FirstByAccountId(accountId)
	if dataExistErr != nil {
		if errors.Is(dataExistErr, gorm.ErrRecordNotFound) {
			return nil, defined.ErrorDataNotFound
		}

		return nil, defined.ErrorDataSelectError
	}

	var (
		timeNow          = time.Now()
		isOnline    int8 = 1
		lastLoginIp      = data.ClientIp

		accountOnline = new(model.AccountOnline)
	)

	accountOnline.AccountId = accountId
	accountOnline.LoginTime = timeNow
	accountOnline.LoginIp = data.ClientIp
	accountOnline.ClientAddr = data.ClientAddr
	accountOnline.ServerAddr = data.ServerAddr
	accountOnline.DeviceId = data.DeviceId
	accountOnline.Os = data.Os
	accountOnline.System = data.System

	accountOnlineQuery := dbrepo.NewDefaultDbQuery(r.data.DbRepo()).AccountOnline
	if err := accountOnlineQuery.WithContext(ctx).Create(accountOnline); err != nil {
		return nil, err
	}

	account := originAccount
	account.IsOnline = isOnline
	account.LastLoginTime = &timeNow
	account.LastLoginIp = lastLoginIp
	account.UpdatedAt = timeNow

	assignExpr := []field.AssignExpr{
		q.IsOnline.Value(isOnline),
		q.LastLoginTime.Value(timeNow),
		q.LastLoginIp.Value(lastLoginIp),
		q.UpdatedAt.Value(timeNow),
	}

	if originAccount.FirstLoginTime == nil {
		if accountOnlineExistsResult, err := accountOnlineQuery.WithContext(ctx).ExistsByAccountId(originAccount.AccountId); err == nil {
			if existsResult, existsResultOk := accountOnlineExistsResult["ok"]; existsResultOk {
				existsValue, existsValueOk := existsResult.(int64)
				if existsValueOk && existsValue == 0 {
					account.FirstLoginTime = &timeNow
					assignExpr = append(assignExpr, q.FirstLoginTime.Value(timeNow))
				}
			}
		}
	}

	if _, updateDataErr := q.WithContext(ctx).Where(q.AccountId.Eq(originAccount.AccountId)).UpdateSimple(assignExpr...); updateDataErr != nil {
		r.tools.Logger().UseSQL(ctx).Error("更新账号登录信息错误", zap.Any("origin_account", originAccount), zap.Any("account", account), zap.Error(updateDataErr))
		return nil, defined.ErrorDataUpdateError
	}

	return &account, nil
}
