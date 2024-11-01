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

package generator

import (
	"fmt"

	cthulstruct "cthul.io/cthul/pkg/domain/structure"
)

// Attach installs / locks all devices that are required by the domain config.
func (l *LibvirtGenerator) Attach(config *cthulstruct.Domain) error {
	for _, device := range config.BlockDevices {
		// PoC: l.granit.AttachBlock(device.GranitBlockDeviceId)
		_ = device
	}

	for _, device := range config.NetworkDevices {
		// PoC: l.proton.AttachInterface(device.ProtonNetworkDeviceId)
		_ = device
	}

	for _, device := range config.SerialDevices {
		// PoC: l.wave.AttachSerial(device.WaveSerialDeviceId)
		_ = device
	}

	for _, device := range config.GraphicDevices {
		// PoC: l.wave.AttachGraphic(device.WaveGraphicDeviceId)
		_ = device
	}
	
	return fmt.Errorf("not implemented biatch")
}
