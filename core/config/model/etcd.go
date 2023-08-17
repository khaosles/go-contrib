package model

/*
   @File: etcd.go
   @Author: khaosles
   @Time: 2023/6/29 00:16
   @Desc:
*/

type Etcd struct {
	Nodes    []string `mapstructure:"nodes" default:"" yaml:"nodes" json:"nodes"`
	Username string   `mapstructure:"username" default:"" yaml:"username" json:"username"`
	Password string   `mapstructure:"password" default:"" yaml:"password" json:"password"`
	Timeout  int      `mapstructure:"timeout" default:"10" yaml:"timeout" json:"timeout"`
	TTl      int64    `mapstructure:"ttl" default:"10" yaml:"ttl" json:"ttl"`
}
