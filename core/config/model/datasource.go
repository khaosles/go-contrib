package model

/*
   @File: datasource.go
   @Author: khaosles
   @Time: 2023/4/10 23:44
   @Desc:
*/

type Datasource struct {
	Host         string `mapstructure:"host" default:"" yaml:"host" json:"host"`
	Port         string `mapstructure:"port" default:"" yaml:"port" json:"port"`
	DbName       string `mapstructure:"db-name" default:"" yaml:"db-name" json:"dbName"`
	Username     string `mapstructure:"username" default:"" yaml:"username" json:"username"`
	Password     string `mapstructure:"password" default:"" yaml:"password" json:"password"`
	Config       string `mapstructure:"config" default:"" yaml:"config" json:"config"`
	Prefix       string `mapstructure:"prefix" default:"" yaml:"prefix" json:"prefix"`
	Engine       string `mapstructure:"engine" json:"engine" yaml:"engine" default:"InnoDB"` // 数据库引擎，默认InnoDB
	Singular     bool   `mapstructure:"singular" default:"true" yaml:"singular" json:"singular"`
	MaxIdleConns int    `mapstructure:"max-idle-conns" default:"10" yaml:"max-idle-conns" json:"maxIdleConns"`
	MaxOpenConns int    `mapstructure:"max-open-conns" default:"100" yaml:"max-open-conns" json:"maxOpenConns"`
	MaxLifeTime  int    `mapstructure:"max-life-time" default:"5" yaml:"max-life-time" json:"maxLifeTime"`
	LogMode      string `mapstructure:"log-mode" default:"debug" yaml:"log-mode" json:"logMode"`
	LogZap       bool   `mapstructure:"log-zap" default:"false" yaml:"log-zap" json:"logZap"`
}
