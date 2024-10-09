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
	"cthul.io/cthul/pkg/lifecycle"
	"cthul.io/cthul/pkg/log/bootstrap"
	"cthul.io/cthul/pkg/log/runtime"
	"go.etcd.io/etcd/client/v3"
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

	dbClient, err := clientv3.New(clientv3.Config{
		Endpoints: []string{config.Database.Addr},
		Username: config.Database.Username,
		Password: config.Database.Password,
		DialTimeout: time.Second * time.Duration(config.Database.TimeoutTTL),
	})
	if err!=nil {
		return err
	}

	res, err := dbClient.Auth.AuthStatus(context.TODO())
	if err!=nil {
		return err
	}
	fmt.Println(res.Enabled)
	
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


	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	exitSignal := <-signalChan
	bootLogger.Info("service",
		fmt.Sprintf("received %s; service is being shutdown...", exitSignal.String()),
	)
	
	return nil
}
