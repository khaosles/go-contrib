package consumer

import (
	"context"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/rlog"
	"github.com/khaosles/go-contrib/core/config"
	glog "github.com/khaosles/go-contrib/core/log"
	"github.com/khaosles/go-contrib/rocket"

	"github.com/khaosles/giz/g"
)

/*
   @File: consumer.go
   @Author: khaosles
   @Time: 2023/8/17 11:34
   @Desc:
*/

var pushConsumer rocketmq.PushConsumer

type Consumer struct {
	rocket.Rocketmq `yaml:",inline" mapstructure:",squash"`
}

func init() {
	var c Consumer
	// 解析参数
	if err := config.Configuration(rocket.APP, &c); err != nil {
		glog.Fatal(err)
	}
	rlog.SetLogLevel(c.LogLevel)
	// 启动实例
	if err := c.run(); err != nil {
		glog.Fatal(err)
	}
	glog.Info("Push consumer connect succeed")
}

func (c Consumer) run() error {
	var err error
	// push
	pushConsumer, err = rocketmq.NewPushConsumer(
		consumer.WithNameServer(c.NameServer),
		consumer.WithGroupName(c.GroupName),
		consumer.WithRetry(c.Retry),
		consumer.WithCredentials(primitive.Credentials{
			AccessKey: c.AccessKey,
			SecretKey: c.SecretKey,
		}),
	)
	if err != nil {
		return err
	}
	return nil
}

func Subscribe(topic string, selector consumer.MessageSelector,
	cb func(context.Context, ...*primitive.MessageExt) (consumer.ConsumeResult, error)) error {
	// 监听退出
	go g.Exit(func() {
		// 停止消费
		_ = pushConsumer.Shutdown()
		// 取消订阅
		_ = pushConsumer.Unsubscribe(topic)
	})
	// 订阅主题
	if err := pushConsumer.Subscribe(topic, selector, cb); err != nil {
		return err
	}
	// 开始消费
	if err := pushConsumer.Start(); err != nil {
		return err
	}
	return nil
}
