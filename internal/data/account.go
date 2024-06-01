package data

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"mt/internal/biz"
	"mt/internal/constant/defined"
	"mt/internal/constant/types"
	"mt/internal/repositories/dbrepo"
	"mt/internal/repositories/dbrepo/model"
	"mt/pkg/logger"
	"time"
)

type accountRepo struct {
	data *Data
	log  *logger.Logger
}

func NewAccountRepo(data *Data, logger *logger.Logger) biz.AccountRepo {
	return &accountRepo{
		data: data,
		log:  logger,
	}
}

func (r *accountRepo) Create(ctx context.Context, data types.AccountCreateData) (*model.Account, error) {
	account := &model.Account{
		AccountId: data.AccountId,
		Nickname:  data.Nickname,
		Avatar:    data.Avatar,
		IsAdmin:   data.IsAdmin,
	}
	account.CreatedAt = time.Now()

	q := dbrepo.NewDefaultDbQuery(r.data.DbRepo).Account
	if _, dataExistErr := q.WithContext(ctx).FirstByAccountId(account.AccountId); !errors.Is(dataExistErr, gorm.ErrRecordNotFound) {
		return nil, defined.ErrorDataAlreadyExists
	}
	if createDataErr := q.WithContext(ctx).Create(account); createDataErr != nil {
		r.log.UseSQL(ctx).Error("创建账号错误", zap.Any("data", account), zap.Error(createDataErr))
		return nil, defined.ErrorDataAddError
	}

	return account, nil
}
