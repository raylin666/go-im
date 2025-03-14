package data

import (
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewHeartbeatRepo, NewAccountRepo, NewMessageRepo)
