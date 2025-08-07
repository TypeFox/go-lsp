// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package protocol

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/TypeFox/go-lsp/internal/event"
	"golang.org/x/exp/jsonrpc2"
	"github.com/TypeFox/go-lsp/internal/xcontext"
)

var (
	// RequestCancelledError should be used when a request is cancelled early.
	RequestCancelledError = jsonrpc2.NewError(-32800, "JSON RPC cancelled")
)

type ClientCloser interface {
	Client
	io.Closer
}

type connSender interface {
	io.Closer

	Notify(ctx context.Context, method string, params any) error
	Call(ctx context.Context, method string, params, result any) error
}

type clientDispatcher struct {
	sender connSender
}

func (c *clientDispatcher) Close() error {
	return c.sender.Close()
}

// ClientDispatcher returns a Client that dispatches LSP requests across the
// given jsonrpc2 connection.
func ClientDispatcher(conn *jsonrpc2.Connection) ClientCloser {
	return &clientDispatcher{sender: clientConn{conn}}
}

type clientConn struct {
	conn *jsonrpc2.Connection
}

func (c clientConn) Close() error {
	return c.conn.Close()
}

func (c clientConn) Notify(ctx context.Context, method string, params any) error {
	return c.conn.Notify(ctx, method, params)
}

func (c clientConn) Call(ctx context.Context, method string, params any, result any) error {
	call := c.conn.Call(ctx, method, params)
	err := call.Await(ctx, result)
	if ctx.Err() != nil {
		detached := xcontext.Detach(ctx)
		c.conn.Notify(detached, "$/cancelRequest", &CancelParams{ID: call.ID().Raw()})
	}
	return err
}

// ServerDispatcher returns a Server that dispatches LSP requests across the
// given jsonrpc2 connection.
func ServerDispatcher(conn *jsonrpc2.Connection) Server {
	return &serverDispatcher{sender: clientConn{conn}}
}

type serverDispatcher struct {
	sender connSender
}

func ClientHandler(client Client) jsonrpc2.HandlerFunc {
	return func(ctx context.Context, req *jsonrpc2.Request) (any, error) {
		if ctx.Err() != nil {
			return nil, RequestCancelledError
		}
		return clientDispatch(ctx, client, req)
	}
}

func ServerHandler(server Server) jsonrpc2.HandlerFunc {
	return func(ctx context.Context, req *jsonrpc2.Request) (any, error) {
		if ctx.Err() != nil {
			return nil, RequestCancelledError
		}
		return serverDispatch(ctx, server, req)
	}
}

func Call(ctx context.Context, conn *jsonrpc2.Connection, method string, params any, result any) error {
	call := conn.Call(ctx, method, params)
	err := call.Await(ctx, result)
	if ctx.Err() != nil {
		conn.Notify(xcontext.Detach(ctx), "$/cancelRequest", &CancelParams{ID: call.ID().Raw()})
	}
	return err
}

func cancelCall(ctx context.Context, sender connSender, id any) {
	if ctx.Err() == nil {
		return
	}
	ctx = xcontext.Detach(ctx)
	// Note that only *jsonrpc2.ID implements json.Marshaler.
	sender.Notify(ctx, "$/cancelRequest", &CancelParams{ID: id})
}

func writeError(ctx context.Context, err error) error {
	if err == nil {
		event.Error(ctx, "jsonrpc2 internal error", fmt.Errorf("null error"))
		err = jsonrpc2.ErrInternal
	}
	event.Error(ctx, "jsonrpc2 error", err)
	return err
}

func sendParseError(ctx context.Context, err error) error {
	return writeError(ctx, fmt.Errorf("%w: %s", jsonrpc2.ErrParse, err))
}

// UnmarshalJSON unmarshals msg into the variable pointed to by
// params. In JSONRPC, optional messages may be
// "null", in which case it is a no-op.
func UnmarshalJSON(msg json.RawMessage, v any) error {
	if len(msg) == 0 || bytes.Equal(msg, []byte("null")) {
		return nil
	}
	return json.Unmarshal(msg, v)
}

// NonNilSlice returns x, or an empty slice if x was nil.
//
// (Many slice fields of protocol structs must be non-nil
// to avoid being encoded as JSON "null".)
func NonNilSlice[T comparable](x []T) []T {
	if x == nil {
		return []T{}
	}
	return x
}

// EncodeMessage is a utility to encode LSP protocol messages to JSON.
func EncodeMessage(msg any) ([]byte, error) {
	data, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal message: %w", err)
	}
	
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("Content-Length: %d\r\n\r\n", len(data)))
	buf.Write(data)
	return buf.Bytes(), nil
}
