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

package app

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"cthul.io/cthul/internal/wave/api"
	"cthul.io/cthul/internal/wave/scheduler"
	"cthul.io/cthul/pkg/db/etcdv3"
	"cthul.io/cthul/pkg/elect"
	"cthul.io/cthul/pkg/lifecycle"
	"cthul.io/cthul/pkg/log/adapter"
	"cthul.io/cthul/pkg/log/bootstrap"
	"cthul.io/cthul/pkg/log/runtime"
	"cthul.io/cthul/pkg/wave/domain"
	"cthul.io/cthul/pkg/wave/node"
	"google.golang.org/grpc/grpclog"
)

// Run is the root entrypoint of the service.
// This function does only fail if a critical error occurs while setting up the system,
// otherwise it will run until an os level signal (SIGINT/TERM) is received.
func Run(config *BaseConfig) error {
	loggerIOLock := &sync.Mutex{}
	
	bootLogger := bootstrap.NewBootstrapLogger("wave",
		bootstrap.WithIOLock(loggerIOLock),
		bootstrap.WithLevel(config.Logging.Level),
		bootstrap.WithTrace(config.Logging.Trace),
	)

	terminationManager := lifecycle.NewTerminationManager(
		lifecycle.WithLogger(bootLogger),
	)
	defer terminationManager.TerminateParallel(
		time.Second * time.Duration(config.Lifecycle.TerminationTTL),
	)

	coreLogger := runtime.NewRuntimeLogger("wave",
		runtime.WithIOLock(loggerIOLock),
		runtime.WithLevel(config.Logging.Level),
		runtime.WithTrace(config.Logging.Trace),
		runtime.WithLogBuffer(config.Logging.Buffer),
	)
	coreLogger.ServeAndDetach()
	terminationManager.AddHook(coreLogger.Terminate)

	grpclog.SetLoggerV2(adapter.NewGrpcLogAdapter("grpc",
		adapter.WithWarnLog(coreLogger.Warn),
		adapter.WithErrLog(coreLogger.Err),
		adapter.WithCritLog(coreLogger.Crit),
	))

	dbClient := etcdv3.NewEtcdClient([]string{config.Database.Addr},
		etcdv3.WithAuth(config.Database.Username, config.Database.Password),
		etcdv3.WithDialTimeout(time.Second * time.Duration(config.Database.TimeoutTTL)),
		etcdv3.WithSkipVerify(config.Database.SkipVerify),
	)
	terminationManager.AddHook(dbClient.Terminate)

	if config.Database.Healthcheck {
		ctx, cancel := context.WithTimeout(
			context.Background(), time.Second * time.Duration(config.Database.TimeoutTTL),
		)
		if err := dbClient.CheckEndpointHealth(ctx); err!=nil {
			cancel()
			return fmt.Errorf("database healthcheck failed: %s", err.Error())
		}
		cancel()
	}

  nodeController := node.NewController(dbClient)
  domainController := domain.NewController(dbClient)
	scheduler := scheduler.NewScheduler(dbClient,
    scheduler.WithLogger(coreLogger),
    scheduler.WithCycleTTL(config.Scheduler.CycleTTL),
    scheduler.WithRescheduleCycles(config.Scheduler.RescheduleCycles),
	)
	scheduler.ServeAndDetach()
	terminationManager.AddHook(scheduler.Terminate)

	apiCertificate, err := tls.LoadX509KeyPair(config.Api.CertFile, config.Api.KeyFile)
	if err!=nil {
		return err
	}
	apiEndpoint := api.NewApiEndpoint(config.Api.Addr, apiCertificate,
		api.WithApplicationLog(coreLogger),
		api.WithSystemLog(coreLogger),
		api.WithIdleTimeout(time.Second * time.Duration(config.Api.IdleTTL)),
	)
	if err := apiEndpoint.ServeAndDetach(); err!=nil {
		return err
	}
	terminationManager.AddHook(apiEndpoint.Terminate)

	electController := elect.NewElectController(dbClient, "/WAVE/LEADER",
		elect.WithLocalLeader(config.Election.Contest, config.NodeId, config.Election.Cash),
		elect.WithContestTTL(config.Election.ContestTTL),
		elect.WithContestHooks(scheduler.SetLeaderState, apiEndpoint.SetLeaderState),
		elect.WithLogger(coreLogger),
	)
	electController.ServeAndDetach()
	terminationManager.AddHook(electController.Terminate)
	

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	exitSignal := <-signalChan
	bootLogger.Info("service",
		fmt.Sprintf("received %s; service is being shutdown...", exitSignal.String()),
	)
	
	return nil
}
