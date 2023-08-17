package model

/*
   @File: system.go
   @Author: khaosles
   @Time: 2023/4/17 12:27
   @Desc:
*/

type System struct {
	Proxy string `mapstructure:"proxy" yaml:"proxy" json:"proxy" default:""`
}
