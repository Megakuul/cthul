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

package resource

import (
	"cthul.io/cthul/pkg/db"
	"cthul.io/cthul/pkg/log"
	"cthul.io/cthul/pkg/log/discard"
)

// ResourceOperator provides operations to fetch and update node & domain resources.
type ResourceOperator struct {
	client db.Client
	logger log.Logger
}

type ResourceOperatorOption func(*ResourceOperator)

// NewResourceOperator creates a new resource operator.
func NewResourceOperator(client db.Client, opts ...ResourceOperatorOption) *ResourceOperator {
	operator := &ResourceOperator{
		client: client,
		logger: discard.NewDiscardLogger(),
	}

	for _, opt := range opts {
		opt(operator)
	}

	return operator
}

// WithLogger enables logging for the operator.
func WithLogger(logger log.Logger) ResourceOperatorOption {
	return func(r *ResourceOperator) {
		r.logger = logger
	}
}
