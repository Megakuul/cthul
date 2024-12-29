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

package node

type NodeOperator struct {
	rootCtx       context.Context
	rootCtxCancel context.CancelFunc

	workCtx       context.Context
	workCtxCancel context.CancelFunc

	// finChan is used to send the absolute exist signal
	// if the channel emits, this indicates that the operator is fully cleaned up.
	finChan chan struct{}

	client db.Client
	logger log.Logger

	// cycleTTL specifies the interval for scheduler cycles.
	cycleTTL int64
}

type NodeOperatorOption func(*NodeOperator)

func NewNodeOperator(opts ...NodeOperatorOption) *NodeOperator {
	operator := &NodeOperator{

	}

	for _, opt := range opts {
		opt(operator)
	}

	return operator
}

