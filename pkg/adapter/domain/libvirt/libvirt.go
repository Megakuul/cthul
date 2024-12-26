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

package libvirt

import (
	"context"
	"encoding/hex"
	"fmt"
	"net/url"
	"strings"
	"sync"

	"cthul.io/cthul/pkg/adapter/domain/libvirt/generator"
	"cthul.io/cthul/pkg/log"
	"cthul.io/cthul/pkg/log/discard"
	"github.com/digitalocean/go-libvirt"
)

type LibvirtAdapter struct {
	initLock sync.Mutex
	client *libvirt.Libvirt
	logger log.Logger

	generator *generator.LibvirtGenerator
}

type LibvirtAdapterOption func(*LibvirtAdapter)

func NewLibvirtAdapter(generator *generator.LibvirtGenerator, opts ...LibvirtAdapterOption) *LibvirtAdapter {
	controller := &LibvirtAdapter{
		initLock: sync.Mutex{},
		client: nil,
		logger: discard.NewDiscardLogger(),
		generator: generator,
	}

	for _, opt := range opts {
		opt(controller)
	}

	return controller
}

// WithLogger adds a logger to the libvirt controller.
func WithLogger(logger log.Logger) LibvirtAdapterOption {
	return func (l *LibvirtAdapter) {
		l.logger = logger
	}
}

// initClient creates the underlying libvirt connection client if not already initialized.
func (l *LibvirtAdapter) initClient() error {
	l.initLock.Lock()
	defer l.initLock.Unlock()
	if l.client!=nil {
		return nil
	}
	
	uri, _ := url.Parse(string(libvirt.QEMUSystem))
	client, err := libvirt.ConnectToURI(uri)
	if err!=nil {
		return err
	}
	l.client  = client
	return nil
}

// parseUUID tries to convert a uuid string (either with or without hyphens) into a libvirt.UUID.
func (l *LibvirtAdapter) parseUUID(id string) (libvirt.UUID, error) {
	rawStr := strings.ReplaceAll(id, "-", "")
	if len(rawStr) != 2 * libvirt.UUIDBuflen {
		return [libvirt.UUIDBuflen]byte{}, fmt.Errorf(
			"failed to parse uuid: expected %d characters", libvirt.UUIDBuflen,
		)
	}

	uuidBuffer, err := hex.DecodeString(rawStr)
	if err!=nil {
		return [libvirt.UUIDBuflen]byte{}, fmt.Errorf("failed to parse uuid: cannot hex decode id")
	}
	uuid := [libvirt.UUIDBuflen]byte{}
	copy(uuid[:], uuidBuffer)
	return uuid, nil
}

// serializeUUID converts a libvirt uuid into a uuid string with hyphens.
func (l *LibvirtAdapter) serializeUUID(uuid libvirt.UUID) (string, error) {
	uuidStr := hex.EncodeToString(uuid[:])
	if len(uuidStr) != 32 {
		return "", fmt.Errorf("failed to serialize uuid: expected encoded hex string with %d characters", 32)
	}
	return fmt.Sprintf(
		"%s-%s-%s-%s-%s",
		uuidStr[:8], uuidStr[8:12], uuidStr[12:16], uuidStr[16:20], uuidStr[20:32],
	), nil
}

// Terminate stops and closes the libvirt controller.
// The context is currently not utilized due to the lack of context handling in the underlying libvirt library.
func (l* LibvirtAdapter) Terminate(ctx context.Context) error {
	if l.client != nil {
		return l.client.Disconnect()
	}
	return nil
}