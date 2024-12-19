package method

import "gorm.io/gen"

type Account interface {
	// ExistsByAccountId SELECT EXISTS (SELECT * FROM @@table WHERE `account_id` = @accountId) AS `ok`
	ExistsByAccountId(accountId string) (gen.M, error)
	// FirstByAccountId WHERE("`account_id` = @accountId")
	FirstByAccountId(accountId string) (gen.T, error)
}
