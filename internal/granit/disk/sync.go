/**
 * Cthul System
 *
 * Copyright (C) 2025 Linus Ilian Moser <linus.moser@megakuul.ch>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program. If not, see <https://www.gnu.org/licenses/>.
 */

package disk

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"cthul.io/cthul/pkg/api/granit/v1/disk"
	"google.golang.org/protobuf/proto"
)

func (o *Operator) synchronize() {
	primaryMap, primaryMapLock := map[string]bool{}, sync.RWMutex{}
	configMap, configMapLock  := map[string]string{}, sync.RWMutex{}
	syncChan := make(chan string)

	o.operationWg.Add(1)
	go func() {
		defer o.operationWg.Done()
		for {
			device := ""
			select {
			case <-o.rootCtx.Done():
				return
			case id := <-syncChan:
				device = id
			}

      ctx, cancel := context.WithTimeout(o.rootCtx, time.Duration(o.syncCycleTTL))
      defer cancel()

      configMapLock.RLock()
			config, ok := configMap[device]
      configMapLock.RUnlock()
			if !ok {
				o.logger.Debug("aborting triggered sync: device is not configured for this node",
					"device", device, "node", o.nodeId,
				)
				continue
			}
      primaryMapLock.RLock()
			primary, ok := primaryMap[device]
      primaryMapLock.RUnlock()
			if !ok {
				o.logger.Debug("aborting triggered sync: device is not configured for this node",
					"device", device, "node", o.nodeId,
				)
        continue
			}

			if primary {
				err := o.applyConfig(ctx, config, true)
				if err != nil {
					o.logger.Error(err.Error(), "device", device, "node", o.nodeId)
					continue
				}
        oldNode, err := o.client.Set(ctx, fmt.Sprintf("/GRANIT/DISK/NODE/%s", device), o.nodeId, 0)
				if err != nil {
					o.logger.Error(err.Error(), "device", device, "node", o.nodeId)
					continue
				}
        if oldNode != o.nodeId {
          o.logger.Info("successfully attached device to a new cluster node",
            "device", device, "node", o.nodeId,
          )
        }
			} else {
				err := o.applyConfig(ctx, config, false)
				if err != nil {
					o.logger.Error(err.Error(), "device", device, "node", o.nodeId)
					continue
				}
			}
		}
	}()

	o.syncer.Add("/GRANIT/DISK/CLUSTERNODES/", o.updateCycleTTL, func(ctx context.Context, k, rawCluster string) error {
		id := strings.TrimPrefix(k, "/GRANIT/DISK/CLUSTERNODES/")
		configKey := fmt.Sprintf("/GRANIT/DISK/CONFIG/%s", id)

		cluster := &disk.DiskCluster{}
		err := proto.Unmarshal([]byte(rawCluster), cluster)
		if err != nil {
			return err
		}
		if _, ok := cluster.Nodes[o.nodeId]; ok {
			o.syncer.Add(configKey, o.syncCycleTTL, func(ctx context.Context, k, v string) error {
        configMapLock.Lock()
        configMap[id] = v
        configMapLock.Unlock()
        syncChan <- id
				return nil
			})
		} else {
			o.syncer.Remove(configKey, true)

      configMapLock.Lock()
      delete(configMap, id)
      configMapLock.Unlock()

      primaryMapLock.Lock()
      delete(primaryMap, id)
      primaryMapLock.Unlock()

			o.pruneDevice(ctx, id)
		}
		return nil
	})

	o.syncer.Add("/GRANIT/DISK/REQNODE/", o.updateCycleTTL, func(ctx context.Context, k, reqnode string) error {
		id := strings.TrimPrefix(k, "/GRANIT/DISK/REQNODE/")
    primaryMapLock.Lock()
    primaryMap[id] = reqnode == o.nodeId
    primaryMapLock.Unlock()
    syncChan <- id
		return nil
	})
}

func (o *Operator) applyConfig(ctx context.Context, rawConfig string, primary bool) error {
	config := &disk.DiskConfig{}
	err := proto.Unmarshal([]byte(rawConfig), config)
	if err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}
	// drbdadm up r0
	// drbdadm primary r0 o.nodeid
	// drbdadm secondary r0 o.nodeId
}

func (o *Operator) pruneDevice(ctx context.Context, id string) {

	// umount /dev/drbdxy
	// umount /dev/loopdev
	// rm -rf /device.img
}
