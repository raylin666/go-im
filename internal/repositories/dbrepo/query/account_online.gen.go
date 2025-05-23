// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package query

import (
	"context"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	"mt/internal/repositories/dbrepo/model"
)

func newAccountOnline(db *gorm.DB, opts ...gen.DOOption) accountOnline {
	_accountOnline := accountOnline{}

	_accountOnline.accountOnlineDo.UseDB(db, opts...)
	_accountOnline.accountOnlineDo.UseModel(&model.AccountOnline{})

	tableName := _accountOnline.accountOnlineDo.TableName()
	_accountOnline.ALL = field.NewAsterisk(tableName)
	_accountOnline.ID = field.NewInt(tableName, "id")
	_accountOnline.AccountId = field.NewString(tableName, "account_id")
	_accountOnline.LoginTime = field.NewTime(tableName, "login_time")
	_accountOnline.LogoutTime = field.NewTime(tableName, "logout_time")
	_accountOnline.LoginIp = field.NewString(tableName, "login_ip")
	_accountOnline.LogoutIp = field.NewString(tableName, "logout_ip")
	_accountOnline.LogoutState = field.NewInt8(tableName, "logout_state")
	_accountOnline.ClientAddr = field.NewString(tableName, "client_addr")
	_accountOnline.ServerAddr = field.NewString(tableName, "server_addr")
	_accountOnline.DeviceId = field.NewString(tableName, "device_id")
	_accountOnline.Os = field.NewString(tableName, "os")
	_accountOnline.System = field.NewString(tableName, "system")

	_accountOnline.fillFieldMap()

	return _accountOnline
}

type accountOnline struct {
	accountOnlineDo accountOnlineDo

	ALL         field.Asterisk
	ID          field.Int
	AccountId   field.String
	LoginTime   field.Time
	LogoutTime  field.Time
	LoginIp     field.String
	LogoutIp    field.String
	LogoutState field.Int8
	ClientAddr  field.String
	ServerAddr  field.String
	DeviceId    field.String
	Os          field.String
	System      field.String

	fieldMap map[string]field.Expr
}

func (a accountOnline) Table(newTableName string) *accountOnline {
	a.accountOnlineDo.UseTable(newTableName)
	return a.updateTableName(newTableName)
}

func (a accountOnline) As(alias string) *accountOnline {
	a.accountOnlineDo.DO = *(a.accountOnlineDo.As(alias).(*gen.DO))
	return a.updateTableName(alias)
}

func (a *accountOnline) updateTableName(table string) *accountOnline {
	a.ALL = field.NewAsterisk(table)
	a.ID = field.NewInt(table, "id")
	a.AccountId = field.NewString(table, "account_id")
	a.LoginTime = field.NewTime(table, "login_time")
	a.LogoutTime = field.NewTime(table, "logout_time")
	a.LoginIp = field.NewString(table, "login_ip")
	a.LogoutIp = field.NewString(table, "logout_ip")
	a.LogoutState = field.NewInt8(table, "logout_state")
	a.ClientAddr = field.NewString(table, "client_addr")
	a.ServerAddr = field.NewString(table, "server_addr")
	a.DeviceId = field.NewString(table, "device_id")
	a.Os = field.NewString(table, "os")
	a.System = field.NewString(table, "system")

	a.fillFieldMap()

	return a
}

func (a *accountOnline) WithContext(ctx context.Context) *accountOnlineDo {
	return a.accountOnlineDo.WithContext(ctx)
}

func (a accountOnline) TableName() string { return a.accountOnlineDo.TableName() }

func (a accountOnline) Alias() string { return a.accountOnlineDo.Alias() }

func (a *accountOnline) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := a.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (a *accountOnline) fillFieldMap() {
	a.fieldMap = make(map[string]field.Expr, 12)
	a.fieldMap["id"] = a.ID
	a.fieldMap["account_id"] = a.AccountId
	a.fieldMap["login_time"] = a.LoginTime
	a.fieldMap["logout_time"] = a.LogoutTime
	a.fieldMap["login_ip"] = a.LoginIp
	a.fieldMap["logout_ip"] = a.LogoutIp
	a.fieldMap["logout_state"] = a.LogoutState
	a.fieldMap["client_addr"] = a.ClientAddr
	a.fieldMap["server_addr"] = a.ServerAddr
	a.fieldMap["device_id"] = a.DeviceId
	a.fieldMap["os"] = a.Os
	a.fieldMap["system"] = a.System
}

func (a accountOnline) clone(db *gorm.DB) accountOnline {
	a.accountOnlineDo.ReplaceConnPool(db.Statement.ConnPool)
	return a
}

