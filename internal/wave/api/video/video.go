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
	"errors"
	"io"

	"connectrpc.com/connect"
	"cthul.io/cthul/pkg/api/wave/v1/video"
	videoctrl "cthul.io/cthul/pkg/wave/video"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

type Service struct {
	controller *videoctrl.Controller
}

func New(controller *videoctrl.Controller) *Service {
	return &Service{
		controller: controller,
	}
}

func (d *Service) Get(ctx context.Context, r *connect.Request[video.GetRequest]) (*connect.Response[video.GetResponse], error) {
	// TODO: authorize
	result, err := d.controller.Lookup(ctx, r.Msg.Id)
	if err != nil {
		var mismatchErr *videoctrl.NodeMismatchErr
		if errors.As(err, &mismatchErr) {
			rpcErr := connect.NewError(connect.CodeNotFound, mismatchErr)
			rpcErr.Meta().Add("Location", mismatchErr.Node)
			return nil, rpcErr
		}
		return nil, err
	}

	return &connect.Response[video.GetResponse]{
		Msg: &video.GetResponse{Video: result},
	}, nil
}


func (d *Service) Connect(ctx context.Context, r *connect.BidiStream[video.ConnectRequest, video.ConnectResponse]) error {
	// TODO: authorize
	reader, writer := make(chan<- []byte), make(<-chan []byte)

	err := d.controller.Connect(ctx, r.RequestHeader().Get("id"), reader, writer)
	if err != nil {
		var mismatchErr *videoctrl.NodeMismatchErr
		if errors.As(err, &mismatchErr) {
			rpcErr := connect.NewError(connect.CodeNotFound, mismatchErr)
			rpcErr.Meta().Add("Location", mismatchErr.Node)
			return rpcErr
		}
		return err
	}

	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		errChan := make(chan error)
		// wrapped into goroutine because r.Receive() is not manually cancellable.
		// -> it is automatically cleaned up after Connect() exits.
		go func() {
			for {
				input, err := r.Receive()
				if err != nil {
					errChan <- err
					return
				}
				reader <- input.Input
			}
		}()
		select {
		case err := <-errChan:
			return err
		case <-gCtx.Done():
			return nil
		}
	})

	g.Go(func() error {
		for {
			select {
			case output := <-writer:
				err := r.Send(&video.ConnectResponse{
					Output: output,
				})
				if err != nil {
					if errors.Is(err, io.EOF) {
						return nil
					}
					return err
				}
			case <-gCtx.Done():
				return nil
			}
		}
	})

	return g.Wait()
}

func (d *Service) List(ctx context.Context, r *connect.Request[video.ListRequest]) (*connect.Response[video.ListResponse], error) {
	// TODO: authorize
	result, err := d.controller.List(ctx)
	if err != nil {
		var mismatchErr *videoctrl.NodeMismatchErr
		if errors.As(err, &mismatchErr) {
			rpcErr := connect.NewError(connect.CodeNotFound, mismatchErr)
			rpcErr.Meta().Add("Location", mismatchErr.Node)
			return nil, rpcErr
		}
		return nil, err
	}

	return &connect.Response[video.ListResponse]{
		Msg: &video.ListResponse{Videos: result},
	}, nil
}

func (d *Service) Create(ctx context.Context, r *connect.Request[video.CreateRequest]) (*connect.Response[video.CreateResponse], error) {
	// TODO: authorize
	id := uuid.New().String()
	err := d.controller.Apply(ctx, id, r.Msg.Config)
	if err != nil {
		var mismatchErr *videoctrl.NodeMismatchErr
		if errors.As(err, &mismatchErr) {
			rpcErr := connect.NewError(connect.CodeNotFound, mismatchErr)
			rpcErr.Meta().Add("Location", mismatchErr.Node)
			return nil, rpcErr
		}
		return nil, err
	}

	return &connect.Response[video.CreateResponse]{
		Msg: &video.CreateResponse{Id: id},
	}, nil
}

func (d *Service) Update(ctx context.Context, r *connect.Request[video.UpdateRequest]) (*connect.Response[video.UpdateResponse], error) {
	// TODO: authorize
	err := d.controller.Apply(ctx, r.Msg.Id, r.Msg.Config)
	if err != nil {
		var mismatchErr *videoctrl.NodeMismatchErr
		if errors.As(err, &mismatchErr) {
			rpcErr := connect.NewError(connect.CodeNotFound, mismatchErr)
			rpcErr.Meta().Add("Location", mismatchErr.Node)
			return nil, rpcErr
		}
		return nil, err
	}

	return &connect.Response[video.UpdateResponse]{
		Msg: &video.UpdateResponse{},
	}, nil
}

func (d *Service) Delete(ctx context.Context, r *connect.Request[video.DeleteRequest]) (*connect.Response[video.DeleteResponse], error) {
	// TODO: authorize
	err := d.controller.Delete(ctx, r.Msg.Id)
	if err != nil {
		var mismatchErr *videoctrl.NodeMismatchErr
		if errors.As(err, &mismatchErr) {
			rpcErr := connect.NewError(connect.CodeNotFound, mismatchErr)
			rpcErr.Meta().Add("Location", mismatchErr.Node)
			return nil, rpcErr
		}
		return nil, err
	}

	return &connect.Response[video.DeleteResponse]{
		Msg: &video.DeleteResponse{},
	}, nil
}
