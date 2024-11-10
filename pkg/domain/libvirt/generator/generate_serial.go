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
	libvirtstruct "cthul.io/cthul/pkg/domain/libvirt/structure"
	cthulstruct "cthul.io/cthul/pkg/domain/structure"
)

// generateSerial generates a libvirt serial device from the cthul serial device.
func (l *LibvirtGenerator) generateSerial(device *cthulstruct.SerialDevice) (*libvirtstruct.Serial, error) {

	return nil, nil
}
