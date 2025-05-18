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

package app

import (
	"github.com/BurntSushi/toml"
	"os"
)

// BaseConfig represents the configuration model
type BaseConfig struct {
	NodeId    string          `toml:"nodeid" validate:"required"`
	Lifecycle LifecycleConfig `toml:"lifecycle"`
	Logging   LoggingConfig   `toml:"logging"`
	Database  DatabaseConfig  `toml:"db"`
	Scheduler SchedulerConfig `toml:"scheduler"`
	Api       ApiConfig       `toml:"api"`
}

type LifecycleConfig struct {
	TerminationTTL int64 `toml:"termination_ttl" validate:"required"`
}

type LoggingConfig struct {
	Level  string `toml:"level" validate:"required,oneof=debug info warning error critical"`
	Trace  bool   `toml:"trace"`
	Buffer int64  `toml:"buffer" validate:"gte=0,lte=4096"`
}

type DatabaseConfig struct {
	Addr        string `toml:"addr" validate:"required,tcp_addr|unix_addr"`
	Username    string `toml:"username" validate:"required"`
	Password    string `toml:"password" validate:"required"`
	TimeoutTTL  int64  `toml:"timeout_ttl" validate:"required"`
	Healthcheck bool   `toml:"healthcheck"`
	SkipVerify  bool   `toml:"skipverify"`
}

type SchedulerConfig struct {
	CycleTTL int64 `toml:"cycle_ttl" validate:"required"`
  RescheduleCycles int64 `toml:"reschedule_cycles" validate:"required"`
	CpuThreshold int64 `toml:"cpu_threshold" validate:"gte=0,lte=1000"`
	MemThreshold int64 `toml:"mem_threshold" validate:"gte=0,lte=1000"`
}

type ApiConfig struct {
	Addr     string `toml:"addr" validate:"required,tcp_addr"`
	CertFile string `toml:"cert_file" validate:"required"`
	KeyFile  string `toml:"key_file" validate:"required"`
}

// LoadConfig reads the configuration file, decodes it (toml) and validates it.
// If a step fails, an error is returned.
func LoadConfig(path string) (*BaseConfig, error) {
	rawConfig, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := &BaseConfig{}
	_, err = toml.Decode(string(rawConfig), config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
