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
	"os"
	"github.com/BurntSushi/toml"
	"github.com/go-playground/validator/v10"
)

// BaseConfig represents the root configuration.
type BaseConfig struct {
	Logging LoggingConfig `toml:"logging"`
}

// LoggingConfig defines log related configuration.
type LoggingConfig struct {
	Level string `toml:"level" validate:"required,oneof=debug info warning error critical"`
	Trace bool `toml:"trace"`
	Buffer int64 `toml:"buffer" validate:"gte=0,lte=4096"`
}

// LoadConfig reads the configuration file, decodes it (toml) and validates it.
// If a step fails, an error is returned.
func LoadConfig(path string) (*BaseConfig, error) {
	rawConfig, err := os.ReadFile(path)
	if err!=nil {
		return nil, err
	}

	config := &BaseConfig{}
	_, err = toml.Decode(string(rawConfig), config)
	if err!=nil {
		return nil, err
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(config)
	if err!=nil {
		return nil, err
	}

	return config, nil
}
