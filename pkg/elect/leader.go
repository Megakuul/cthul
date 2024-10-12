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

package elect

import (
	"encoding/json"
	"fmt"
)

// clusterLeader contains information about the cluster leader node.
type clusterLeader struct {
	Id   string `json:"id"`
	Cash int64  `json:"cash"`
}

// parseClusterLeader parses the leader string into a cluster leader.
func parseClusterLeader(leaderStr string) (*clusterLeader, error) {
	leader := clusterLeader{}
	err := json.Unmarshal([]byte(leaderStr), &leader)
	if err != nil {
		return nil, fmt.Errorf("cannot parse node information")
	}

	return &leader, nil
}

// serializeClusterLeader serializes the cluster leader into a string.
func serializeClusterLeader(leader *clusterLeader) string {
	leaderStr, err := json.Marshal(leader)
	if err != nil {
		return ""
	}
	return string(leaderStr)
}
