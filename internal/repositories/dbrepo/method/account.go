package method

import "gorm.io/gen"

type Account interface {
	// FirstByUserId where("`user_id`=@userId")
	FirstByUserId(userId string) (gen.T, error)
}
