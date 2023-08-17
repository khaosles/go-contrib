package sqlite

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/khaosles/go-contrib/core/config"
	glog "github.com/khaosles/go-contrib/core/log"
	"github.com/khaosles/go-contrib/gorm/internal"
)

/*
   @File: sqlite.go
   @Author: khaosles
   @Time: 2023/4/23 01:02
   @Desc:
*/

var DB *gorm.DB

const APP = "sqlite"

type Sqlite struct {
	internal.Datasource `yaml:",inline" mapstructure:",squash"`
}

// Dsn 基于配置文件获取 dsn
func (s *Sqlite) Dsn() string {
	return s.DbName
}

func init() {
	var err error
	var sql Sqlite
	if err = config.Configuration(APP, &sql); err != nil {
		glog.Fatal(err)
	}
	if DB, err = New(&sql); err != nil {
		glog.Fatal(err)
	}
}

func New(cfg *Sqlite) (*gorm.DB, error) {
	var err error
	var db *gorm.DB
	if db, err = gorm.Open(sqlite.Open(cfg.Dsn()), internal.Gorm.Config(
		cfg.Prefix, cfg.Singular, cfg.LogMode, cfg.LogZap),
	); err != nil {
		glog.Error("Sqlite connect failed ===> ", cfg.Dsn())
		return nil, err
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
		sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
		glog.Info("Sqlite connect succeed ===> ", cfg.Dsn())
	}
	return db, nil
}
