// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lsp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sync/atomic"
	"time"

	"golang.org/x/exp/jsonrpc2"
)

var (
	// RequestCancelledError should be used when a request is cancelled early.
	RequestCancelledError = jsonrpc2.NewError(-32800, "JSON RPC cancelled")
)

// detach returns a context that keeps all the values of its parent context
// but detaches from the cancellation and error handling.
func detach(ctx context.Context) context.Context { return detachedContext{ctx} }

type detachedContext struct{ parent context.Context }

func (v detachedContext) Deadline() (time.Time, bool) { return time.Time{}, false }
func (v detachedContext) Done() <-chan struct{}       { return nil }
func (v detachedContext) Err() error                  { return nil }
func (v detachedContext) Value(key any) any           { return v.parent.Value(key) }

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
		detached := detach(ctx)
		_ = c.conn.Notify(detached, "$/cancelRequest", &CancelParams{ID: cancelParamsFromID(call.ID())})
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
		result, err := clientDispatch(ctx, client, req)
		return nullResult(req, result, err)
	}
}

// ServerHandler returns a jsonrpc2 handler that dispatches incoming LSP
// messages to server.
//
// The handler tracks the shutdown state of the connection it serves, so a
// separate handler must be created for every connection. Once a shutdown
// request has been handled successfully, all further requests are answered
// with an InvalidRequest error and all notifications except exit are dropped,
// as mandated by the LSP specification.
func ServerHandler(server Server) jsonrpc2.HandlerFunc {
	var shutdown atomic.Bool
	return func(ctx context.Context, req *jsonrpc2.Request) (any, error) {
		if shutdown.Load() && req.Method != "exit" {
			if req.IsCall() {
				return nil, fmt.Errorf("%w: %q received after shutdown", jsonrpc2.ErrInvalidRequest, req.Method)
			}
			return nil, nil
		}
		if ctx.Err() != nil {
			return nil, RequestCancelledError
		}
		result, err := serverDispatch(ctx, server, req)
		if req.Method == "shutdown" && err == nil {
			shutdown.Store(true)
		}
		return nullResult(req, result, err)
	}
}

// nullResult turns the missing result of a successfully handled request
// (such as shutdown) into a JSON null. jsonrpc2 would otherwise reject the
// missing result and respond with an internal error.
func nullResult(req *jsonrpc2.Request, result any, err error) (any, error) {
	if result == nil && err == nil && req.IsCall() {
		return json.RawMessage("null"), nil
	}
	return result, err
}

func Call(ctx context.Context, conn *jsonrpc2.Connection, method string, params any, result any) error {
	call := conn.Call(ctx, method, params)
	err := call.Await(ctx, result)
	if ctx.Err() != nil {
		_ = conn.Notify(detach(ctx), "$/cancelRequest", &CancelParams{ID: cancelParamsFromID(call.ID())})
	}
	return err
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

func cancelParamsFromID(id jsonrpc2.ID) CancelParamsId {
	switch v := id.Raw().(type) {
	case int32:
		return CancelParamsId{Int32: &v}
	case string:
		return CancelParamsId{String: &v}
	default:
		// Invalid, return empty value
		return CancelParamsId{}
	}
}
