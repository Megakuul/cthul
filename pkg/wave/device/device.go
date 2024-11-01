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

package device

import "cthul.io/cthul/pkg/db"

// DeviceController provides a controller interface for the wave device operator.
// It provides operations to request, manage and release wave devices.
type DeviceController struct {
	client db.Client
}

type DeviceControllerOption func(*DeviceController)

func NewDeviceController(client db.Client, opts ...DeviceControllerOption) *DeviceController {
	controller := &DeviceController{
		client: client,
	}

	for _, opt := range opts {
		opt(controller)
	}

	return controller
}


