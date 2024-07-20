package websocket

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
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
func SetAccountOnline(accountId string, account *AccountOnline) {
	valueByte, err := json.Marshal(account)
	if err != nil {
		return
	}

	ctx := context.TODO()
	cacheKey := getAccountsOnlineCacheKey()
	redisClient := redisrepo.NewDefaultClient(RedisRepo())
	redisClient.Watch(ctx, func(tx *redis.Tx) error {
		tx.TxPipelined(ctx, func(pipeliner redis.Pipeliner) error {
			isOk, _ := pipeliner.Exists(ctx, cacheKey).Result()
			if isOk <= 0 {

			}

			return nil
		})

		return nil
	}, cacheKey)

	isOk, _ := redisClient.Exists(ctx, cacheKey).Result()
	jsonString := string(valueByte)
	_, err = redisClient.HSet(ctx, cacheKey, accountId, jsonString).Result()
	if err != nil {
		Logger(ctx).Error("设置在线账号存储缓存数据失败", zap.String("cache_key", cacheKey), zap.String("data", jsonString), zap.Error(err))

		return
	}

	if isOk == 0 {
		redisClient.Expire(ctx, cacheKey, cacheAccountsOnlineHashExpire*time.Second)
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
