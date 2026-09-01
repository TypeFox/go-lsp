// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lsp

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sync"
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

// CancelHandler returns a jsonrpc2.Preempter and jsonrpc2.Handler pair that
// implement LSP '$/cancelRequest' support around the given handler. Install
// both on the same connection:
//
//	preempter, handler := lsp.CancelHandler(lsp.ServerHandler(server))
//	opts := jsonrpc2.ConnectionOptions{Preempter: preempter, Handler: handler}
//
// The preempter intercepts '$/cancelRequest' notifications before they are
// queued and cancels the context of the matching in-flight request, which then
// responds with RequestCancelledError (-32800). A cancellation that arrives
// before its request has been delivered to the handler is a no-op; the request
// runs to completion (best effort, as permitted by the LSP spec).
//
// When combined with AsyncBinder, wrap with CancelHandler first, then pass the
// resulting options to AsyncBinder; see the AsyncBinder documentation for the
// full example and why the nesting order matters.
func CancelHandler(handler jsonrpc2.Handler) (jsonrpc2.Preempter, jsonrpc2.Handler) {
	c := &canceller{handler: handler, pending: make(map[jsonrpc2.ID]context.CancelFunc)}
	return c, c
}

type canceller struct {
	handler jsonrpc2.Handler
	mu      sync.Mutex
	pending map[jsonrpc2.ID]context.CancelFunc
}

func (c *canceller) Preempt(ctx context.Context, req *jsonrpc2.Request) (any, error) {
	if req.Method != "$/cancelRequest" {
		return nil, jsonrpc2.ErrNotHandled
	}
	var params CancelParams
	if err := UnmarshalJSON(req.Params, &params); err != nil {
		return nil, fmt.Errorf("%w: %v", jsonrpc2.ErrParse, err)
	}
	var id jsonrpc2.ID
	switch {
	case params.ID.Int32 != nil:
		id = jsonrpc2.Int64ID(int64(*params.ID.Int32))
	case params.ID.String != nil:
		id = jsonrpc2.StringID(*params.ID.String)
	default:
		return nil, fmt.Errorf("%w: $/cancelRequest without a request ID", jsonrpc2.ErrInvalidParams)
	}
	c.mu.Lock()
	cancel := c.pending[id]
	c.mu.Unlock()
	if cancel != nil {
		cancel()
	}
	return nil, nil
}

func (c *canceller) Handle(ctx context.Context, req *jsonrpc2.Request) (any, error) {
	if !req.IsCall() {
		return c.handler.Handle(ctx, req)
	}
	ctx, cancel := context.WithCancel(ctx)
	c.mu.Lock()
	c.pending[req.ID] = cancel
	c.mu.Unlock()
	defer func() {
		c.mu.Lock()
		delete(c.pending, req.ID)
		c.mu.Unlock()
		cancel()
	}()
	result, err := c.handler.Handle(ctx, req)
	if err != nil && ctx.Err() != nil {
		return nil, RequestCancelledError
	}
	return result, err
}

// AsyncBinder returns a jsonrpc2.Binder that installs the given options with
// their Handler wrapped so that each request (call) runs in its own goroutine
// and replies asynchronously, while notifications remain synchronous and are
// processed strictly in arrival order. A notification such as
// 'textDocument/didChange' is therefore always applied before any request that
// arrives after it is handled.
//
// Because requests may then run concurrently with each other and with
// notifications, the wrapped handler must synchronize access to shared state
// itself (e.g. with a sync.RWMutex around the workspace).
//
// To combine async dispatch with '$/cancelRequest' support, wrap the handler
// with CancelHandler first and give the resulting options to AsyncBinder:
//
//	preempter, handler := lsp.CancelHandler(lsp.ServerHandler(server))
//	binder := lsp.AsyncBinder(jsonrpc2.ConnectionOptions{
//		Preempter: preempter,
//		Handler:   handler,
//	})
//	conn, err := jsonrpc2.Dial(ctx, dialer, binder)
func AsyncBinder(opts jsonrpc2.ConnectionOptions) jsonrpc2.Binder {
	return asyncBinder{opts: opts}
}

type asyncBinder struct {
	opts jsonrpc2.ConnectionOptions
}

func (b asyncBinder) Bind(ctx context.Context, conn *jsonrpc2.Connection) (jsonrpc2.ConnectionOptions, error) {
	opts := b.opts
	if opts.Handler != nil {
		opts.Handler = &asyncHandler{conn: conn, inner: opts.Handler}
	}
	return opts, nil
}

type asyncHandler struct {
	conn  *jsonrpc2.Connection
	inner jsonrpc2.Handler
}

func (h *asyncHandler) Handle(ctx context.Context, req *jsonrpc2.Request) (any, error) {
	if !req.IsCall() {
		return h.inner.Handle(ctx, req)
	}
	go func() {
		result, err := h.inner.Handle(ctx, req)
		if errors.Is(err, jsonrpc2.ErrNotHandled) {
			// deliverMessages would normally do this conversion; the async
			// path bypasses it.
			result, err = nil, fmt.Errorf("%w: %q", jsonrpc2.ErrMethodNotFound, req.Method)
		}
		_ = h.conn.Respond(req.ID, result, err)
	}()
	return nil, jsonrpc2.ErrAsyncResponse
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
	case int64:
		// jsonrpc2 stores integer IDs as int64 (see jsonrpc2.Int64ID).
		v32 := int32(v)
		return CancelParamsId{Int32: &v32}
	case string:
		return CancelParamsId{String: &v}
	default:
		// Invalid, return empty value
		return CancelParamsId{}
	}
}
