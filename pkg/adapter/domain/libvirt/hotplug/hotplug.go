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

package hotplug

import (
	"fmt"

	libvirtstruct "cthul.io/cthul/pkg/adapter/domain/libvirt/structure"
	"github.com/digitalocean/go-libvirt"
)

// Hotplugger provides methods to hotplug libvirt changes based on the xml configuration.
// Simply defining the xml does only change the domains persistent config but does not hotplug updates,
// the hotplugger takes components that are hotpluggable and updates them with the appropriate libvirt rpc calls.
type Hotplugger struct {
}

func New() *Hotplugger {
	return &Hotplugger{
	}
}


func (l *Hotplugger) Hotplug(libvirt *libvirt.Libvirt, config *libvirtstruct.Domain) error {
	return fmt.Errorf("not implemented")
}
