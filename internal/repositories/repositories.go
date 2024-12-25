package repositories

import (
	"github.com/google/wire"
	"github.com/raylin666/go-utils/cache/redis"
	"mt/internal/repositories/dbrepo"
	"mt/internal/repositories/dbrepo/query"
	"mt/internal/repositories/redisrepo"
	"mt/pkg/repositories"
)

// ProviderSet is repo providers.
var ProviderSet = wire.NewSet(NewRepositories)

var _ DataRepo = (*Repositories)(nil)

type DataRepo interface {
	DefaultDbQuery() *query.Query
	DefaultRedisClient() redis.Client
}

type Repositories struct {
	Db struct {
		DefaultQuery *query.Query
	}
	Redis struct {
		DefaultClient redis.Client
	}
}

func NewRepositories(repo repositories.DataRepo) DataRepo {
	var dataRepo = new(Repositories)
	dataRepo.Db.DefaultQuery = dbrepo.NewDefaultDbQuery(repo.DbRepo())
	dataRepo.Redis.DefaultClient = redisrepo.NewDefaultClient(repo.RedisRepo())
	return dataRepo
}

func (repositories *Repositories) DefaultDbQuery() *query.Query {
	//TODO implement me

	return repositories.Db.DefaultQuery
}

func (repositories *Repositories) DefaultRedisClient() redis.Client {
	//TODO implement me

	return repositories.Redis.DefaultClient
}
