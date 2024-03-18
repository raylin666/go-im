package method

import "gorm.io/gen"

type App interface {
	// FirstByKeyAndSecret where("`key`=@key and `secret`=@secret")
	FirstByKeyAndSecret(key uint64, secret string) (gen.T, error)
}
