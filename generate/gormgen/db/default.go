package db

import (
	"gorm.io/gen"
	"mt/internal/repositories/dbrepo/method"
	"mt/internal/repositories/dbrepo/model"
	"mt/pkg/db"
)

func NewGeneratorDefaultDb(dbInterface db.Db, outPath string) {
	g := gen.NewGenerator(gen.Config{
		// 生成目录存放位置
		OutPath: outPath,
		// WithContext 模式
		Mode: gen.WithDefaultQuery,
	})

	g.UseDB(dbInterface.Get().DB())

	accountModel := model.Account{}
	accountOnlineModel := model.AccountOnline{}

	// apply basic crud api on structs or table models which is specified by table name with function
	// GenerateModel/GenerateModelAs. And generator will generate table models' code when calling Execute.
	g.ApplyBasic(
		accountModel,
		accountOnlineModel,
	)

	// apply diy interfaces on structs or table models
	g.ApplyInterface(func(method method.Account) {}, accountModel)
	g.ApplyInterface(func(method method.AccountOnline) {}, accountOnlineModel)

	g.Execute()
}
