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

// VIDEO_TYPE specifies what protocol / type of host-side video adapter is used.
// Currently only spice is supported which uses the qemu host spice server.
type VIDEO_TYPE string
const (
	VIDEO_SPICE VIDEO_TYPE = "spice"
)

// Video holds all information about a video adapter device.
type Video struct {
  Reqnode string `json:"reqnode"`
	Node string `json:"node"`
	Type VIDEO_TYPE `json:"type"`
	Path string `json:"path"`
}
