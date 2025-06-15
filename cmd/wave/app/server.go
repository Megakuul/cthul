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
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cthul.io/cthul/internal/wave/api"
	"cthul.io/cthul/internal/wave/scheduler"
	"cthul.io/cthul/pkg/adapter/domain/libvirt"
	"cthul.io/cthul/pkg/adapter/domain/libvirt/generator"
	"cthul.io/cthul/pkg/adapter/domain/libvirt/hotplug"
	"cthul.io/cthul/pkg/db/etcdv3"
	"cthul.io/cthul/pkg/granit/disk"
	"cthul.io/cthul/pkg/lifecycle"
	"cthul.io/cthul/pkg/proton/inter"
	"cthul.io/cthul/pkg/wave/domain"
	"cthul.io/cthul/pkg/wave/node"
	"cthul.io/cthul/pkg/wave/serial"
	"cthul.io/cthul/pkg/wave/video"
	domainop "cthul.io/cthul/internal/wave/domain"
	nodeop "cthul.io/cthul/internal/wave/node"
	serialop "cthul.io/cthul/internal/wave/serial"
	videoop "cthul.io/cthul/internal/wave/video"
	"github.com/lmittmann/tint"
)

// Run is the root entrypoint of the service.
// This function does only fail if a critical error occurs while setting up the system,
// otherwise it will run until an os level signal (SIGINT/TERM) is received.
func Run(config *BaseConfig) error {
  logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{
    AddSource: true,
    TimeFormat: time.Kitchen,
  }))

	lifecycleManager := lifecycle.NewManager(logger.With("comp", "lifecycle-manager"))
	defer lifecycleManager.TerminateParallel(
		time.Second * time.Duration(config.Lifecycle.TerminationTTL),
	)

	dbClient := etcdv3.New([]string{config.Database.Addr},
		etcdv3.WithAuth(config.Database.Username, config.Database.Password),
		etcdv3.WithDialTimeout(time.Second * time.Duration(config.Database.TimeoutTTL)),
		etcdv3.WithSkipVerify(config.Database.SkipVerify),
	)
	lifecycleManager.AddHook(dbClient.Terminate)

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

  nodeController := node.New(config.NodeId, dbClient)

  videoController := video.New(config.NodeId, dbClient)
  serialController := serial.New(config.NodeId, dbClient)
  diskController := disk.New(config.NodeId, dbClient)
  interController := inter.New(config.NodeId, dbClient)
  domainAdapter := libvirt.New(
    generator.New(config.NodeId, 
      videoController, 
      serialController, 
      diskController, 
      interController,
    ),
    hotplug.New(),
  )
  domainController := domain.New(config.NodeId, dbClient, domainAdapter, domain.WithRunRoot(
    "/run/cthul/wave/",
  ))
	scheduler := scheduler.New(logger.With("comp", "scheduler"), dbClient, 
    domainController, nodeController,
    scheduler.WithCycleTTL(config.Scheduler.CycleTTL),
    scheduler.WithRescheduleCycles(config.Scheduler.RescheduleCycles),
	)
	scheduler.ServeAndDetach()
	lifecycleManager.AddHook(scheduler.Terminate)

	domainOperator := domainop.New(logger.With("comp", "domain-operator"), dbClient, domainAdapter,
		// TODO
	)
	domainOperator.ServeAndDetach()
	lifecycleManager.AddHook(domainOperator.Terminate)

	nodeOperator := nodeop.New(logger.With("comp", "node-operator"), dbClient, 
		nodeop.WithAffinity("todo", "todo2"),
		// TODO
	)
	nodeOperator.ServeAndDetach()
	lifecycleManager.AddHook(nodeOperator.Terminate)

	serialOperator := serialop.New(logger.With("comp", "serial-operator"), dbClient, 
		// TODO
	)
	serialOperator.ServeAndDetach()
	lifecycleManager.AddHook(serialOperator.Terminate)
	
	videoOperator := videoop.New(logger.With("comp", "video-operator"), dbClient, 
		// TODO
	)
	videoOperator.ServeAndDetach()
	lifecycleManager.AddHook(videoOperator.Terminate)

	apiCertificate, err := tls.LoadX509KeyPair(config.Api.CertFile, config.Api.KeyFile)
	if err!=nil {
		return err
	}
  _ = apiCertificate
	apiEndpoint := api.New(config.Api.Addr,
    api.WithLogger(logger.With("comp", "api-endpoint")),
    api.WithOrigins(config.Api.Origins),
    api.WithDomain(domainController),
    api.WithVideo(videoController),
    api.WithSerial(serialController),
    api.WithNode(nodeController),
	)
	if err := apiEndpoint.ServeAndDetach(); err!=nil {
		return err
	}
	lifecycleManager.AddHook(apiEndpoint.Terminate)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	exitSignal := <-signalChan
  logger.Info(fmt.Sprintf("received %s; service is being shutdown...", exitSignal.String()))
	
	return nil
}
