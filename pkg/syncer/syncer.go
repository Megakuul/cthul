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

package syncer

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"cthul.io/cthul/pkg/db"
	"cthul.io/cthul/pkg/log"
)

// Syncer is a utility component that helps operators to apply a state from the database to the system.
// The syncer uses multiple goroutines that incrementally and periodically call the operator specific function
// to apply state from the database. This two way system allows the system to behave deterministic by
// periodically applying the defined state, while also using incremental updates to apply changes immediately.
type Syncer struct {
	// rootCtx is used as root for every context that is created on the syncer. This ensures that
	// canceling it will stop every started goroutine tracked by operationWg.
	rootCtx context.Context
	rootCtxCancel context.CancelFunc

	client db.Client
	logger log.Logger
	
	// operationWg tracks every single operation started on the syncer, this ensures that even goroutines
	// that were removed from the trackMap without waiting for them, are not leaking.
	operationWg sync.WaitGroup

	// trackMap tracks all operational syncer routines by holding a cancel function of each one.
	// Routines may be removed from the trackMap without waiting for them to finish by passing false to the
	// wait flag, to ensure goroutines are not leaking, they are tracked by the trackMap AND the operationWg.
	trackMapLock sync.Mutex
	trackMap map[string]func(bool)
}

// Add adds a routine to the syncer. This means that the syncer starts two goroutines one that incrementally
// watches $prefix$uuid and one that is executed periodically in the specified interval. Both goroutines
// fire $fn periodically / on change, passing the value of $prefix$uuid to $fn.
func (s *Syncer) Add(uuid, prefix string, interval int64, fn func(context.Context, string) error) error {
	s.trackMapLock.Lock()
	defer s.trackMapLock.Unlock()

	if _, ok := s.trackMap[uuid]; ok {
		return fmt.Errorf("syncer routine for device '%s' is already running", uuid)
	}

	syncKey := fmt.Sprintf("%s/%s", strings.TrimSuffix(prefix, "/"), uuid)
	
	funcWg := sync.WaitGroup{}
	funcCtx, funcCtxCancel := context.WithCancel(s.rootCtx)

	s.operationWg.Add(1)
	funcWg.Add(1)
	go func() {
		defer s.operationWg.Done()
		defer funcWg.Done()
		for {
			ctx, cancel := context.WithTimeout(funcCtx, time.Duration(interval) * time.Second)
			defer cancel()
				
			state, err := s.client.Get(ctx, syncKey)
			if err!=nil {
				s.logger.Err("syncer", fmt.Sprintf("failed to load key '%s': %s",
					syncKey, err.Error(),
				))
			} else {
				err = fn(ctx, state)
				if err!=nil {
					s.logger.Err("syncer", fmt.Sprintf("cannot apply state to '%s': %s", uuid, err.Error()))
				} else {
					s.logger.Debug("syncer", fmt.Sprintf("applied state '%s' to '%s'", state, uuid))
				}
			}
			
			select {
			case <-funcCtx.Done():
				return
			case <-ctx.Done():
			}
		}
	}()

	s.operationWg.Add(1)
	funcWg.Add(1)
	go func() {
		defer s.operationWg.Done()
		defer funcWg.Done()
		err := s.client.WatchRange(funcCtx, syncKey, func(_, state string, err error) {
			if err!=nil {
				s.logger.Err("syncer", fmt.Sprintf("failed to load key '%s': %s",
					syncKey, err.Error(),
				))
				return
			}
			err = fn(funcCtx, state)
			if err!=nil {
				s.logger.Err("syncer", fmt.Sprintf("cannot apply state to '%s': %s", uuid, err.Error()))
			} else {
				s.logger.Debug("syncer", fmt.Sprintf("applied state '%s' to '%s'", state, uuid))
			}
		})
		if err!=nil {
			s.logger.Crit("syncer", fmt.Sprintf(
				"failed to watch '%s' state: %s; exiting state watcher...", syncKey, err.Error(),
			))
		}
	}()

	s.trackMap[uuid] = func(wait bool) {
		funcCtxCancel()
		if wait {
			funcWg.Wait()
		}
	}

	return nil
}

// Remove stops and removes a syncer routine (idempotent). Specifing the wait flag, ensures the function
// waits before returning until the routine is fully closed. Without wait flag, the syncer stops the context
// of the function and removes it without waiting for it to exit.
func (s *Syncer) Remove(uuid string, wait bool) {
	s.trackMapLock.Lock()
	defer s.trackMapLock.Unlock()

	cancel, ok := s.trackMap[uuid]
	if !ok {
		return
	}

	cancel(wait)
	delete(s.trackMap, uuid)
}

// Shutdown cancels all running syncers, whether they are tracked by the trackMap or not. It waits
// until every single goroutine is done.
func (s *Syncer) Shutdown() {
	s.rootCtxCancel()
	s.operationWg.Wait()
}
