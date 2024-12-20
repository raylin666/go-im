package method

import "gorm.io/gen"

type AccountOnline interface {
	// CheckClientIsOnline SELECT EXISTS (SELECT * FROM @@table WHERE `client_addr`=@clientAddr AND `server_addr` = @serverAddr AND `logout_time` IS NULL) AS `ok`
	CheckClientIsOnline(clientAddr string, serverAddr string) (gen.M, error)
	// FirstByOnlineId WHERE("`id` = @onlineId")
	FirstByOnlineId(onlineId int) (gen.T, error)
}
