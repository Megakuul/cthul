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

type VIDEO_TYPE string
const (
	VIDEO_SPICE VIDEO_TYPE = "spice"
)

// Video holds all information about a video adapter device.
type Video struct {
	Type VIDEO_TYPE
	Path string // unix socket path
}

func (d *DeviceController) ListVideo(id string) error {

}

// CreateVideo creates a new wave video adapter device.
func (d *DeviceController) CreateVideo(id string) error {
	
}

func (d *DeviceController) DeleteVideo(id string) error {

}

func (d *DeviceController) LookupVideo(id string) error {

}

func (d *DeviceController) GetPort(id string) (<-chan []byte, chan<-[]byte, error) {

}

func (d *DeviceController) AttachVideo(id string) error {

}

func (d *DeviceController) DetachVideo(id string) error {

}