func (a accountOnline) replaceDB(db *gorm.DB) accountOnline {
	a.accountOnlineDo.ReplaceDB(db)
	return a
}

type accountOnlineDo struct{ gen.DO }

// ClientIsOnline SELECT EXISTS (SELECT * FROM @@table WHERE `client_addr`=@clientAddr AND `server_addr` = @serverAddr AND `logout_time` IS NULL) AS `ok`
func (a accountOnlineDo) ClientIsOnline(clientAddr string, serverAddr string) (result map[string]interface{}, err error) {
	var params []interface{}

	var generateSQL strings.Builder
	params = append(params, clientAddr)
	params = append(params, serverAddr)
	generateSQL.WriteString("SELECT EXISTS (SELECT * FROM account_online WHERE `client_addr`=? AND `server_addr` = ? AND `logout_time` IS NULL) AS `ok` ")

	result = make(map[string]interface{})
	var executeSQL *gorm.DB

	executeSQL = a.UnderlyingDB().Raw(generateSQL.String(), params...).Take(result)
	err = executeSQL.Error
	return
}

// IsOnline SELECT EXISTS (SELECT * FROM @@table WHERE `account_id`=@accountId AND `logout_time` IS NULL) AS `ok`
func (a accountOnlineDo) IsOnline(accountId string) (result map[string]interface{}, err error) {
	var params []interface{}

	var generateSQL strings.Builder
	params = append(params, accountId)
	generateSQL.WriteString("SELECT EXISTS (SELECT * FROM account_online WHERE `account_id`=? AND `logout_time` IS NULL) AS `ok` ")

	result = make(map[string]interface{})
	var executeSQL *gorm.DB

	executeSQL = a.UnderlyingDB().Raw(generateSQL.String(), params...).Take(result)
	err = executeSQL.Error
	return
}

// FirstByOnlineId WHERE("`id` = @onlineId")
func (a accountOnlineDo) FirstByOnlineId(onlineId int) (result model.AccountOnline, err error) {
	var params []interface{}

	var generateSQL strings.Builder
	params = append(params, onlineId)
	generateSQL.WriteString("`id` = ? ")

	var executeSQL *gorm.DB

	executeSQL = a.UnderlyingDB().Where(generateSQL.String(), params...).Take(&result)
	err = executeSQL.Error
	return
}

func (a accountOnlineDo) Debug() *accountOnlineDo {
	return a.withDO(a.DO.Debug())
}

func (a accountOnlineDo) WithContext(ctx context.Context) *accountOnlineDo {
	return a.withDO(a.DO.WithContext(ctx))
}

func (a accountOnlineDo) ReadDB() *accountOnlineDo {
	return a.Clauses(dbresolver.Read)
}

func (a accountOnlineDo) WriteDB() *accountOnlineDo {
	return a.Clauses(dbresolver.Write)
}

func (a accountOnlineDo) Session(config *gorm.Session) *accountOnlineDo {
	return a.withDO(a.DO.Session(config))
}

func (a accountOnlineDo) Clauses(conds ...clause.Expression) *accountOnlineDo {
	return a.withDO(a.DO.Clauses(conds...))
}

func (a accountOnlineDo) Returning(value interface{}, columns ...string) *accountOnlineDo {
	return a.withDO(a.DO.Returning(value, columns...))
}

func (a accountOnlineDo) Not(conds ...gen.Condition) *accountOnlineDo {
	return a.withDO(a.DO.Not(conds...))
}

func (a accountOnlineDo) Or(conds ...gen.Condition) *accountOnlineDo {
	return a.withDO(a.DO.Or(conds...))
}

func (a accountOnlineDo) Select(conds ...field.Expr) *accountOnlineDo {
	return a.withDO(a.DO.Select(conds...))
}

func (a accountOnlineDo) Where(conds ...gen.Condition) *accountOnlineDo {
	return a.withDO(a.DO.Where(conds...))
}

func (a accountOnlineDo) Exists(subquery interface{ UnderlyingDB() *gorm.DB }) *accountOnlineDo {
	return a.Where(field.CompareSubQuery(field.ExistsOp, nil, subquery.UnderlyingDB()))
}

func (a accountOnlineDo) Order(conds ...field.Expr) *accountOnlineDo {
	return a.withDO(a.DO.Order(conds...))
}

func (a accountOnlineDo) Distinct(cols ...field.Expr) *accountOnlineDo {
	return a.withDO(a.DO.Distinct(cols...))
}

func (a accountOnlineDo) Omit(cols ...field.Expr) *accountOnlineDo {
	return a.withDO(a.DO.Omit(cols...))
}

