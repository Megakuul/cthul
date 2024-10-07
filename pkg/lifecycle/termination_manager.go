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

package lifecycle

import (
	"context"
	"fmt"
	"sync"
	"time"

	"cthul.io/cthul/pkg/log"
	"cthul.io/cthul/pkg/log/discard"
)

// TerminationManager provides a basic termination lifecycle.
// It captures termination hooks of various detached components, allowing the programm
// to terminate all of those components in a controlled manner.
type TerminationManager struct {
	logger log.Logger
	hooks []func(context.Context) error
}

type TerminationManagerOption func(*TerminationManager)

// NewTerminationManager creates a new termination manager.
func NewTerminationManager(opts ...TerminationManagerOption) *TerminationManager {
	manager := &TerminationManager{
		logger: discard.NewDiscardLogger(),
		hooks: []func(context.Context) error{},
	}

	for _, opt := range opts {
		opt(manager)
	}

	return manager
}

// WithLogger adds a logger to the termination manager
// The logger is used to log errors and some debug information while terminating.
func WithLogger(logger log.Logger) TerminationManagerOption {
	return func(t *TerminationManager) {
		t.logger = logger
	}
}

// AddHook adds a termination hook. The provided terminateFunc should try to gracefully close / shutdown
// the component. If the context exceeds, the function is expected to immediately forcefully close / shutdown
// the component. In case of an error, the error should be returned.
func (t *TerminationManager) AddHook(terminateFunc func(context.Context) error) {
	t.hooks = append(t.hooks, terminateFunc)
}

// TerminateParallel executes all termination hooks at the same time.
// After the provided timeout, the hooks stop the graceful shutdown process and immediately
// and forcefully close the components. The function blocks until all hooks returned.
func (t *TerminationManager) TerminateParallel(timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	t.logger.Debug("termination_manager",
		fmt.Sprintf("initiating termination of %d hooks in %v...", len(t.hooks), timeout))

	exitWg := sync.WaitGroup{}
	for i, hook := range t.hooks {
		exitWg.Add(1)
		go func() {
			defer exitWg.Done()
			if err := hook(ctx); err!=nil {
				t.logger.Warn("termination_manager",
					fmt.Sprintf("termination of hook %d failed: %v", i, err))
				return
			}
			t.logger.Debug("termination_manager",
				fmt.Sprintf("hook %d terminated successfully", i))
		}()
	}

	exitWg.Wait()
}
