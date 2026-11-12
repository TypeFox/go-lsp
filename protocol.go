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

// CancelBinder returns a [jsonrpc2.Binder] that adds LSP '$/cancelRequest'
// support to the connection options produced by the given binder. Since
// [jsonrpc2.ConnectionOptions] is itself a binder, the simple case is:
//
//	binder := lsp.CancelBinder(jsonrpc2.ConnectionOptions{
//		Handler: lsp.ServerHandler(server),
//	})
//	conn, err := jsonrpc2.Dial(ctx, dialer, binder)
//
// The binder installs a preempter that intercepts '$/cancelRequest'
// notifications before they are queued and cancels the matching request through
// [jsonrpc2.Connection.Cancel]. The connection tracks every request from the
// moment it is read, so a request can be cancelled both while it is being
// handled and while it is still waiting in the queue. An in-flight request has
// its context cancelled and a queued request is never passed to the handler.
// Either way the request responds with RequestCancelledError (-32800). A
// cancellation for an unknown or already answered request is a no-op.
//
// A Preempter already present in the options still sees every other message.
//
// CancelBinder composes with [AsyncBinder] in either order.
func CancelBinder(binder jsonrpc2.Binder) jsonrpc2.Binder {
	return cancelBinder{inner: binder}
}

type cancelBinder struct {
	inner jsonrpc2.Binder
}

func (b cancelBinder) Bind(ctx context.Context, conn *jsonrpc2.Connection) (jsonrpc2.ConnectionOptions, error) {
	opts, err := b.inner.Bind(ctx, conn)
	if err != nil {
		return opts, err
	}
	if opts.Framer == nil {
		opts.Framer = jsonrpc2.HeaderFramer()
	}
	opts.Framer = cancelFramer{opts.Framer}
	opts.Preempter = cancelPreempter{conn: conn, next: opts.Preempter}
	return opts, nil
}

// cancelPreempter handles '$/cancelRequest' and passes every other message on
// to the next preempter, if any.
type cancelPreempter struct {
	conn *jsonrpc2.Connection
	next jsonrpc2.Preempter
}

func (p cancelPreempter) Preempt(ctx context.Context, req *jsonrpc2.Request) (any, error) {
	if req.Method != "$/cancelRequest" {
		if p.next != nil {
			return p.next.Preempt(ctx, req)
		}
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
	p.conn.Cancel(id)
	return nil, nil
}

// cancelFramer wraps a framer so that responses failing with context.Canceled
// are sent as RequestCancelledError.
//
// This translation lives in the framing layer because it is the only hook the
// connection offers for it: a request cancelled while still queued is answered
// by the connection itself with the raw context error, without ever reaching a
// handler that could map it.
type cancelFramer struct {
	jsonrpc2.Framer
}

func (f cancelFramer) Writer(w io.Writer) jsonrpc2.Writer {
	return cancelWriter{f.Framer.Writer(w)}
}

type cancelWriter struct {
	jsonrpc2.Writer
}

func (w cancelWriter) Write(ctx context.Context, msg jsonrpc2.Message) (int64, error) {
	if resp, ok := msg.(*jsonrpc2.Response); ok && errors.Is(resp.Error, context.Canceled) {
		resp.Error = RequestCancelledError
	}
	return w.Writer.Write(ctx, msg)
}

// AsyncBinder returns a [jsonrpc2.Binder] that wraps the Handler of the
// connection options produced by the given binder, so that each request (call)
// runs in its own goroutine and replies asynchronously, while notifications
// remain synchronous and are processed strictly in arrival order. A
// notification such as 'textDocument/didChange' is therefore always applied
// before any request that arrives after it is handled.
//
// Because requests may then run concurrently with each other and with
// notifications, the wrapped handler must synchronize access to shared state
// itself (e.g. with a sync.RWMutex around the workspace).
//
// To combine async dispatch with '$/cancelRequest' support, nest it with
// [CancelBinder]:
//
//	binder := lsp.CancelBinder(lsp.AsyncBinder(jsonrpc2.ConnectionOptions{
//		Handler: lsp.ServerHandler(server),
//	}))
//	conn, err := jsonrpc2.Dial(ctx, dialer, binder)
func AsyncBinder(binder jsonrpc2.Binder) jsonrpc2.Binder {
	return asyncBinder{inner: binder}
}

type asyncBinder struct {
	inner jsonrpc2.Binder
}

func (b asyncBinder) Bind(ctx context.Context, conn *jsonrpc2.Connection) (jsonrpc2.ConnectionOptions, error) {
	opts, err := b.inner.Bind(ctx, conn)
	if err != nil {
		return opts, err
	}
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
		// notifications are deliberately handled inline to preserve their arrival order.
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
