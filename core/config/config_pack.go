//go:build pack

package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/khaosles/giz/fileutil"
	"github.com/spf13/viper"

	"github.com/khaosles/go-contrib/core/config/model"
)

/*
  @File: cfg_configure.go
  @Author: khaosles
  @Time: 2023/2/19 11:04
  @Desc:
*/

var GCfg model.Config  // 配置类
var Viper *viper.Viper // 配置源

func init() {
	cfg := "config.yaml"
	rootPath, _ := os.Executable()
	cfg = filepath.Join(fileutil.Dirname(rootPath), cfg)

	// 配置文件不存在
	if !fileutil.Exist(cfg) {
		log.Fatal("Configure not exists ===> ", cfg)
	}

	// 创建viper
	Viper = viper.New()

	Viper.SetConfigFile(cfg)
	Viper.SetConfigType("yaml")

	// 读取配置文件
	err := Viper.ReadInConfig()
	if err != nil {
		log.Fatal("Configure reading error.")
	}
	Viper.WatchConfig()
	if err = Viper.Unmarshal(&GCfg); err != nil {
		log.Fatal("Configure parse error.")
	}

	// 识别 default 标签
	setDefaults(reflect.ValueOf(&GCfg).Elem())
	// 解析配置文件
	Viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
		if err = Viper.Unmarshal(&GCfg); err != nil {
			log.Fatal("Configure parse error.")
		}
		// 识别 default 标签
		setDefaults(reflect.ValueOf(&GCfg).Elem())
	})
}

func Configuration(name string, conf any) error {
	// 解析参数
	err := Viper.UnmarshalKey(name, conf)
	if err != nil {
		return err
	}
	// 如果参数为空，则设置默认值
	setDefaults(reflect.ValueOf(conf).Elem())
	return nil
}

// setDefaults 设置结构体默认值
func setDefaults(v reflect.Value) {
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		defaultTag := v.Type().Field(i).Tag.Get("default")
		if defaultTag != "" && field.Interface() == reflect.Zero(field.Type()).Interface() {
			var defaultValue any
			switch field.Type().Kind() {
			case reflect.Int:
				defaultValue, _ = strconv.Atoi(defaultTag)
			case reflect.Int32:
				defaultValue, _ = strconv.ParseInt(defaultTag, 8, 32)
			case reflect.Int64:
				if field.Type().Name() == reflect.TypeOf(int64(0)).Name() {
					defaultValue, _ = strconv.ParseInt(defaultTag, 0, 0)
				} else {
					defaultValue, _ = time.ParseDuration(defaultTag)
				}
			case reflect.Float32:
				defaultValue, _ = strconv.ParseFloat(defaultTag, 32)
			case reflect.Float64:
				defaultValue, _ = strconv.ParseFloat(defaultTag, 64)
			case reflect.Bool:
				defaultValue, _ = strconv.ParseBool(defaultTag)
			default:
				defaultValue = defaultTag
			}
			field.Set(reflect.ValueOf(defaultValue).Convert(field.Type()))
		}
		if field.Kind() == reflect.Struct {
			setDefaults(field)
		}
	}
}
