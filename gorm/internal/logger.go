package internal

import (
	"fmt"

	"gorm.io/gorm/logger"

	glog "github.com/khaosles/go-contrib/core/log"
)

type writer struct {
	logger.Writer
	LogZap bool
}

// NewWriter writer 构造函数
func NewWriter(logZap bool, w logger.Writer) *writer {
	return &writer{Writer: w, LogZap: logZap}
}

// Printf 格式化打印日志
func (w *writer) Printf(message string, data ...interface{}) {
	if w.LogZap {
		glog.Debug(fmt.Sprintf(message+"\n", data...))
	} else {
		w.Writer.Printf(message, data...)
	}
}
