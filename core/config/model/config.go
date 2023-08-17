package model

/*
   @File: config.go
   @Author: khaosles
   @Time: 2023/4/10 23:27
   @Desc:
*/

type Config struct {
	System  *System  `mapstructure:"system" json:"system" yaml:"system"`
	Server  *Server  `mapstructure:"server" json:"server" yaml:"server"`
	Pgsql   *Pgsql   `mapstructure:"pgsql" json:"pgsql" yaml:"pgsql"`
	Mysql   *Mysql   `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	Sqlite  *Sqlite  `mapstructure:"sqlite" json:"sqlite" yaml:"sqlite"`
	Logging *Logging `mapstructure:"logging" json:"logging" yaml:"logging"`
	Redis   *Redis   `mapstructure:"redis" json:"redis" yaml:"redis"`
	Etcd    *Etcd    `mapstructure:"etcd" json:"etcd" yaml:"etcd"`
	// Rpc      *Rpc      `mapstructure:"rpc" json:"rpc" yaml:"rpc"`
	Rocketmq *Rocketmq `mapstructure:"rocketmq" json:"rocketmq" yaml:"rocketmq"`
}
