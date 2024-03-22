package data

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"mt/internal/biz"
	"mt/internal/constant/defined"
	"mt/internal/constant/types"
	"mt/internal/lib"
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
		UserId:   data.UserId,
		Username: data.Username,
		Avatar:   data.Avatar,
		IsAdmin:  data.IsAdmin,
	}
	account.CreatedAt = time.Now()

	appid, err := lib.HeaderAppID(ctx)
	if err != nil {
		return nil, defined.ErrorNotVisitAuth
	}

	tableName := model.AccountTableName(appid.Key)
	q := dbrepo.NewDefaultDbQuery(r.data.DbRepo).Account.Table(tableName)
	if _, dataExistErr := q.WithContext(ctx).FirstByUserId(account.UserId); !errors.Is(dataExistErr, gorm.ErrRecordNotFound) {
		return nil, defined.ErrorDataAlreadyExists
	}
	if createDataErr := q.WithContext(ctx).Create(account); createDataErr != nil {
		r.log.UseSQL(ctx).Error("创建应用账号错误", zap.String("table_name", tableName), zap.Any("data", account), zap.Error(createDataErr))
		return nil, defined.ErrorDataAddError
	}

	return account, nil
}
