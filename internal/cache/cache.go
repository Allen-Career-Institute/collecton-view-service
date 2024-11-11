package cache

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewRedisStore, NewCacheImpl)