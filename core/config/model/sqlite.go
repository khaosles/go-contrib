package model

/*
   @File: sqlite.go
   @Author: khaosles
   @Time: 2023/4/23 00:58
   @Desc:
*/

type Sqlite struct {
	Datasource `yaml:",inline" mapstructure:",squash"`
}

// Dsn 基于配置文件获取 dsn
func (s *Sqlite) Dsn() string {
	return s.DbName
}
