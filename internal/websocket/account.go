package websocket

import (
	"context"
	"encoding/json"
	"go.uber.org/zap"
	"mt/internal/repositories/redisrepo"
	"time"
)

const (
	/*
		cacheAccountsOnlineHash 管理所有在线账号, 无法区分账号所在哪个服务器登录, 只是标识在线状态

		cacheAccountOnlineHash 具体的在线账号缓存, 用来标识账号所在服务器及具体信息
	*/

	// 缓存存储用户在线用户哈希
	cacheAccountsOnlineHash       = "im:accounts:online"
	cacheAccountsOnlineHashExpire = 365 * 24 * 60 * 60

	cacheAccountOnlineHash       = "im:account:online:%s"
	cacheAccountOnlineHashExpire = 3 * 24 * 60 * 60
)

type Account struct {
	ID       string `json:"id"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	OnlineId int    `json:"online_id"` // 账号在线ID, account_online 表
}

// NewAccount 创建账号
func NewAccount(id, nickname, avatar string) *Account {
	return &Account{
		ID:       id,
		Nickname: nickname,
		Avatar:   avatar,
	}
}

// WithOnlineId 添加账号在线ID
func (account *Account) WithOnlineId(id int) *Account {
	account.OnlineId = id
	return account
}

// AccountOnline 在线账号信息
type AccountOnline struct {
	Ip             string     `json:"ip"`               // 服务IP
	Port           string     `json:"port"`             // 服务端口
	ClientIp       string     `json:"client_ip"`        // 客户端IP
	ClientPort     string     `json:"client_port"`      // 客户端端口
	AccountId      string     `json:"account_id"`       // 账号ID
	Nickname       string     `json:"nickname"`         // 账号昵称
	Avatar         string     `json:"avatar"`           // 账号头像
	IsAdmin        bool       `json:"is_admin"`         // 是否管理员
	FirstLoginTime *time.Time `json:"first_login_time"` // 账号首次登录时间
	LastLoginTime  *time.Time `json:"last_login_time"`  // 账号最后登录时间
}

// getAccountsOnlineCacheKey 获取所有账号在线缓存数据KEY
func getAccountsOnlineCacheKey() string { return cacheAccountsOnlineHash }

// SetAccountOnline 设置在线账号缓存数据
func SetAccountOnline(accountOnline *AccountOnline) {
	var accountId = accountOnline.AccountId
	valueByte, err := json.Marshal(accountOnline)
	if err != nil {
		return
	}

	ctx := context.TODO()
	cacheKey := getAccountsOnlineCacheKey()
	redisClient := redisrepo.NewDefaultClient(RedisRepo())
	isOk, _ := redisClient.Exists(ctx, cacheKey).Result()
	jsonString := string(valueByte)
	_, err = redisClient.HSet(ctx, cacheKey, accountId, jsonString).Result()
	if err != nil {
		Logger(ctx).Error("设置账号存储缓存数据失败", zap.String("cache_key", cacheKey), zap.String("data", jsonString), zap.Error(err))

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
