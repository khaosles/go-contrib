package model

/*
   @File: rocketmq.go
   @Author: khaosles
   @Time: 2023/7/1 09:56
   @Desc:
*/

type Rocketmq struct {
	NameServer    []string `mapstructure:"name-server" default:"" yaml:"name-server" json:"nameServer""`
	AccessKey     string   `mapstructure:"access-key" default:"" yaml:"access-key" json:"accessKey"`
	SecretKey     string   `mapstructure:"secret-key" default:"" yaml:"secret-key" json:"secretKey"`
	SecurityToken string   `mapstructure:"security-token" default:"" yaml:"security-token" json:"securityToken"`
	Topic         string   `mapstructure:"topic" default:"" yaml:"topic" json:"topic"`
	LogLevel      string   `mapstructure:"log-level" default:"" yaml:"log-level" json:"logLevel"`
	Retry         int      `mapstructure:"retry" default:"" yaml:"lretry" json:"retry"`
}
