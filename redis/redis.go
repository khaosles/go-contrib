package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	gerr "github.com/khaosles/gtools2/utils/err"
	"go.uber.org/zap"

	"github.com/khaosles/go-contrib/core/config"
	glog "github.com/khaosles/go-contrib/core/log"
)

/*
   @File: redis.go
   @Author: khaosles
   @Time: 2023/3/3 22:34
   @Desc:
*/

var Cli *redis.Client

const APP = "redis"

type Redis struct {
	Addr               string        `mapstructure:"addr" default:"" yaml:"addr" json:"addr"`                                                   // Redis 服务器的地址
	Password           string        `mapstructure:"password" default:"" yaml:"password" json:"password"`                                       // Redis 服务器的密码
	DB                 int           `mapstructure:"db" default:"" yaml:"db" json:"db"`                                                         // Redis 服务器的数据库编号
	MaxRetries         int           `mapstructure:"max-retries" default:"3" yaml:"max-retries" json:"max-retries"`                             // Redis 客户端重试连接的最大次数
	DialTimeout        time.Duration `mapstructure:"dial-timeout" default:"5" yaml:"dial-timeout"  json:"dial-timeout"`                         // 连接 Redis 服务器的超时时间
	ReadTimeout        time.Duration `mapstructure:"read-timeout" default:"3" yaml:"read-timeout"  json:"read-timeout"`                         // 从 Redis 服务器读取数据的超时时间
	WriteTimeout       time.Duration `mapstructure:"write-timeout" default:"3" yaml:"write-timeout"  json:"write-timeout"`                      // 向 Redis 服务器写入数据的超时时间
	PoolSize           int           `mapstructure:"pool-size" default:"10" yaml:"pool-size" json:"pool-size" `                                 // 连接池的最大连接数
	MinIdleConns       int           `mapstructure:"min-idle-conns" default:"0" yaml:"min-idle-conns"  json:"min-idle-conns"`                   // 连接池的最小空闲连接数
	MaxConnAge         time.Duration `mapstructure:"max-conn-age" default:"0" yaml:"max-conn-age"  json:"max-conn-age"`                         // 连接池中连接的最大寿命，超过这个时间将被关闭
	PoolTimeout        time.Duration `mapstructure:"pool-timeout" default:"4" yaml:"pool-timeout"  json:"pool-timeout"`                         // 从连接池获取连接的超时时间
	IdleTimeout        time.Duration `mapstructure:"idle-timeout" default:"5" yaml:"idle-timeout"  json:"idle-timeout"`                         // 连接池中空闲连接的超时时间
	IdleCheckFrequency time.Duration `mapstructure:"idle-check-frequency" default:"1" yaml:"idle-check-frequency"  json:"idle-check-frequency"` // 连接池中空闲连接的检查频率
}

func init() {
	var err error
	var rds *Redis
	if err = config.Configuration(APP, rds); err != nil {
		glog.Fatal(err)
	}
	if Cli, err = New(rds); err != nil {
		glog.Fatal(err)
	}
}

func New(rds *Redis) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:               rds.Addr,
		Password:           rds.Password,
		DB:                 rds.DB,
		MaxRetries:         rds.MaxRetries,
		DialTimeout:        rds.DialTimeout,
		ReadTimeout:        rds.ReadTimeout,
		WriteTimeout:       rds.WriteTimeout,
		PoolSize:           rds.PoolSize,
		MinIdleConns:       rds.MinIdleConns,
		MaxConnAge:         rds.MaxConnAge,
		PoolTimeout:        rds.PoolTimeout,
		IdleTimeout:        rds.IdleTimeout,
		IdleCheckFrequency: rds.IdleCheckFrequency,
	})
	pong, err := Cli.Ping(context.Background()).Result()
	if err != nil {
		glog.Error("redis connect ping failed")
		return nil, err
	}
	glog.Info("redis connect ping response ===> ", zap.String("pong", pong))
	return rdb, nil
}

//  /////////////////////////// 字符串 ///////////////////////////

// Set 设置string
func Set(key string, value interface{}) error {
	err := Cli.Set(context.Background(), key, value, 0).Err()
	return err
}

// SetExpire 设置string
func SetExpire(key string, value interface{}, expireTime time.Duration) error {
	err := Cli.Set(context.Background(), key, value, expireTime).Err()
	return err
}

