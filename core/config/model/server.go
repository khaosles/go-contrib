package model

/*
   @File: server.go
   @Author: khaosles
   @Time: 2023/4/10 23:28
   @Desc:
*/

type Server struct {
	Host string `mapstructure:"host" default:"0.0.0.0" yaml:"host" json:"host"`
	Port string `mapstructure:"port" default:"8000" yaml:"port" json:"port"`
	Mode string `mapstructure:"mode" default:"debug" yaml:"mode" json:"mode"`
}
