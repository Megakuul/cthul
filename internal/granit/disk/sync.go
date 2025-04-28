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

	"cthul.io/cthul/pkg/api/granit/v1/disk"
	"google.golang.org/protobuf/proto"
)

func (o *Operator) synchronize() {
	primaryChan := make(chan bool)
	configChan := make(chan string)

	o.operationWg.Add(1)
	go func() {
		defer o.operationWg.Done()
		primary := false
		config := ""
		for {
			select {
			case <-o.rootCtx.Done():
				return
			case primary = <-primaryChan:
			case config = <-configChan:
			}

      if config == "" {
			  o.pruneDevice(ctx, id)
        continue 
      }

			if primary {
        err := o.applyConfig(config, true)
				if err != nil {
          // logbliblab
          continue
				}
				_, err = o.client.Set(ctx, fmt.Sprintf("/GRANIT/DISK/NODE/%s", id), o.nodeId, 0)
				if err != nil {
          // logbliblab
          continue
				}
			} else {
        err := o.applyConfig(config, false)
				if err != nil {
          // logbliblab
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
        configChan <- v
        return nil
			})
		} else {
			o.syncer.Remove(configKey, true)
      configChan <- ""
		}
		return nil
	})

	o.syncer.Add("/GRANIT/DISK/REQNODE/", o.updateCycleTTL, func(ctx context.Context, k, reqnode string) error {
    primaryChan <- reqnode == o.nodeId
		return nil
	})
}

func (o *Operator) applyPrimary(rawConfig string) error {
	// drbdadm up r0
	// drbdadm primary r0 o.nodeid
}

func (o *Operator) applySecondary(rawConfig string) error {
	// drbdadm secondary r0 o.nodeId
}

func (o *Operator) pruneDevice(ctx context.Context, id, node string) {

	// umount /dev/drbdxy
	// umount /dev/loopdev
	// rm -rf /device.img
}

func (o *Operator) applyConfig(rawConfig string, primary bool) error {
	config := &disk.DiskConfig{}
	err := proto.Unmarshal([]byte(rawConfig), config)
	if err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	cleanPath := filepath.Join(o.runRoot, config.Path)
	if !strings.HasPrefix(cleanPath, o.runRoot) {
		return fmt.Errorf("device socket path is not allowed to escape the run root ('%s' => '%s')", o.runRoot, cleanPath)
	}
	err = os.MkdirAll(filepath.Dir(cleanPath), 0600)
	if err != nil {
		return err
	}
	return os.Chmod(filepath.Dir(cleanPath), 0600)
}
