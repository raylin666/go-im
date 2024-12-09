package websocket

import "github.com/google/wire"

// ProviderSet is websocket providers.
var ProviderSet = wire.NewSet(NewManagement)
