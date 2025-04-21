/**
 * Cthul System
 *
 * Copyright (C) 2025 Linus Ilian Moser <linus.moser@megakuul.ch>
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
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

// block adapter provides an abstraction layer over a distributed block replication system.
// it is used primarily by granit to provision block clusters and allocate / attach blockdevs.
package block

import (
)

type Adapter interface {
  // PoC
  Apply()
  // drbdadm up <disk>
  // if node!=reqnode -> drbdadm secondary <node>
  // drbdadm primary <reqnode>
  Destroy()
  // if localnode==node && localnode!=reqnode -> drbdadm secondary <localnode>
  // drbadm down <disk>
  // umount /dev/drbdxy
  // umount /dev/loopdev
  // rm -rf /device
}