// Get 获取string
func Get(key string) (string, error) {
	// 判断key是否存在
	if !ExistKey(key) {
		return "", gerr.RedisKeyNotFoundException.New(key)
	}
	// 获取key
	val, err := Cli.Get(context.Background(), key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

//  /////////////////////////// hash表 ///////////////////////////

// HSet 设置单个哈希表字段
func HSet(key, field string, value any) error {
	err := Cli.HSet(context.Background(), key, field, value).Err()
	return err
}

// HSetMap 设置map对象位hash表
func HSetMap(key string, fields map[string]any) error {
	// 遍历字段
	for field, value := range fields {
		err := HSet(key, field, value)
		if err != nil {
			return err
		}
	}
	return nil
}

// HSetMapExpire 设置map对象位hash表
func HSetMapExpire(key string, fields map[string]any, expireTime time.Duration) error {
	// 遍历字段
	for field, value := range fields {
		err := HSet(key, field, value)
		if err != nil {
			return err
		}
	}
	return ExpireTime(key, expireTime)
}

// HGet 获取hash字段值
func HGet(key, field string) (string, error) {
	// 判断key是否存在
	if !ExistKey(key) {
		return "", gerr.RedisKeyNotFoundException.New(key)
	}
	// 获取数据
	result, err := Cli.HGet(context.Background(), key, field).Result()
	if err != nil {
		return "", err
	}
	return result, nil
}

// HGetAll 获取所有key 转成map
func HGetAll(key string) (map[string]string, error) {
	// 判断key是否存在
	if !ExistKey(key) {
		return nil, gerr.RedisKeyNotFoundException.New(key)
	}
	// 获取数据
	fields, err := Cli.HGetAll(context.Background(), key).Result()
	if err != nil {
		return nil, err
	}
	return fields, nil
}

// HExists 检查哈希表中给定字段是否存在
func HExists(key, field string) (bool, error) {
	return Cli.HExists(context.Background(), key, field).Result()
}

// HDel 删除哈希表中的一个或多个字段
func HDel(key string, fields ...string) (int64, error) {
	return Cli.HDel(context.Background(), key, fields...).Result()
}

// HKeys 获取哈希表中所有的字段
func HKeys(key string) ([]string, error) {
	return Cli.HKeys(context.Background(), key).Result()
}

// HVals 获取哈希表中所有字段的值
func HVals(key string) ([]string, error) {
	return Cli.HVals(context.Background(), key).Result()
}

// HLen 获取哈希表中字段的数量
func HLen(key string) (int64, error) {
	return Cli.HLen(context.Background(), key).Result()
}

//  /////////////////////////// 列表 ///////////////////////////

// LPush 将元素插入 Redis 的 list 头部
func LPush(key string, values ...interface{}) error {
	return Cli.LPush(context.Background(), key, values...).Err()
}

// RPush 将元素插入 Redis 的 list 尾部
func RPush(key string, values ...interface{}) error {
	return Cli.RPush(context.Background(), key, values...).Err()
}

// LRange 获取 Redis 的 list 中指定范围的元素
func LRange(key string, start, stop int64) ([]string, error) {
	// 判断key是否存在
	if !ExistKey(key) {
		return []string{}, gerr.RedisKeyNotFoundException.New(key)
	}
	return Cli.LRange(context.Background(), key, start, stop).Result()
}

// LIndex 获取 Redis 的 list 中指定值的元素
func LIndex(key string, index int64) (string, error) {
	// 判断key是否存在
	if !ExistKey(key) {
		return "", gerr.RedisKeyNotFoundException.New(key)
	}
	return Cli.LIndex(context.Background(), key, index).Result()
}

// LPop 获取 Redis 的 list 删除左边的元素
func LPop(key string) (string, error) {
	// 判断key是否存在
	if !ExistKey(key) {
		return "", gerr.RedisKeyNotFoundException.New(key)
	}
	return Cli.LPop(context.Background(), key).Result()
}

// RPop 获取 Redis 的 list 删除右边的元素
func RPop(key string) (string, error) {
	// 判断key是否存在
	if !ExistKey(key) {
		return "", gerr.RedisKeyNotFoundException.New(key)
	}
	return Cli.RPop(context.Background(), key).Result()
}

//  /////////////////////////// 无序集合 ///////////////////////////

// SAdd 无序集合添加元素
func SAdd(key string, members ...interface{}) error {
	return Cli.SAdd(context.Background(), key, members...).Err()
}

// SMembers 无序集合获取所有成员
func SMembers(key string) ([]string, error) {
	// 判断key是否存在
	if !ExistKey(key) {
		return []string{}, gerr.RedisKeyNotFoundException.New(key)
	}
	return Cli.SMembers(context.Background(), key).Result()
}

// SIsMember 无序集合判断是否属于该集合
func SIsMember(key string, member interface{}) (bool, error) {
	// 判断key是否存在
	if !ExistKey(key) {
		return false, gerr.RedisKeyNotFoundException.New(key)
	}
	return Cli.SIsMember(context.Background(), key, member).Result()
}

// SRem 无序集合删除成员
func SRem(key string, members ...interface{}) error {
	// 判断key是否存在
	if !ExistKey(key) {
		return gerr.RedisKeyNotFoundException.New(key)
	}
	return Cli.SRem(context.Background(), key, members...).Err()
}

//  /////////////////////////// 有序集合 ///////////////////////////

// ZAdd 添加元素
func ZAdd(key string, members ...*redis.Z) error {
	return Cli.ZAdd(context.Background(), key, members...).Err()
}

// ZRange 获取范围内元素
func ZRange(key string, start, stop int64) ([]string, error) {
	// 判断key是否存在
	if !ExistKey(key) {
		return []string{}, gerr.RedisKeyNotFoundException.New(key)
	}
	return Cli.ZRange(context.Background(), key, start, stop).Result()
}

// ZRank 有序集合中指定成员的序号
func ZRank(key string, member string) (int64, error) {
	// 判断key是否存在
	if !ExistKey(key) {
		return -1, gerr.RedisKeyNotFoundException.New(key)
	}
	return Cli.ZRank(context.Background(), key, member).Result()
}

// ZRem 删除
func ZRem(key string, members ...interface{}) error {
	// 判断key是否存在
	if !ExistKey(key) {
		return gerr.RedisKeyNotFoundException.New(key)
	}
	return Cli.ZRem(context.Background(), key, members...).Err()
}

// Rename 更改key的名字
func Rename(oldKey, newKey string) (bool, error) {
	// 判断key是否存在
	if !ExistKey(oldKey) {
		return false, gerr.RedisKeyNotFoundException.New(oldKey)
	}
	return Cli.RenameNX(context.Background(), oldKey, newKey).Result()
}

// ExistKey key是否存在
func ExistKey(key string) bool {
	ok, _ := Cli.Exists(context.Background(), key).Result()
	return ok == 1
}

// Delete 删除key
func Delete(keys ...string) error {
	return Cli.Del(context.Background(), keys...).Err()
}

// ExpireTime 设置key过期时间
func ExpireTime(key string, t time.Duration) error {
	// 判断key是否存在
	if !ExistKey(key) {
		return gerr.RedisKeyNotFoundException.New(key)
	}
	return Cli.Expire(context.Background(), key, t).Err()
}

// GetExpire 获得 key 的过期时间
func GetExpire(key string) (time.Duration, error) {
	if !ExistKey(key) {
		return -1, gerr.RedisKeyNotFoundException.New(key)
	}
	return Cli.TTL(context.Background(), key).Result()
}

// RemoveExpire 删除key过期时间
func RemoveExpire(key string) error {
	// 判断key是否存在
	if !ExistKey(key) {
		return gerr.RedisKeyNotFoundException.New(key)
	}
	return Cli.Persist(context.Background(), key).Err()
}

// Subscribe 订阅消息
func Subscribe(channel string, cb func(string)) *redis.PubSub {
	pubsub := Cli.Subscribe(context.Background(), channel)
	// 处理接收到的消息
	for msg := range pubsub.Channel() {
		cb(msg.Payload)
	}
	return pubsub
}

// Publish 发布消息
func Publish(channel, message string) {
	Cli.Publish(context.Background(), channel, message)
}

func Eval(script string, keys []string, args ...interface{}) error {
	return Cli.Eval(context.Background(), script, keys, args...).Err()
}