func (a accountOnlineDo) Join(table schema.Tabler, on ...field.Expr) *accountOnlineDo {
	return a.withDO(a.DO.Join(table, on...))
}

func (a accountOnlineDo) LeftJoin(table schema.Tabler, on ...field.Expr) *accountOnlineDo {
	return a.withDO(a.DO.LeftJoin(table, on...))
}

func (a accountOnlineDo) RightJoin(table schema.Tabler, on ...field.Expr) *accountOnlineDo {
	return a.withDO(a.DO.RightJoin(table, on...))
}

func (a accountOnlineDo) Group(cols ...field.Expr) *accountOnlineDo {
	return a.withDO(a.DO.Group(cols...))
}

func (a accountOnlineDo) Having(conds ...gen.Condition) *accountOnlineDo {
	return a.withDO(a.DO.Having(conds...))
}

func (a accountOnlineDo) Limit(limit int) *accountOnlineDo {
	return a.withDO(a.DO.Limit(limit))
}

func (a accountOnlineDo) Offset(offset int) *accountOnlineDo {
	return a.withDO(a.DO.Offset(offset))
}

func (a accountOnlineDo) Scopes(funcs ...func(gen.Dao) gen.Dao) *accountOnlineDo {
	return a.withDO(a.DO.Scopes(funcs...))
}

func (a accountOnlineDo) Unscoped() *accountOnlineDo {
	return a.withDO(a.DO.Unscoped())
}

func (a accountOnlineDo) Create(values ...*model.AccountOnline) error {
	if len(values) == 0 {
		return nil
	}
	return a.DO.Create(values)
}

func (a accountOnlineDo) CreateInBatches(values []*model.AccountOnline, batchSize int) error {
	return a.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (a accountOnlineDo) Save(values ...*model.AccountOnline) error {
	if len(values) == 0 {
		return nil
	}
	return a.DO.Save(values)
}

func (a accountOnlineDo) First() (*model.AccountOnline, error) {
	if result, err := a.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.AccountOnline), nil
	}
}

func (a accountOnlineDo) Take() (*model.AccountOnline, error) {
	if result, err := a.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.AccountOnline), nil
	}
}

func (a accountOnlineDo) Last() (*model.AccountOnline, error) {
	if result, err := a.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.AccountOnline), nil
	}
}

func (a accountOnlineDo) Find() ([]*model.AccountOnline, error) {
	result, err := a.DO.Find()
	return result.([]*model.AccountOnline), err
}

func (a accountOnlineDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.AccountOnline, err error) {
	buf := make([]*model.AccountOnline, 0, batchSize)
	err = a.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (a accountOnlineDo) FindInBatches(result *[]*model.AccountOnline, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return a.DO.FindInBatches(result, batchSize, fc)
}

func (a accountOnlineDo) Attrs(attrs ...field.AssignExpr) *accountOnlineDo {
	return a.withDO(a.DO.Attrs(attrs...))
}

func (a accountOnlineDo) Assign(attrs ...field.AssignExpr) *accountOnlineDo {
	return a.withDO(a.DO.Assign(attrs...))
}

func (a accountOnlineDo) Joins(fields ...field.RelationField) *accountOnlineDo {
	for _, _f := range fields {
		a = *a.withDO(a.DO.Joins(_f))
	}
	return &a
}

func (a accountOnlineDo) Preload(fields ...field.RelationField) *accountOnlineDo {
	for _, _f := range fields {
		a = *a.withDO(a.DO.Preload(_f))
	}
	return &a
}

func (a accountOnlineDo) FirstOrInit() (*model.AccountOnline, error) {
	if result, err := a.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.AccountOnline), nil
	}
}

func (a accountOnlineDo) FirstOrCreate() (*model.AccountOnline, error) {
	if result, err := a.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.AccountOnline), nil
	}
}

func (a accountOnlineDo) FindByPage(offset int, limit int) (result []*model.AccountOnline, count int64, err error) {
	result, err = a.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = a.Offset(-1).Limit(-1).Count()
	return
}

func (a accountOnlineDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = a.Count()
	if err != nil {
		return
	}

	err = a.Offset(offset).Limit(limit).Scan(result)
	return
}

func (a accountOnlineDo) Scan(result interface{}) (err error) {
	return a.DO.Scan(result)
}

func (a accountOnlineDo) Delete(models ...*model.AccountOnline) (result gen.ResultInfo, err error) {
	return a.DO.Delete(models)
}

func (a *accountOnlineDo) withDO(do gen.Dao) *accountOnlineDo {
	a.DO = *do.(*gen.DO)
	return a
}
