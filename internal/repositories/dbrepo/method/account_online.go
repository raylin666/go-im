package method

import "gorm.io/gen"

type AccountOnline interface {
	// ExistsByAccountId SELECT EXISTS (SELECT * FROM @@table WHERE `account_id`=@accountId) AS `ok`
	ExistsByAccountId(accountId string) (gen.M, error)
}
