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
	"os"
	"path/filepath"
	"strings"

	"cthul.io/cthul/pkg/api/wave/v1/video"
	"google.golang.org/protobuf/proto"
)

func (o *Operator) synchronize() {
  o.syncer.Add("/WAVE/VIDEO/REQNODE/", o.updateCycleTTL, func(ctx context.Context, k, reqnode string) error {
    id := strings.TrimPrefix(k, "/WAVE/VIDEO/REQNODE/")
    configKey := fmt.Sprintf("/WAVE/VIDEO/CONFIG/%s", id)
    if reqnode == o.nodeId {
      o.syncer.Add(configKey, o.syncCycleTTL, func(ctx context.Context, k, v string) error {
        err := o.applyConfig(v)
        if err!=nil {
          return err
        }
        _, err = o.client.Set(ctx, fmt.Sprintf("/WAVE/VIDEO/NODE/%s", id), reqnode, 0)
        if err!=nil {
          return err
        }
        return nil
      })
    } else {
      o.syncer.Remove(configKey, false)
    }
    return nil
  })
}

func (o *Operator) applyConfig(rawConfig string) error {
  config := &video.VideoConfig{}
  err := proto.Unmarshal([]byte(rawConfig), config)
  if err!=nil {
		return fmt.Errorf("failed to parse config: %w", err)
  }

  cleanPath := filepath.Join(o.runRoot, config.Path)  
  if !strings.HasPrefix(cleanPath, o.runRoot) {
    return fmt.Errorf("device socket path is not allowed to escape the run root ('%s' => '%s')", o.runRoot, cleanPath)
  }
  err = os.MkdirAll(filepath.Dir(cleanPath), 0600)
  if err!=nil {
    return err
  }
  return os.Chmod(filepath.Dir(cleanPath), 0600)
}
