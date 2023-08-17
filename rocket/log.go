package rocket

import (
	"fmt"

	"github.com/apache/rocketmq-client-go/v2/rlog"

	glog "github.com/khaosles/go-contrib/core/log"
)

/*
   @File: log.go
   @Author: khaosles
   @Time: 2023/7/2 09:33
   @Desc:
*/

var logger = new(mqLogger)

type mqLogger struct {
	rlog.Logger
}

func (l mqLogger) Debug(msg string, fields map[string]interface{}) {
	glog.Debug(fmt.Sprintf("%s -> %+v", msg, fields))
}

func (l mqLogger) Info(msg string, fields map[string]interface{}) {
	glog.Info(fmt.Sprintf("%s -> %+v", msg, fields))
}

func (l mqLogger) Warning(msg string, fields map[string]interface{}) {
	glog.Warn(fmt.Sprintf("%s -> %+v", msg, fields))
}

func (l mqLogger) Error(msg string, fields map[string]interface{}) {
	glog.Error(fmt.Sprintf("%s -> %+v", msg, fields))
}

func (l mqLogger) Fatal(msg string, fields map[string]interface{}) {
	glog.Panic(fmt.Sprintf("%s -> %+v", msg, fields))
}

func (l mqLogger) Level(msglevel string) {
}

func (l mqLogger) OutputPath(path string) (err error) {
	return nil
}
