package model

import "time"

/*
   @File: logging.go
   @Author: khaosles
   @Time: 2023/4/11 21:35
   @Desc:
*/

type Logging struct {
	LevelConsole string        `mapstructure:"level-console" default:"debug" yaml:"level-console" json:"level-console"`
	LevelFile    string        `mapstructure:"level-file" default:"info" yaml:"level-file" json:"level-file"`
	Prefix       string        `mapstructure:"prefix" default:"" yaml:"prefix" json:"prefix"`
	Path         string        `mapstructure:"path" default:"" yaml:"path" json:"path"`
	MaxHistory   time.Duration `mapstructure:"max-history" default:"7" yaml:"max-history" json:"maxHistory"`
	LogInConsole bool          `mapstructure:"log-in-console" default:"true" yaml:"log-in-console" json:"logInConsole"`
	LogInFile    bool          `mapstructure:"log-in-file" default:"false" yaml:"log-in-file" json:"logInFile"`
	ShowLine     bool          `mapstructure:"show-line" default:"false" yaml:"show-line" json:"show-line"`
}
