package method

import "gorm.io/gen"

type AccountOnline interface {
	// ClientIsOnline SELECT EXISTS (SELECT * FROM @@table WHERE `client_addr`=@clientAddr AND `server_addr` = @serverAddr AND `logout_time` IS NULL) AS `ok`
	ClientIsOnline(clientAddr string, serverAddr string) (gen.M, error)
	// IsOnline SELECT EXISTS (SELECT * FROM @@table WHERE `account_id`=@accountId AND `logout_time` IS NULL) AS `ok`
	IsOnline(accountId string) (gen.M, error)
	// FirstByOnlineId WHERE("`id` = @onlineId")
	FirstByOnlineId(onlineId int) (gen.T, error)
}
