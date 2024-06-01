package method

import "gorm.io/gen"

type Account interface {
	// FirstByAccountId where("`account_id`=@accountId")
	FirstByAccountId(accountId string) (gen.T, error)
}
