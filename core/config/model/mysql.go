package model

/*
   @File: mysql.go
   @Author: khaosles
   @Time: 2023/4/15 22:29
   @Desc:
*/

type Mysql struct {
	Datasource `yaml:",inline" mapstructure:",squash"`
}

func (m *Mysql) Dsn() string {
	return m.Username + ":" + m.Password + "@tcp(" + m.Host + ":" + m.Port + ")/" + m.DbName + "?" + m.Config
}

func (m *Mysql) DsnHide() string {
	return m.Username + ":" + "******" + "@tcp(" + m.Host + ":" + m.Port + ")/" + m.DbName + "?" + m.Config
}

func (m *Mysql) GetLogMode() string {
	return m.LogMode
}
