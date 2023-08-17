package pgsql

import (
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/khaosles/go-contrib/core/config"
	glog "github.com/khaosles/go-contrib/core/log"
	"github.com/khaosles/go-contrib/gorm/internal"
)

/*
   @File: pgsql.go
   @Author: khaosles
   @Time: 2023/4/22 23:39
   @Desc:
*/

var DB *gorm.DB

const APP = "pgsql"

type Pgsql struct {
	internal.Datasource `yaml:",inline" mapstructure:",squash"`
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

func init() {
	var err error
	var psql Pgsql
	if err = config.Configuration(APP, &psql); err != nil {
		glog.Fatal(err)
	}
	if DB, err = New(&psql); err != nil {
		glog.Fatal(err)
	}
}

func New(psql *Pgsql) (*gorm.DB, error) {
	var err error
	var pdb *gorm.DB
	pgsqlConfig := postgres.Config{
		DSN:                  psql.Dsn(), // DSN data source name
		PreferSimpleProtocol: false,
	}
	if pdb, err = gorm.Open(postgres.New(pgsqlConfig), internal.Gorm.Config(psql.Prefix, psql.Singular, psql.LogMode, psql.LogZap)); err != nil {
		glog.Error("Pgsql connect failed ===> ", psql.DsnHide())
		return nil, err
	} else {
		sqlDB, _ := pdb.DB()
		sqlDB.SetMaxIdleConns(psql.MaxIdleConns)
		sqlDB.SetMaxOpenConns(psql.MaxOpenConns)
		sqlDB.SetConnMaxLifetime(time.Minute * time.Duration(psql.MaxLifeTime))
		glog.Info("Pgsql connect succeed ===> ", psql.DsnHide())
	}
	return pdb, nil
}
