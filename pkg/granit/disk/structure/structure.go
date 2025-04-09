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

package structure

type DISK_TYPE string

const (
  DISK_BLOCK DISK_TYPE = "block"
  DISK_FILE DISK_TYPE = "file"
)

type DISK_FORMAT string

const (
  DISK_RAW DISK_FORMAT = "raw"
  DISK_QCOW2 DISK_FORMAT = "qcow2"
)

type Disk struct {
  Reqnode string `json:"reqnode"`
	Node string `json:"node"`
  Type DISK_TYPE `json:"type"`
  Format DISK_FORMAT `json:"format"`
	Path string `json:"path"` // core
  Readonly bool `json:"readonly"`
  Error error `json:"-"`
}
