package etcd

import (
	"context"

	"go.etcd.io/etcd/client/v3/concurrency"
)

/*
   @File: distributed.go
   @Author: khaosles
   @Time: 2023/8/6 11:28
   @Desc:
*/

// ExecOnceTask 基于分布式的唯一执行任务
func ExecOnceTask(key string, ttl int, task func()) error {
	ctx := context.Background()
	// 建立session
	session, err := concurrency.NewSession(
		Client,
		concurrency.WithTTL(ttl),
		concurrency.WithContext(ctx),
	)
	if err != nil {
		return err
	}
	mutex := concurrency.NewMutex(session, key)
	// 尝试获取锁
	err = mutex.TryLock(ctx)
	if err != nil {
		return err
	}
	// 执行操作
	task()
	// 释放锁
	_ = mutex.Unlock(ctx)
	return nil
}

// ExecAllTask 基于分布式的全部执行任务
func ExecAllTask(key string, ttl int, task func()) error {
	ctx := context.Background()
	// 建立session
	session, err := concurrency.NewSession(
		Client,
		concurrency.WithTTL(ttl),
		concurrency.WithContext(ctx),
	)
	if err != nil {
		return err
	}
	mutex := concurrency.NewMutex(session, key)
	// 尝试获取锁
	err = mutex.Lock(ctx)
	if err != nil {
		return err
	}
	// 执行操作
	task()
	// 释放锁
	_ = mutex.Unlock(ctx)
	return nil
}
