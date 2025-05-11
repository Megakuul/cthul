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

package drbd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"cthul.io/cthul/pkg/api/granit/v1/disk"
	"golang.org/x/sys/unix"
)

type Adapter struct {
	executable string
	storage    string
}

type Option func(*Adapter)

func New(opts ...Option) *Adapter {
	adapter := &Adapter{
		executable: "drbdadm",
		storage:    "/var/lib/cthul/granit/",
	}

	for _, opt := range opts {
		opt(adapter)
	}

	return adapter
}

func WithExecutable(path string) Option {
	return func(a *Adapter) {
		a.executable = path
	}
}

func WithStorage(path string) Option {
	return func(a *Adapter) {
		a.storage = path
	}
}

func (a *Adapter) Apply(ctx context.Context, id string, config *disk.DiskConfig, cluster *disk.DiskCluster) error {
	path := filepath.Join(a.storage, id)
	if !strings.HasPrefix(filepath.Clean(path), a.storage) {
		return fmt.Errorf("device uses a path that escapes the storage root '%s'", a.storage)
	}

  requiresInit := false

  file, err := os.OpenFile(path, os.O_RDWR, 0600)
  if err!=nil && !os.IsNotExist(err) {
    return err
  } else if os.IsNotExist(err) {
    requiresInit = true
    os.MkdirAll(filepath.Dir(path), 0755)
		file, err = os.CreateTemp("", fmt.Sprintf("granit-prep-disk-%s-", id))
		if err != nil {
			return err
		}
		cmd := exec.CommandContext(ctx, "drbdadm", "create-md", id)
		err = cmd.Run()
		if err != nil {
			return err
		}
  }
  defer file.Close()

  devPath := fmt.Sprintf("/dev/cthul/granit/%s", id)
	_, err = os.Lstat(devPath)
  if err == nil {
    if !requiresInit {
      return nil
    }

    // file doesnt exist but the device does 
    // this might happen if a sysadmin deletes the file manually.
    // recover by removing the link and proceeding to allocate a new loopdev.
   
    detachLoopDev(devPath)
    err = os.Remove(devPath)
    if err!=nil {
      return fmt.Errorf("failed to remove old device link: %w", err)
    }
  } 

  os.MkdirAll(filepath.Dir(devPath), 0755)
  loopDevPath, err := attachLoopDev(file)
  if err!=nil {
    return err
  }

  err = os.Symlink(loopDevPath, devPath)
  if err!=nil {
    rErr := detachLoopDev(loopDevPath)
    if rErr != nil {
      return fmt.Errorf("%w; rollback failed: %v", err, rErr)
    }
    return err
  }

  if requiresInit {
	  err = exec.CommandContext(ctx,
      a.executable, "create-md", id,
    ).Run()
    if err!=nil {
      rErr := detachLoopDev(loopDevPath)
      if rErr != nil {
        return fmt.Errorf("%w; rollback failed: %v", err, rErr)
      }
      rErr = os.Remove(devPath)
      if rErr!= nil {
        return fmt.Errorf("%w; rollback failed: %v", err, rErr)
      }
      return fmt.Errorf("failed to initialize drbd disk: %w", err)
    }
    err = os.Rename(file.Name(), path)
    if err!=nil {
      return fmt.Errorf("failed to commit: %w", err)
    } 
  }

  return nil
}

// attaches the file to a free loop device. returns the loopdev path.
func attachLoopDev(file *os.File) (string, error) {
	ctrl, err := unix.Open("/dev/loop-control", os.O_RDONLY, 0)
	if err != nil {
		return "", fmt.Errorf("failed to open loopcontroller: %w", err)
	}
	defer unix.Close(ctrl)

	loopNumber, err := unix.IoctlRetInt(ctrl, unix.LOOP_CTL_GET_FREE)
	if err != nil {
		return "", fmt.Errorf("failed to get free loopdev: %w", err)
	}

  loopDevPath := fmt.Sprintf("/dev/loop%d", loopNumber)
	loop, err := unix.Open(loopDevPath, os.O_RDWR, 0)
	if err != nil {
		return "", fmt.Errorf("failed to open loopdev: %w", err)
	}
	defer unix.Close(loop)

  err = unix.IoctlLoopConfigure(loop, &unix.LoopConfig{
    Fd: uint32(file.Fd()),
    Info: unix.LoopInfo64{
      Flags: unix.LO_FLAGS_DIRECT_IO, // avoid double fs caching
    },
  })
  if err!=nil {
    return "", fmt.Errorf("failed to attach loopdev: %w", err)
  }

  return loopDevPath, nil
}

// detaches the loop device from its associated file. 
func detachLoopDev(path string) error {
	loop, err := unix.Open(path, os.O_RDWR, 0)
	if err != nil {
		return fmt.Errorf("failed to open loopdev: %w", err)
	}
	defer unix.Close(loop)

  err = unix.IoctlSetInt(loop, unix.LOOP_CLR_FD, 0)
  if err!=nil {
    return fmt.Errorf("failed to detach loopdev: %w", err)
  }
  return nil
}

func Destroy(ctx context.Context, id string) error {

}

func Primary(ctx context.Context) error {

  err = exec.CommandContext(ctx, 
    a.executable, "up", id,
  ).Run()
  if err!=nil {
    return fmt.Errorf("failed to attach drbd disk: %w", err)
  }
}

func Secondary(ctx context.Context) error {

  err = exec.CommandContext(ctx, 
    a.executable, "up", id,
  ).Run()
  if err!=nil {
    return fmt.Errorf("failed to attach drbd disk: %w", err)
  }
}
