/**
 * Cthul System
 *
 * Copyright (C) 2024 Linus Ilian Moser <linus.moser@megakuul.ch>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package etcdv3

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

// EtcdClient provides a Client implementation for etcdv3.
type EtcdClient struct {
	config clientv3.Config
	client *clientv3.Client
}

type EtcdClientOption func(*EtcdClient)

// NewEtcdClient creates a new etcdv3 client.
func NewEtcdClient(endpoints []string, opts ...EtcdClientOption) *EtcdClient {
	etcdClient := &EtcdClient{
		config: clientv3.Config{
			Endpoints: endpoints,
			TLS: &tls.Config{},
			DialTimeout: time.Second * 2,
			Logger: zap.NewNop(),
		},
		client: nil,
	}

	for _, opt := range opts {
		opt(etcdClient)
	}
	
	return etcdClient
}

// WithAuth adds credentials for basic authentication to the client.
func WithAuth(username, password string) EtcdClientOption {
	return func (c *EtcdClient) {
		c.config.Username = username
		c.config.Password = password
	}
}

// WithDialTimeout defines a custom dial timeout.
func WithDialTimeout(timeout time.Duration) EtcdClientOption {
	return func (c *EtcdClient) {
		c.config.DialTimeout = timeout
	}
}

// WithSkipVerify skips tls public cert verification.
func WithSkipVerify(skip bool) EtcdClientOption {
	return func (c *EtcdClient) {
		c.config.TLS.InsecureSkipVerify = skip
	}
}


// initClient creates the underlying etcdv3 client if not already initialized.
func (c *EtcdClient) initClient() error {
	if c.client != nil {
		return nil
	}
	client, err := clientv3.New(c.config)
	if err!=nil {
		return err
	}
	c.client = client
	return nil
}

// CheckEndpointHealth reads the status for every database client endpoint.
// If an error is found on one of the checked endpoints, it is returned immediately.
func (c *EtcdClient) CheckEndpointHealth(ctx context.Context) error {
	if err := c.initClient(); err!=nil {
		return err
	}
	for _, endpoint := range c.config.Endpoints {
		res, err := c.client.Maintenance.Status(ctx, endpoint)
		if err!=nil {
			return err
		}
		for _, err := range res.Errors {
			return fmt.Errorf(err)
		}
	}
	return nil
}

// Get returns a single key. If the key is empty or not existent, an empty string is returned.
func (c *EtcdClient) Get(ctx context.Context, key string) (string, error) {
	if err := c.initClient(); err!=nil {
		return "", err
	}
	res, err := c.client.KV.Get(ctx, key)
	if err!=nil {
		return "", err
	}

	if len(res.Kvs) < 1 {
		return "", nil
	}

	return string(res.Kvs[1].Value), nil
}

// GetRange returns a kv map with all keys that match the prefix.
func (c *EtcdClient) GetRange(ctx context.Context, prefix string) (map[string]string, error) {
	if err := c.initClient(); err!=nil {
		return nil, err
	}
	res, err := c.client.KV.Get(ctx, prefix, clientv3.WithPrefix())
	if err!=nil {
		return nil, err
	}

	kvMap := map[string]string{}
	for _, kv := range res.Kvs {
		kvMap[string(kv.Key)] = string(kv.Value)
	}

	return kvMap, nil
}


// Set upserts a kv to the database. If ttl is set to 0 the kv never expires.
func (c *EtcdClient) Set(ctx context.Context, key, value string, ttl int64) error {
	if err := c.initClient(); err!=nil {
		return err
	}
	if ttl!=0 {
		// utilizing lease checking, keep-alive, revokation, etc. hardly overcomplicates this usecase
		// therefore we just grant a new lease every time the key is set.
		lease, err := c.client.Lease.Grant(ctx, ttl)
		if err!=nil {
			return err
		}
		_, err = c.client.KV.Put(ctx, key, value, clientv3.WithLease(lease.ID))
		if err!=nil {
			return err
		}
	} else {
		_, err := c.client.KV.Put(ctx, key, value)
		if err!=nil {
			return err
		}
	}

	return nil
}

// Delete deletes one specific kv by key.
func (c *EtcdClient) Delete(ctx context.Context, key string) error {
	if err := c.initClient(); err!=nil {
		return err
	}
	_, err := c.client.KV.Delete(ctx, key)
	if err!=nil {
		return err
	}

	return nil
}

// DeleteRange deletes all kvs that match the prefix.
func (c *EtcdClient) DeleteRange(ctx context.Context, prefix string) error {
	if err := c.initClient(); err!=nil {
		return err
	}
	_, err := c.client.KV.Delete(ctx, prefix, clientv3.WithPrefix())
	if err!=nil {
		return err
	}

	return nil
}

// Watch starts a blocking listener that reacts to changes on the specified key.
// The event function is triggered on every event, containing the key, value and an error on failure.
// Stop the watcher by cancelling the context.
func (c *EtcdClient) Watch(ctx context.Context, key string, event func(string, string, error)) error {
	if err := c.initClient(); err!=nil {
		return err
	}
	return c.startWatchCycle(c.client.Watcher.Watch(ctx, key), event)
}

// WatchRange starts a blocking listener that reacts to changes on keys in the specified prefix.
// The event function is triggered on every event, containing the key, value and an error on failure.
// Stop the watcher by cancelling the context.
func (c *EtcdClient) WatchRange(ctx context.Context, prefix string, event func(string, string, error)) error {
	if err := c.initClient(); err!=nil {
		return err
	}
	return c.startWatchCycle(c.client.Watcher.Watch(ctx, prefix, clientv3.WithPrefix()), event)
}

// startWatchCycle implements the actual watch cycle.
func (c *EtcdClient) startWatchCycle(watchChan clientv3.WatchChan, eventFunc func(string, string, error)) error {
	for {
		select {
		case event, ok := <- watchChan:
			if !ok || event.Canceled {
				return nil
			}

			for _,update := range event.Events {
				eventFunc(string(update.Kv.Key), string(update.Kv.Value), event.Err())
			}
		}
	}
}


// Terminate cleans up the underlying etcd client, terminating all client connections.
// Connections are terminated forcefully, the context is only provided to match the cthul terminate pattern.
func (c *EtcdClient) Terminate(ctx context.Context) error {
	// c.rootCtxCancel()
	if c.client!=nil {
		return c.client.Close()
	}
	return nil
}
