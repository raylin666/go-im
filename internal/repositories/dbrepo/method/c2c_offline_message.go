package method

import "gorm.io/gen"

type C2COfflineMessage interface {
	// FirstByAccount where("`from_account`=@fromAccount and to_account`=@toAccount")
	FirstByAccount(fromAccount string, toAccount string) (gen.T, error)
}
