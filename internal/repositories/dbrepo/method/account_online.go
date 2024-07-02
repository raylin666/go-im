package method

import "gorm.io/gen"

type AccountOnline interface {
	// SELECT 1 FROM @@table a WHERE EXISTS (select * from @@table where a.`account_id`=@accountId)
	ExistsByAccountId(accountId string) (gen.T, error)
}
