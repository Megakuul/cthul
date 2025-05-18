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

package node

import (
	"connectrpc.com/connect"
	"context"
	"cthul.io/cthul/pkg/api/wave/v1/node"
	nodectrl "cthul.io/cthul/pkg/wave/node"
	"errors"
)

type Service struct {
	controller *nodectrl.Controller
}

func New(controller *nodectrl.Controller) *Service {
	return &Service{
    controller: controller,
  }
}

func (d *Service) Get(ctx context.Context, r *connect.Request[node.GetRequest]) (*connect.Response[node.GetResponse], error) {
  // TODO: authorize
	result, err := d.controller.Lookup(ctx, r.Msg.Id)
	if err != nil {
		var mismatchErr *nodectrl.NodeMismatchErr
		if errors.As(err, &mismatchErr) {
			rpcErr := connect.NewError(connect.CodeNotFound, mismatchErr)
			rpcErr.Meta().Add("Location", mismatchErr.Node)
			return nil, rpcErr
		}
		return nil, err
	}

	return &connect.Response[node.GetResponse]{
		Msg: &node.GetResponse{Node: result},
	}, nil
}

func (d *Service) List(ctx context.Context, r *connect.Request[node.ListRequest]) (*connect.Response[node.ListResponse], error) {
  // TODO: authorize
	result, err := d.controller.List(ctx)
	if err != nil {
		var mismatchErr *nodectrl.NodeMismatchErr
		if errors.As(err, &mismatchErr) {
			rpcErr := connect.NewError(connect.CodeNotFound, mismatchErr)
			rpcErr.Meta().Add("Location", mismatchErr.Node)
			return nil, rpcErr
		}
		return nil, err
	}

	return &connect.Response[node.ListResponse]{
		Msg: &node.ListResponse{Nodes: result},
	}, nil
}
