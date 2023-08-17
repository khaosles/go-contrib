package etcd

import (
	"context"
	"fmt"
	"time"

	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"

	"github.com/khaosles/go-contrib/core/config"
	glog "github.com/khaosles/go-contrib/core/log"
)

/*
   @File: etcd.go
   @Author: khaosles
   @Time: 2023/7/1 23:35
   @Desc:
*/

var Client *clientv3.Client

const APP = "etcd"

type Etcd struct {
	Nodes    []string `mapstructure:"nodes" default:"" yaml:"nodes" json:"nodes"`
	Username string   `mapstructure:"username" default:"" yaml:"username" json:"username"`
	Password string   `mapstructure:"password" default:"" yaml:"password" json:"password"`
	Timeout  int      `mapstructure:"timeout" default:"10" yaml:"timeout" json:"timeout"`
	TTL      int64    `mapstructure:"ttl" default:"10" yaml:"ttl" json:"ttl"`
}

func init() {
	var e *Etcd
	var err error
	if err = config.Configuration(APP, e); err != nil {
		glog.Fatal(err)
	}
	if Client, err = New(e); err != nil {
		glog.Fatal(err)
	}
}

func New(e *Etcd) (*clientv3.Client, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   e.Nodes,
		DialTimeout: time.Duration(e.Timeout) * time.Second,
		Username:    e.Username,
		Password:    e.Password,
	})
	if err != nil {
		glog.Error("etcd connect failed")
		return nil, err
	}
	glog.Info("etcd connect succeed")
	return cli, nil
}

func Register(serverName, addr string, ttl int64) error {
	em, err := endpoints.NewManager(Client, serverName)
	if err != nil {
		return err
	}
	lease, _ := Client.Grant(context.TODO(), ttl)
	err = em.AddEndpoint(
		context.TODO(),
		fmt.Sprintf("%s/%s", serverName, addr),
		endpoints.Endpoint{Addr: addr},
		clientv3.WithLease(lease.ID),
	)
	if err != nil {
		return err
	}
	alive, err := Client.KeepAlive(context.TODO(), lease.ID)
	if err != nil {
		return err
	}
	go func() {
		for {
			<-alive
		}
	}()
	return nil
}

func Unregister(serverName, addr string) error {
	em, err := endpoints.NewManager(Client, serverName)
	if err != nil {
		return err
	}
	err = em.DeleteEndpoint(context.TODO(), fmt.Sprintf("%s/%s", serverName, addr))
	if err != nil {
		return err
	}
	return nil
}

func Put(key, value string, opts ...clientv3.OpOption) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	_, err := Client.Put(ctx, key, value, opts...)
	cancel()
	if err != nil {
		return err
	}
	return nil
}

func Get(key string, opts ...clientv3.OpOption) ([]*mvccpb.KeyValue, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	resp, err := Client.Get(ctx, key, opts...)
	cancel()
	if err != nil {
		return nil, err
	}
	return resp.Kvs, nil
}

func Del(key string, opts ...clientv3.OpOption) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	_, err := Client.Delete(ctx, key, opts...)
	cancel()
	if err != nil {
		return err
	}
	return nil
}

func PutWithLease(key, val string, ttl int64, opts ...clientv3.OpOption) error {
	// 设置租约
	lease, err := Client.Grant(context.Background(), ttl)
	if err != nil {
		return err
	}
	opts = append(opts, clientv3.WithLease(lease.ID))
	// 添加带租约的key
	_, err = Client.Put(context.Background(), key, val, opts...)
	if err != nil {
		return err
	}
	// 续约
	alive, err := Client.KeepAlive(context.Background(), lease.ID)
	if err != nil {
		return err
	}
	go func() {
		for {
			<-alive
		}
	}()
	return nil
}
