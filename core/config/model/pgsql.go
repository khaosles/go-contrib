package model

/*
   @File: pgsql.go
   @Author: khaosles
   @Time: 2023/4/15 22:27
   @Desc:
*/

type Pgsql struct {
	Datasource `yaml:",inline" mapstructure:",squash"`
}

// Dsn 基于配置文件获取 dsn
func (p *Pgsql) Dsn() string {
	return "host=" + p.Host + " user=" + p.Username + " password=" + p.Password + " dbname=" + p.DbName + " port=" + p.Port + " " + p.Config
}

func (p *Pgsql) DsnHide() string {
	return "host=" + p.Host + " user=" + p.Username + " password=" + "******" + " dbname=" + p.DbName + " port=" + p.Port + " " + p.Config
}

// LinkDsn 根据 dbname 生成 dsn
func (p *Pgsql) LinkDsn(dbname string) string {
	return "host=" + p.Host + " user=" + p.Username + " password=" + p.Password + " dbname=" + dbname + " port=" + p.Port + " " + p.Config
}

func (m *Pgsql) GetLogMode() string {
	return m.LogMode
}
