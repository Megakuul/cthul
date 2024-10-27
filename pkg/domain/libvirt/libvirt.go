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
	"net/url"

	"github.com/digitalocean/go-libvirt"
)

type LibvirtController struct {
	client *libvirt.Libvirt
}

type LibvirtControllerOption func(*LibvirtController)

func NewLibvirtController(opts ...LibvirtControllerOption) *LibvirtController {
	controller := &LibvirtController{

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
