package producer

import (
	"context"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/apache/rocketmq-client-go/v2/rlog"
	"github.com/khaosles/giz/g"
	"github.com/khaosles/go-contrib/core/config"
	glog "github.com/khaosles/go-contrib/core/log"
	"github.com/khaosles/go-contrib/rocket"
)

/*
   @File: producer.go
   @Author: khaosles
   @Time: 2023/8/17 15:29
   @Desc:
*/

var pro rocketmq.Producer

type Producer struct {
	rocket.Rocketmq `yaml:",inline" mapstructure:",squash"`
}

func init() {
	var c Producer
	// 解析参数
	if err := config.Configuration(rocket.APP, &c); err != nil {
		glog.Fatal(err)
	}
	rlog.SetLogLevel(c.LogLevel)
	// 启动实例
	if err := c.run(); err != nil {
		glog.Fatal(err)
	}
	glog.Info("Producer connect succeed")
}

func (c Producer) run() error {
	var err error
	// 生产者
	pro, err = rocketmq.NewProducer(
		producer.WithNameServer(c.NameServer),
		producer.WithGroupName(c.GroupName),
		producer.WithRetry(c.Retry),
		producer.WithCredentials(primitive.Credentials{
			AccessKey: c.AccessKey,
			SecretKey: c.SecretKey,
		}),
	)
	if err != nil {
		return err
	}
	// 开启消费者
	if err := pro.Start(); err != nil {
		return err
	}
	// 检测程序退出关闭消费者
	go g.Exit(func() { _ = pro.Shutdown() })
	return nil
}

// SendSync 同步消息
func SendSync(ctx context.Context, mq ...*primitive.Message) (*primitive.SendResult, error) {
	result, err := pro.SendSync(ctx, mq...)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// SendAsync 异步消息
func SendAsync(ctx context.Context, mq func(ctx context.Context, result *primitive.SendResult, err error),
	msg ...*primitive.Message) error {
	if err := pro.SendAsync(ctx, mq, msg...); err != nil {
		return err
	}
	return nil
}

// SendOneWay 发送单方面消息
func SendOneWay(ctx context.Context, mq ...*primitive.Message) error {
	return pro.SendOneWay(ctx, mq...)
}
