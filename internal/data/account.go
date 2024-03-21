package data

import (
	"context"
	"fmt"
	"mt/internal/biz"
	"mt/internal/constant/types"
	"mt/internal/middleware/auth"
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

	fmt.Println(ctx.Value(auth.AppID))

	/*q := dbrepo.NewDefaultDbQuery(r.data.DbRepo).Account.Table(model.AccountTableName())
	if createDataErr := tx.WithContext(ctx).Create(app).Error; createDataErr != nil {
		tx.Rollback()
		r.log.UseSQL(ctx).Error("创建应用错误", zap.Any("data", app), zap.Error(createDataErr))
		return nil, defined.ErrorDataAddError
	}*/

	return account, nil
}
