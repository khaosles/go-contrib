package resolver

import (
	"context"
	"sync"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	gresolver "google.golang.org/grpc/resolver"
)

/*
   @File: resolver.go
   @Author: khaosles
   @Time: 2023/7/6 17:04
   @Desc:
*/

type Option int

const (
	Add Option = iota
	Delete
)

type Update struct {
	Op   Option
	Key  string
	Addr string
}

type WatchChan chan []*Update

type Builder struct {
	client *clientv3.Client
}

func NewBuilder(client *clientv3.Client) *Builder {
	return &Builder{client: client}
}

func (b Builder) Build(target gresolver.Target, cc gresolver.ClientConn, opts gresolver.BuildOptions) (gresolver.Resolver, error) {
	r := NewResolver(b.client, target.Endpoint(), cc)
	return r, nil
}

func (b Builder) Scheme() string {
	return "gogo"
}

type Resolver struct {
	client *clientv3.Client
	target string
	cc     gresolver.ClientConn
	wch    WatchChan
	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
}

func NewResolver(client *clientv3.Client, target string, cc gresolver.ClientConn) *Resolver {
	var resolver Resolver
	resolver.cc = cc
	resolver.client = client
	resolver.target = target
	resolver.init()
	return &resolver
}

func (r *Resolver) init() {
	r.ctx, r.cancel = context.WithCancel(context.TODO())
	resp, err := r.client.Get(r.ctx, r.target, clientv3.WithPrefix(), clientv3.WithSerializable())
	lg := r.client.GetLogger()
	if err != nil {
		lg.Warn("unmarshal get failed", zap.Error(err))
	}
	initUpdates := make([]*Update, 0, len(resp.Kvs))
	for _, kv := range resp.Kvs {
		// var addr string
		// if err := json.Unmarshal(kv.Value, &addr); err != nil {
		//	lg.Warn("unmarshal update failed", zap.String("key", string(kv.Key)), zap.Error(err))
		//	continue
		// }
		up := &Update{
			Op:   Add,
			Key:  string(kv.Key),
			Addr: string(kv.Value),
		}
		initUpdates = append(initUpdates, up)
	}
	r.wch = make(chan []*Update, 1)
	if len(initUpdates) > 0 {
		r.wch <- initUpdates
	}
	go r.watch(resp.Header.Revision + 1)
	r.wg.Add(1)
	go r.update()
}

func (r *Resolver) watch(rev int64) {
	defer close(r.wch)
	lg := r.client.GetLogger()
	opts := []clientv3.OpOption{clientv3.WithRev(rev), clientv3.WithPrefix()}
	wch := r.client.Watch(r.ctx, r.target, opts...)
	for {
		select {
		case <-r.ctx.Done():
			return
		case wresp, ok := <-wch:
			if !ok {
				lg.Warn("watch closed", zap.String("target", r.target))
				return
			}
			if wresp.Err() != nil {
				lg.Warn("watch failed", zap.String("target", r.target), zap.Error(wresp.Err()))
				return
			}

			deltaUps := make([]*Update, 0, len(wresp.Events))
			for _, e := range wresp.Events {
				var addr string
				var err error
				var op Option
				switch e.Type {
				case clientv3.EventTypePut:
					addr = string(e.Kv.Value)
					// err = json.Unmarshal(e.Kv.Value, &addr)
					op = Add
					if err != nil {
						lg.Warn("unmarshal endpoint update failed", zap.String("key", string(e.Kv.Key)), zap.Error(err))
						continue
					}
				case clientv3.EventTypeDelete:
					op = Delete
				default:
					continue
				}
				up := &Update{Op: op, Key: string(e.Kv.Key), Addr: addr}
				deltaUps = append(deltaUps, up)
			}
			if len(deltaUps) > 0 {
				r.wch <- deltaUps
			}
		}
	}
}

func (r *Resolver) update() {
	defer r.wg.Done()
	allUps := make(map[string]*Update)
	for {
		select {
		case <-r.ctx.Done():
			return
		case ups, ok := <-r.wch:
			if !ok {
				return
			}
			for _, up := range ups {
				switch up.Op {
				case Add:
					allUps[up.Key] = up
				case Delete:
					delete(allUps, up.Key)
				default:
					continue
				}
			}
			addrs := convertToGRPCAddress(allUps)
			err := r.cc.UpdateState(gresolver.State{Addresses: addrs})
			if err != nil {
				return
			}
		}
	}
}

func convertToGRPCAddress(ups map[string]*Update) []gresolver.Address {
	var addrs []gresolver.Address
	for _, up := range ups {
		addr := gresolver.Address{
			Addr: up.Addr,
		}
		addrs = append(addrs, addr)
	}
	return addrs
}

func (r *Resolver) ResolveNow(gresolver.ResolveNowOptions) {}

func (r *Resolver) Close() {
	r.cancel()
	r.wg.Wait()
}
