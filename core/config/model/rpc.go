package model

/*
   @File: rpc.go
   @Author: khaosles
   @Time: 2023/6/29 18:03
   @Desc:
*/

type Rpc struct {
	Host string `mapstructure:"host" default:"0.0.0.0" yaml:"host" json:"host"`
	Port string `mapstructure:"port" default:"50000" yaml:"port" json:"port"`
}
