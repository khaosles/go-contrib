package producer

import (
	"context"
	"log"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/apache/rocketmq-client-go/v2/rlog"
	"github.com/khaosles/giz/g"

	"go-contrib/core/config"
	glog "go-contrib/core/log"
	"go-contrib/rocket"
)

/*
   @File: producer.go
   @Author: khaosles
   @Time: 2023/8/17 15:29
   @Desc:
*/

var pro rocketmq.Producer

type cProducer struct {
	NameServer []string `json:"nameServer" yaml:"name-server"`
	AccessKey  string   `json:"accessKey" yaml:"access-key"`
	SecretKey  string   `json:"secretKey" yaml:"secret-key"`
	GroupName  string   `json:"groupName" yaml:"group-name"`
	Retry      int      `json:"retry" yaml:"retry"`
}

func init() {
	rlog.SetLogLevel("error")
	var c *cProducer
	// 解析参数
	if err := config.Configuration(rocket.APP, c); err != nil {
		log.Fatal(err)
	}
	// 启动实例
	if err := c.run(); err != nil {
		log.Fatal(err)
	}
	glog.Info("Producer connect succeed")
}

func (c *cProducer) run() error {
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
	g.Exit(func() { _ = pro.Shutdown() })
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
