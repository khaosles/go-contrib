package mysql

import (
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"go-contrib/core/config"
	glog "go-contrib/core/log"
	"go-contrib/gorm/internal"
)

/*
   @File: mysql.go
   @Author: khaosles
   @Time: 2023/4/15 22:32
   @Desc:
*/

var DB *gorm.DB

const APP = "mysql"

type Mysql struct {
	internal.Datasource `yaml:",inline" mapstructure:",squash"`
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

func init() {
	var msql *Mysql
	var err error
	if err = config.Configuration(APP, msql); err != nil {
		log.Fatal(err)
	}
	if DB, err = New(msql); err != nil {
		log.Fatal(err)
	}
}

func New(msql *Mysql) (*gorm.DB, error) {
	var err error
	var mydb *gorm.DB
	mysqlConfig := mysql.Config{
		DSN:                       msql.Dsn(), // DSN data source name
		DefaultStringSize:         256,        // string 类型字段的默认长度
		SkipInitializeWithVersion: true,       // 根据版本自动配置
	}
	if mydb, err = gorm.Open(mysql.New(mysqlConfig), internal.Gorm.Config(msql.Prefix, msql.Singular, msql.LogMode, msql.LogZap)); err != nil {
		glog.Errorf("Mysql connect failed ===> ", msql.DsnHide())
		return nil, err
	} else {
		mydb.InstanceSet("gorm:table_options", "ENGINE="+msql.Engine)
		sqlDB, _ := mydb.DB()
		sqlDB.SetMaxIdleConns(msql.MaxIdleConns)
		sqlDB.SetMaxOpenConns(msql.MaxOpenConns)
		sqlDB.SetConnMaxLifetime(time.Minute * time.Duration(msql.MaxLifeTime))
		glog.Info("Mysql connect succeed ===> ", msql.DsnHide())
	}

	return mydb, nil
}
