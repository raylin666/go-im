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
	"mt/internal/repositories/dbrepo/query"
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

	dbQuery := dbrepo.NewDefaultDbQuery(r.data.DbRepo())
	if _, dataExistErr := dbQuery.Account.WithContext(ctx).Where().FirstByAccountId(account.AccountId); !errors.Is(dataExistErr, gorm.ErrRecordNotFound) {
		return nil, defined.ErrorDataAlreadyExists
	}
	if createDataErr := dbQuery.Account.WithContext(ctx).Create(account); createDataErr != nil {
		r.tools.Logger().UseSQL(ctx).Error("创建账号错误", zap.Any("account", account), zap.Error(createDataErr))
		return nil, defined.ErrorDataAdd
	}

	return account, nil
}

// Update 更新账号
func (r *accountRepo) Update(ctx context.Context, accountId string, data *typeAccount.UpdateRequest) (*model.Account, error) {
	dbQuery := dbrepo.NewDefaultDbQuery(r.data.DbRepo())
	account, dataExistErr := dbQuery.Account.WithContext(ctx).FirstByAccountId(accountId)
	if dataExistErr != nil {
		if errors.Is(dataExistErr, gorm.ErrRecordNotFound) {
			return nil, defined.ErrorDataNotFound
		}

		return nil, defined.ErrorDataSelect
	}

	originAccount := account
	account.Nickname = data.Nickname
	account.Avatar = data.Avatar
	if data.IsAdmin {
		account.IsAdmin = 1
	} else {
		account.IsAdmin = 0
	}

	if updateDataErr := dbQuery.Account.WithContext(ctx).Where(dbQuery.Account.AccountId.Eq(accountId)).Save(&account); updateDataErr != nil {
		r.tools.Logger().UseSQL(ctx).Error("更新账号错误", zap.Any("origin_account", originAccount), zap.Any("account", account), zap.Error(updateDataErr))
		return nil, defined.ErrorDataUpdate
	}

	return &account, nil
}

// Delete 删除账号
func (r *accountRepo) Delete(ctx context.Context, accountId string) (*model.Account, error) {
	dbQuery := dbrepo.NewDefaultDbQuery(r.data.DbRepo())
	account, dataExistErr := dbQuery.Account.WithContext(ctx).FirstByAccountId(accountId)
	if errors.Is(dataExistErr, gorm.ErrRecordNotFound) {
		return nil, defined.ErrorDataNotFound
	}
	if result, deleteDataErr := dbQuery.Account.WithContext(ctx).Where(dbQuery.Account.AccountId.Eq(accountId)).Delete(&account); deleteDataErr != nil {
		r.tools.Logger().UseSQL(ctx).Error("删除账号错误", zap.Any("account_id", accountId), zap.Any("result", result), zap.Error(deleteDataErr))
		return nil, defined.ErrorDataDelete
	}

	return &account, nil
}

// GetInfo 获取账号信息
func (r *accountRepo) GetInfo(ctx context.Context, accountId string) (*model.Account, error) {
	dbQuery := dbrepo.NewDefaultDbQuery(r.data.DbRepo())
	account, dataExistErr := dbQuery.Account.WithContext(ctx).FirstByAccountId(accountId)
	if dataExistErr != nil {
		if errors.Is(dataExistErr, gorm.ErrRecordNotFound) {
			return nil, defined.ErrorDataNotFound
		}

		return nil, defined.ErrorDataSelect
	}

	return &account, nil
}

