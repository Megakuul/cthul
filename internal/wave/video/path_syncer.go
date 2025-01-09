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
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"
	
	"cthul.io/cthul/pkg/wave/video/structure"
)

// updateStateSyncer updates the state syncer map based on the new node. If the video is located on the local
// node, a state syncer is started. If a state syncer associated with the video is running, but the new node is
// not == the local one, the state syncer is stopped.
// The state syncer consists of two routines, one routine operates in periodical ticks pushing the state in a
// configured interval. The other routine watches the database and pushes state changes immediately upon update.
func (o *Operator) updatePathSyncer(ctx context.Context, uuid, node, nodeRequest string) {
	if nodeRequest == o.nodeId && node != o.nodeId {
		o.setupPath(ctx, uuid)
	} else if nodeRequest != o.nodeId && node == o.nodeId {
		o.cleanupPath(ctx, uuid)
	}
}

// applyState tries to apply the desired power state to the local video.
func (o *Operator) setupPath(uuid, state string) error {
	return nil
}

// applyState tries to apply the desired power state to the local video.
func (o *Operator) cleanupPath(uuid, state string) error {
	return nil
}
