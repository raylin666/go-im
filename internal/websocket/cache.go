package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"mt/internal/app"
	"mt/internal/repositories/redisrepo"
	"time"
)

const (
	/**
	缓存存储用户在线用户哈希

	cacheAccountsOnlineHash 管理在线账号, 无法区分账号所在哪个服务器登录及账号基础信息
	*/
	cacheAccountsOnlineHash       = "im:accounts:online"
	cacheAccountsOnlineHashExpire = 365 * 24 * 60 * 60
)

// AccountOnline 缓存在线账号信息
type AccountOnline struct {
	Nickname string   `json:"nickname"` // 账号昵称
	Avatar   string   `json:"avatar"`   // 账号头像
	IsAdmin  bool     `json:"is_admin"` // 是否管理员
	Clients  []string `json:"clients"`  // 账号所在服务器地址
}

// getAccountsOnlineCacheKey 获取在线账号存储缓存数据KEY
func getAccountsOnlineCacheKey() string { return cacheAccountsOnlineHash }

// SetAccountOnline 设置在线账号存储缓存数据
func SetAccountOnline(accountId string, account *Account) {
	var err error
	var jsonString string
	var accountOnline = AccountOnline{
		Nickname: account.Nickname,
		Avatar:   account.Avatar,
		IsAdmin:  account.IsAdmin,
	}
	ctx := context.TODO()
	cacheKey := getAccountsOnlineCacheKey()
	redisClient := redisrepo.NewDefaultClient(RedisRepo())
	err = redisClient.Watch(ctx, func(tx *redis.Tx) error {
		_, err = tx.TxPipelined(ctx, func(pipeliner redis.Pipeliner) error {
			isOk, _ := pipeliner.HExists(ctx, cacheKey, accountId).Result()
			fmt.Println(isOk)
			if !isOk {
				accountOnline.Clients = append(accountOnline.Clients, app.LocalServerIp)
				valueByte, err := json.Marshal(accountOnline)
				if err != nil {
					return err
				}

				// 缓存数据不存在
				_, err = pipeliner.HSet(ctx, cacheKey, accountId, string(valueByte)).Result()
				if err != nil {
					return err
				}

				pipeliner.Expire(ctx, cacheKey, cacheAccountsOnlineHashExpire*time.Second)

				return nil
			}
			// 缓存数据已存在, 更新信息
			jsonString, err = pipeliner.HGet(ctx, cacheKey, accountId).Result()
			if err != nil {
				return err
			}

			fmt.Println(jsonString)

			return err
		})

		return nil
	}, cacheKey)

	if err != nil {
		Logger(ctx).Error("设置在线账号存储缓存事务处理失败", zap.String("cache_key", cacheKey), zap.String("data", jsonString), zap.Error(err))

		return
	}

	return
}

// SetAccountOffline 设置账号离线(清理缓存数据)
func SetAccountOffline(accountId string) bool {
	ctx := context.TODO()
	cacheKey := getAccountsOnlineCacheKey()
	redisClient := redisrepo.NewDefaultClient(RedisRepo())
	isOk, _ := redisClient.HDel(ctx, cacheKey, accountId).Result()
	if isOk == 0 {
		return false
	}

	// 当服务地址都不存在时, 删除缓存

	return true
}

// GetAccountOnline 获取在线账号缓存数据
func GetAccountOnline(accountId string) *AccountOnline {
	ctx := context.TODO()
	cacheKey := getAccountsOnlineCacheKey()
	redisClient := redisrepo.NewDefaultClient(RedisRepo())
	jsonAccountOnline, err := redisClient.HGet(ctx, cacheKey, accountId).Result()
	if err != nil {
		return nil
	}

	var accountOnline = &AccountOnline{}
	err = json.Unmarshal([]byte(jsonAccountOnline), accountOnline)
	if err != nil {
		return nil
	}

	return accountOnline
}

// IsAccountOnline 判断账号是否在线
func IsAccountOnline(accountId string) bool {
	accountOnline := GetAccountOnline(accountId)
	if accountOnline == nil {
		return false
	}

	return true
}