// Login 登录帐号
func (r *accountRepo) Login(ctx context.Context, accountId string, data *typeAccount.LoginRequest) (*model.Account, *model.AccountOnline, error) {
	var dbQuery = dbrepo.NewDefaultDbQuery(r.data.DbRepo())

	originAccount, dataExistErr := dbQuery.Account.WithContext(ctx).FirstByAccountId(accountId)
	if dataExistErr != nil {
		if errors.Is(dataExistErr, gorm.ErrRecordNotFound) {
			return nil, nil, defined.ErrorAccountNotFound
		}

		return nil, nil, defined.ErrorDataSelect
	}

	// 校验同客户端是否已登录
	if checkClientAccountOnlineResult, err := dbQuery.AccountOnline.WithContext(ctx).CheckClientIsOnline(data.ClientAddr, data.ServerAddr); err == nil {
		if existsResult, existsResultOk := checkClientAccountOnlineResult["ok"]; existsResultOk {
			existsValue, existsValueOk := existsResult.(int64)
			if existsValueOk && existsValue > 0 {
				return nil, nil, defined.ErrorAccountIsLogin
			}
		}
	}

	var (
		timeNow          = time.Now()
		isOnline    int8 = 1
		lastLoginIp      = data.ClientIp

		account       = originAccount
		accountOnline = new(model.AccountOnline)
	)

	err := dbrepo.NewDefaultDbQuery(r.data.DbRepo()).Transaction(func(tx *query.Query) error {
		// 返回任何错误都会回滚事务
		accountOnline.AccountId = accountId
		accountOnline.LoginTime = timeNow
		accountOnline.LoginIp = data.ClientIp
		accountOnline.ClientAddr = data.ClientAddr
		accountOnline.ServerAddr = data.ServerAddr
		accountOnline.DeviceId = data.DeviceId
		accountOnline.Os = data.Os
		accountOnline.System = data.System
		if err := tx.AccountOnline.WithContext(ctx).Create(accountOnline); err != nil {
			r.tools.Logger().UseSQL(ctx).Error("帐号登录失败: 写入帐号在线表失败", zap.Any("account", account), zap.Any("account_online", accountOnline), zap.Error(err))
			return defined.ErrorDataAdd
		}

		account.IsOnline = isOnline
		account.LastLoginTime = &timeNow
		account.LastLoginIp = lastLoginIp
		account.UpdatedAt = timeNow
		assignExpr := []field.AssignExpr{
			dbQuery.Account.IsOnline.Value(isOnline),
			dbQuery.Account.LastLoginTime.Value(timeNow),
			dbQuery.Account.LastLoginIp.Value(lastLoginIp),
			dbQuery.Account.UpdatedAt.Value(timeNow),
		}
		if originAccount.FirstLoginTime == nil {
			account.FirstLoginTime = &timeNow
			assignExpr = append(assignExpr, dbQuery.Account.FirstLoginTime.Value(timeNow))
		}

		if _, updateDataErr := tx.Account.WithContext(ctx).Where(dbQuery.Account.AccountId.Eq(originAccount.AccountId)).UpdateSimple(assignExpr...); updateDataErr != nil {
			r.tools.Logger().UseSQL(ctx).Error("帐号登录失败: 更新账号登录信息错误", zap.Any("origin_account", originAccount), zap.Any("account", account), zap.Any("account_online", accountOnline), zap.Error(updateDataErr))
			return defined.ErrorDataUpdate
		}

		// 返回 nil 提交事务
		return nil
	})

	return &account, accountOnline, err
}

// Logout 登出帐号
func (r *accountRepo) Logout(ctx context.Context, accountId string, data *typeAccount.LogoutRequest) (*model.AccountOnline, error) {
	var dbQuery = dbrepo.NewDefaultDbQuery(r.data.DbRepo())

	accountOnline, dataExistErr := dbQuery.AccountOnline.WithContext(ctx).FirstByOnlineId(data.OnlineId)
	if dataExistErr != nil {
		if errors.Is(dataExistErr, gorm.ErrRecordNotFound) {
			return nil, defined.ErrorDataNotFound
		}

		return nil, defined.ErrorDataSelect
	}

	if data.ClientIp != nil {
		accountOnline.LogoutIp = *data.ClientIp
	}

	var state int8 = model.AccountOnlineLoginStateNormal
	if data.State > 0 {
		state = data.State
	}

	timeNow := time.Now()
	accountOnline.LogoutTime = &timeNow
	accountOnline.LogoutState = state
	if updateDataErr := dbQuery.AccountOnline.WithContext(ctx).Where(dbQuery.AccountOnline.AccountId.Eq(accountId)).Save(&accountOnline); updateDataErr != nil {
		r.tools.Logger().UseSQL(ctx).Error("帐号登出失败: 更新在线账号错误", zap.Any("account_online", accountOnline), zap.Error(updateDataErr))
		return nil, defined.ErrorDataUpdate
	}

	return &accountOnline, nil
}
