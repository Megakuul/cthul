/**
 * Cthul System
 *
 * Copyright (C) 2024 Linus Ilian Moser <linus.moser@megakuul.ch>
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

package generator

import (
	"fmt"

	cthulstruct "cthul.io/cthul/pkg/adapter/domain/structure"
)

// Detach releases all devices that were used by the domain config.
func (l *Generator) Detach(config *cthulstruct.Domain) error {
	for _, device := range config.NetworkDevices {
		// PoC: l.proton.DetachInterface(device.DeviceId)
		_ = device
	}

	for _, device := range config.StorageDevices {
		// PoC: l.granit.DetachStorage(device.DeviceId)
		_ = device
	}
	
	for _, device := range config.SerialDevices {
		// PoC: l.wave.DetachSerial(device.DeviceId)
		_ = device
	}

	for _, device := range config.VideoAdapters {
		// PoC: l.wave.DetachVideo(config.adapter.DeviceId)
		_ =  device
	}

	
	return fmt.Errorf("not implemented mr")
}
