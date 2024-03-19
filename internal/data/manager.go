package data

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"mt/internal/biz"
	"mt/internal/constant/defined"
	"mt/internal/constant/types"
	"mt/internal/repositories/dbrepo"
	"mt/internal/repositories/dbrepo/model"
	"mt/pkg/logger"
	"time"
)

type managerRepo struct {
	data *Data
	log  *logger.Logger
}

func NewManagerRepo(data *Data, logger *logger.Logger) biz.ManagerRepo {
	return &managerRepo{
		data: data,
		log:  logger,
	}
}

func (r *managerRepo) Create(ctx context.Context, data types.ManagerCreateData) (*model.App, error) {
	app := &model.App{
		Ident:     data.Ident,
		Name:      data.Name,
		Key:       data.Key,
		Secret:    data.Secret,
		Status:    data.Status,
		ExpiredAt: data.ExpiredAt,
	}
	app.CreatedAt = time.Now()

	var db = dbrepo.NewDefaultDb(r.data.DbRepo)
	tx := db.Begin()
	if createDataErr := tx.WithContext(ctx).Create(app).Error; createDataErr != nil {
		tx.Rollback()
		r.log.UseSQL(ctx).Error("创建应用错误", zap.Any("data", app), zap.Error(createDataErr))
		return nil, defined.ErrorDataAddError
	}

	// 创建用户账号表
	var tableName = fmt.Sprintf("account_%d", app.Key)
	if tx.Migrator().HasTable(tableName) {
		tx.Rollback()
		return nil, defined.ErrorDataTableAlreadyExists
	}

	acm := &model.Account{}
	if createTableErr := tx.Set("gorm:table_options", "ENGINE=InnoDB COMMENT='应用用户账号表'").Migrator().CreateTable(acm); createTableErr != nil {
		tx.Rollback()
		r.log.UseSQL(ctx).Error("创建应用用户账号数据表错误", zap.Error(createTableErr))
		return nil, defined.ErrorDataTableCreateError
	}

	if renameTableErr := tx.Migrator().RenameTable(acm, tableName); renameTableErr != nil {
		tx.Rollback()
		r.log.UseSQL(ctx).Error("重命名应用用户账号数据表错误", zap.Error(renameTableErr))
		return nil, defined.ErrorDataTableRenameError
	}

	tx.Commit()

	r.log.UseApp(ctx).Info("创建应用成功", zap.Any("data", app))

	return app, nil
}
