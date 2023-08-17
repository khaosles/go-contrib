//go:build assest

package config

import (
	"bytes"
	"log"
	"os"
	"reflect"
	"strconv"
	"time"

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
	// 创建viper
	Viper = viper.New()
	// 从打包后的文件中读取配置
	bytesContent, err := packed.Asset("manifest/config/config.yaml")
	if err != nil {
		panic("Asset() can not found setting file")
	}
	// 设置要读取的文件类型
	Viper.SetConfigType("yaml")
	// 读取
	err = Viper.ReadConfig(bytes.NewBuffer(bytesContent))
	if err != nil {
		log.Println(err)
	}

	err = Viper.Unmarshal(&GCfg)
	if err != nil {
		log.Fatal("Configure parse error.")
	}
	// 识别 default 标签
	setDefaults(reflect.ValueOf(&GCfg).Elem())
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

func exist(path string) bool {
	// path stat
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return false
}
