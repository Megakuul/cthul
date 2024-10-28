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

	"github.com/digitalocean/go-libvirt"
)

type LibvirtController struct {
	client *libvirt.Libvirt
}

type LibvirtControllerOption func(*LibvirtController)

func NewLibvirtController(opts ...LibvirtControllerOption) *LibvirtController {
	controller := &LibvirtController{
		client: nil,
	}

	for _, opt := range opts {
		opt(controller)
	}

	return controller
}


// initClient creates the underlying libvirt connection client if not already initialized.
func (l *LibvirtController) initClient() error {
	if l.client != nil {
		return nil
	}
	uri, _ := url.Parse(string(libvirt.QEMUSystem))
	client, err := libvirt.ConnectToURI(uri)
	if err!=nil {
		return err
	}
	l.client = client
	return nil
}

// parseUUID tries to convert a uuid string (either with or without hyphens) into a libvirt.UUID.
func (l *LibvirtController) parseUUID(id string) (libvirt.UUID, error) {
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

// Terminate stops and closes the libvirt controller.
// The context is currently not utilized due to the lack of context handling in the underlying libvirt library.
func (l* LibvirtController) Terminate(ctx context.Context) error {
	if l.client != nil {
		return l.client.Disconnect()
	}
	return nil
}
