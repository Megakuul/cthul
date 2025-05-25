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

package domain

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"cthul.io/cthul/pkg/api/wave/v1/domain"
	domctrl "cthul.io/cthul/pkg/wave/domain"
	"github.com/google/uuid"
)

type Service struct {
	controller *domctrl.Controller
}

func New(controller *domctrl.Controller) *Service {
	return &Service{
		controller: controller,
	}
}

func (d *Service) Get(ctx context.Context, r *connect.Request[domain.GetRequest]) (*connect.Response[domain.GetResponse], error) {
	// TODO: authorize
	result, err := d.controller.Lookup(ctx, r.Msg.Id)
	if err != nil {
		var mismatchErr *domctrl.NodeMismatchErr
		if errors.As(err, &mismatchErr) {
			rpcErr := connect.NewError(connect.CodeNotFound, mismatchErr)
			rpcErr.Meta().Add("Location", mismatchErr.Node)
			return nil, rpcErr
		}
		return nil, err
	}

	return &connect.Response[domain.GetResponse]{
		Msg: &domain.GetResponse{Domain: result},
	}, nil
}

func (d *Service) Stat(ctx context.Context, r *connect.Request[domain.StatRequest]) (*connect.Response[domain.StatResponse], error) {
	// TODO: authorize
	result, err := d.controller.Stat(ctx, r.Msg.Id)
	if err != nil {
		var mismatchErr *domctrl.NodeMismatchErr
		if errors.As(err, &mismatchErr) {
			rpcErr := connect.NewError(connect.CodeNotFound, mismatchErr)
			rpcErr.Meta().Add("Location", mismatchErr.Node)
			return nil, rpcErr
		}
		return nil, err
	}

	return &connect.Response[domain.StatResponse]{
		Msg: &domain.StatResponse{Stats: result},
	}, nil
}

func (d *Service) List(ctx context.Context, r *connect.Request[domain.ListRequest]) (*connect.Response[domain.ListResponse], error) {
	// TODO: authorize
	result, err := d.controller.List(ctx)
	if err != nil {
		var mismatchErr *domctrl.NodeMismatchErr
		if errors.As(err, &mismatchErr) {
			rpcErr := connect.NewError(connect.CodeNotFound, mismatchErr)
			rpcErr.Meta().Add("Location", mismatchErr.Node)
			return nil, rpcErr
		}
		return nil, err
	}

	return &connect.Response[domain.ListResponse]{
		Msg: &domain.ListResponse{Domains: result},
	}, nil
}

func (d *Service) Create(ctx context.Context, r *connect.Request[domain.CreateRequest]) (*connect.Response[domain.CreateResponse], error) {
	// TODO: authorize
	id := uuid.New().String()
	err := d.controller.Apply(ctx, id, r.Msg.Config)
	if err != nil {
		var mismatchErr *domctrl.NodeMismatchErr
		if errors.As(err, &mismatchErr) {
			rpcErr := connect.NewError(connect.CodeNotFound, mismatchErr)
			rpcErr.Meta().Add("Location", mismatchErr.Node)
			return nil, rpcErr
		}
		return nil, err
	}

	return &connect.Response[domain.CreateResponse]{
		Msg: &domain.CreateResponse{Id: id},
	}, nil
}

func (d *Service) Update(ctx context.Context, r *connect.Request[domain.UpdateRequest]) (*connect.Response[domain.UpdateResponse], error) {
	// TODO: authorize
	err := d.controller.Apply(ctx, r.Msg.Id, r.Msg.Config)
	if err != nil {
		var mismatchErr *domctrl.NodeMismatchErr
		if errors.As(err, &mismatchErr) {
			rpcErr := connect.NewError(connect.CodeNotFound, mismatchErr)
			rpcErr.Meta().Add("Location", mismatchErr.Node)
			return nil, rpcErr
		}
		return nil, err
	}

	return &connect.Response[domain.UpdateResponse]{
		Msg: &domain.UpdateResponse{},
	}, nil
}

func (d *Service) Attach(ctx context.Context, r *connect.Request[domain.AttachRequest]) (*connect.Response[domain.AttachResponse], error) {
	// TODO: authorize
	err := d.controller.Attach(ctx, r.Msg.Id, r.Msg.Node, false)
	if err != nil {
		var mismatchErr *domctrl.NodeMismatchErr
		if errors.As(err, &mismatchErr) {
			rpcErr := connect.NewError(connect.CodeNotFound, mismatchErr)
			rpcErr.Meta().Add("Location", mismatchErr.Node)
			return nil, rpcErr
		}
		return nil, err
	}

	return &connect.Response[domain.AttachResponse]{
		Msg: &domain.AttachResponse{},
	}, nil
}

func (d *Service) Detach(ctx context.Context, r *connect.Request[domain.DetachRequest]) (*connect.Response[domain.DetachResponse], error) {
	// TODO: authorize
	err := d.controller.Detach(ctx, r.Msg.Id)
	if err != nil {
		var mismatchErr *domctrl.NodeMismatchErr
		if errors.As(err, &mismatchErr) {
			rpcErr := connect.NewError(connect.CodeNotFound, mismatchErr)
			rpcErr.Meta().Add("Location", mismatchErr.Node)
			return nil, rpcErr
		}
		return nil, err
	}

	return &connect.Response[domain.DetachResponse]{
		Msg: &domain.DetachResponse{},
	}, nil
}

func (d *Service) Delete(ctx context.Context, r *connect.Request[domain.DeleteRequest]) (*connect.Response[domain.DeleteResponse], error) {
	// TODO: authorize
	err := d.controller.Delete(ctx, r.Msg.Id)
	if err != nil {
		var mismatchErr *domctrl.NodeMismatchErr
		if errors.As(err, &mismatchErr) {
			rpcErr := connect.NewError(connect.CodeNotFound, mismatchErr)
			rpcErr.Meta().Add("Location", mismatchErr.Node)
			return nil, rpcErr
		}
		return nil, err
	}

	return &connect.Response[domain.DeleteResponse]{
		Msg: &domain.DeleteResponse{},
	}, nil
}
