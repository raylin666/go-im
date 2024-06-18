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
		cacheOnlineAccountsHash 管理所有在线账号, 无法区分账号所在哪个服务器登录, 只是标识在线状态

		cacheOnlineAccountHash 具体的在线账号缓存, 用来标识账号所在服务器及具体信息
	*/

	// 缓存存储用户在线用户哈希
	cacheOnlineAccountsHash       = "im:online:accounts"
	cacheOnlineAccountsHashExpire = 365 * 24 * 60 * 60

	cacheOnlineAccountHash       = "im:online:account:%s"
	cacheOnlineAccountHashExpire = 3 * 24 * 60 * 60
)

type Account struct {
	ID       string `json:"id"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

func NewAccount(id, nickname, avatar string) Account {
	return Account{
		ID:       id,
		Nickname: nickname,
		Avatar:   avatar,
	}
}

// OnlineAccount 在线账号信息
type OnlineAccount struct {
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

// getOnlineAccountsCacheKey 获取所有账号在线缓存数据KEY
func getOnlineAccountsCacheKey() string { return cacheOnlineAccountsHash }

// SetOnlineAccount 设置在线账号缓存数据
func SetOnlineAccount(onlineAccount *OnlineAccount) {
	var accountId = onlineAccount.AccountId
	valueByte, err := json.Marshal(onlineAccount)
	if err != nil {
		return
	}

	ctx := context.TODO()
	cacheKey := getOnlineAccountsCacheKey()
	redisClient := redisrepo.NewDefaultClient(RedisRepo())
	isOk, _ := redisClient.Exists(ctx, cacheKey).Result()
	jsonString := string(valueByte)
	_, err = redisClient.HSet(ctx, cacheKey, accountId, jsonString).Result()
	if err != nil {
		Logger(ctx).Error("设置账号存储缓存数据失败", zap.String("cache_key", cacheKey), zap.String("data", jsonString), zap.Error(err))

		return
	}

	if isOk == 0 {
		redisClient.Expire(ctx, cacheKey, cacheOnlineAccountsHashExpire*time.Second)
	}

	return
}

// SetOfflineAccount 设置账号离线(清理缓存数据)
func SetOfflineAccount(accountId string) bool {
	ctx := context.TODO()
	cacheKey := getOnlineAccountsCacheKey()
	redisClient := redisrepo.NewDefaultClient(RedisRepo())
	isOk, _ := redisClient.HDel(ctx, cacheKey, accountId).Result()
	if isOk == 0 {
		return false
	}

	return true
}

// GetOnlineAccount 获取在线账号缓存数据
func GetOnlineAccount(accountId string) *OnlineAccount {
	ctx := context.TODO()
	cacheKey := getOnlineAccountsCacheKey()
	redisClient := redisrepo.NewDefaultClient(RedisRepo())
	jsonOnlineAccount, err := redisClient.HGet(ctx, cacheKey, accountId).Result()
	if err != nil {
		return nil
	}

	var onlineAccount = &OnlineAccount{}
	err = json.Unmarshal([]byte(jsonOnlineAccount), onlineAccount)
	if err != nil {
		return nil
	}

	return onlineAccount
}

// IsOnlineAccount 判断账号是否在线
func IsOnlineAccount(accountId string) bool {
	accountOnline := GetOnlineAccount(accountId)
	if accountOnline == nil {
		return false
	}

	return true
}
