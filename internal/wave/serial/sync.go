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
)

func (o *Operator) synchronize() {
  o.syncer.Add("/WAVE/VIDEO/REQNODE", o.updateCycleTTL, func(ctx context.Context, k, reqnode string) error {
    uuid := strings.TrimPrefix(k, "/WAVE/VIDEO/REQNODE/")
    pathKey := fmt.Sprintf("/WAVE/VIDEO/PATH/%s", uuid)
    if reqnode == o.nodeId {
      o.syncer.Add(pathKey, o.pathCycleTTL, func(ctx context.Context, k, path string) error {
        err := o.ensurePath(o.waveRunRoot, path)
        if err!=nil {
          return err
        }
        _, err = o.client.Set(ctx, fmt.Sprintf("/WAVE/VIDEO/NODE/%s", uuid), reqnode, 0)
        if err!=nil {
          return err
        }
        return nil
      })
    } else {
      o.syncer.Remove(pathKey, false)
    }
    return nil
  })
}

func (o *Operator) ensurePath(base, path string) error {
  cleanPath := filepath.Join(base, path)  
  if !strings.HasPrefix(cleanPath, base) {
    return fmt.Errorf("device socket path is not allowed to escape the run root ('%s' => '%s')", base, cleanPath)
  }
  err := os.MkdirAll(filepath.Dir(cleanPath), 0600)
  if err!=nil {
    return err
  }
  return os.Chmod(filepath.Dir(cleanPath), 0600)
}
