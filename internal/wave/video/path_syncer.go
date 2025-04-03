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
 * along with this program. If not, see <https://www.gnu.org/licenses/>.
 */

package video

import (
	"fmt"
)

// updateStateSyncer updates the state syncer map based on the new node. If the video is located on the local
// node, a state syncer is started. If a state syncer associated with the video is running, but the new node is
// not == the local one, the state syncer is stopped.
// The state syncer consists of two routines, one routine operates in periodical ticks pushing the state in a
// configured interval. The other routine watches the database and pushes state changes immediately upon update.
func (o *Operator) updatePathSyncer(uuid, node, nodeRequest string) {
	o.pathSyncersLock.Lock()
	defer o.pathSyncersLock.Unlock()

	if node!=o.nodeId {
		if cancel, ok := o.pathSyncers[uuid]; ok {
			o.logger.Debug("video-operator", fmt.Sprintf(
				"removing config synchronizer for domain '%s' on node '%s'...", uuid, node,
			))
			cancel()
			delete(o.pathSyncers, uuid)
		}
		return		
	}

	if nodeRequest!="" {
		
	}
	
	o.logger.Debug("video-operator", fmt.Sprintf(
		"setting up path synchronizer for video device '%s' on node '%s'...", uuid, node,
	))
	
	if _, ok := o.pathSyncers[uuid]; ok {
		o.logger.Debug("video-operator", fmt.Sprintf(
			"path synchronization for video device '%s' is already running on '%s'; skipping setup...",
			uuid, node,
		))
		return
	}
}
