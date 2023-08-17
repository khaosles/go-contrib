package consumer

import (
	"context"
	"log"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/rlog"

	"github.com/khaosles/giz/g"

	glog "go-contrib/core/log"

	"go-contrib/core/config"
	"go-contrib/rocket"
)

/*
   @File: consumer.go
   @Author: khaosles
   @Time: 2023/8/17 11:34
   @Desc:
*/

var pushConsumer rocketmq.PushConsumer

type cConsumer struct {
	NameServer []string `json:"nameServer" yaml:"name-server"`
	AccessKey  string   `json:"accessKey" yaml:"access-key"`
	SecretKey  string   `json:"secretKey" yaml:"secret-key"`
	GroupName  string   `json:"groupName" yaml:"group-name"`
	Retry      int      `json:"retry" yaml:"retry"`
}

func init() {
	rlog.SetLogLevel("error")
	var c *cConsumer
	// 解析参数
	if err := config.Configuration(rocket.APP, c); err != nil {
		log.Fatal(err)
	}
	// 启动实例
	if err := c.run(); err != nil {
		log.Fatal(err)
	}
	glog.Info("Push consumer connect succeed")
}

func (c *cConsumer) run() error {
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
	g.Exit(func() {
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
