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
	"github.com/spf13/cobra"
)

type cliFlags struct {
	configPath string
}

func NewRootCmd() *cobra.Command {
	flags := &cliFlags{}
	cmd := &cobra.Command{
		Use:          "wave",
		SilenceUsage: true,
		Long: `wave is the core domain controller; it manages the virtual machines of the local node.
The service is installed on every cthul node that hosts virtual machines and directly communicates with the local libvirtd service.`,

		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := LoadConfig(flags.configPath)
			if err!=nil {
				return err
			}

			return Run(config)
		},
	}

	cmd.PersistentFlags().StringVarP(&flags.configPath,
		"config", "c", "config.toml", "path of the configuration file")

	return cmd
}
