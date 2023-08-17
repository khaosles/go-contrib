package model

import (
	"time"
)

/*
   @File: redis.go
   @Author: khaosles
   @Time: 2023/4/11 21:11
   @Desc:
*/

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
