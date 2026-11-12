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
	return Call(ctx, c.conn, method, params, result)
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
	opts.Framer = cancelFramer{Framer: opts.Framer}
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
	return cancelWriter{Writer: f.Framer.Writer(w)}
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

// AsyncBinder returns a [jsonrpc2.Binder] that runs each request (call) of
// the given binder's Handler in its own goroutine. Notifications are still
// handled synchronously in arrival order, and each one waits for all earlier
// requests to respond. A request therefore sees exactly the notifications
// that arrived before it, and the handler only needs to synchronize state
// shared between concurrent requests.
//
// A slow request delays every notification behind it, so request handlers
// must return promptly once their context is cancelled and must not wait for
// a notification from the client. The context is also cancelled when the
// connection's stream ends.
//
// Nest it with [CancelBinder] so that the client can cancel outdated requests:
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
		// Cancel the running requests as soon as the stream ends. The
		// connection cannot finish shutting down while a queued notification
		// is waiting for them.
		closed := make(chan struct{})
		if opts.Framer == nil {
			opts.Framer = jsonrpc2.HeaderFramer()
		}
		opts.Framer = closeFramer{
			Framer:  opts.Framer,
			onClose: sync.OnceFunc(func() { close(closed) }),
		}
		opts.Handler = &asyncHandler{conn: conn, inner: opts.Handler, closed: closed}
	}
	return opts, nil
}

// closeFramer wraps a framer to report when reading from the stream fails,
// which means the connection is shutting down.
type closeFramer struct {
	jsonrpc2.Framer
	onClose func()
}

func (f closeFramer) Reader(r io.Reader) jsonrpc2.Reader {
	return closeReader{Reader: f.Framer.Reader(r), onClose: f.onClose}
}

type closeReader struct {
	jsonrpc2.Reader
	onClose func()
}

func (r closeReader) Read(ctx context.Context) (jsonrpc2.Message, int64, error) {
	msg, n, err := r.Reader.Read(ctx)
	if err != nil {
		r.onClose()
	}
	return msg, n, err
}

type asyncHandler struct {
	conn  *jsonrpc2.Connection
	inner jsonrpc2.Handler
	// closed is closed once the connection's stream has ended.
	closed <-chan struct{}
	// requests tracks the running requests. It is only added to and waited
	// on from Handle, which the connection calls from a single goroutine.
	requests sync.WaitGroup
}

func (h *asyncHandler) Handle(ctx context.Context, req *jsonrpc2.Request) (any, error) {
	if !req.IsCall() {
		// A notification is handled inline, after every earlier request has
		// responded, so that messages take effect in arrival order.
		h.requests.Wait()
		return h.inner.Handle(ctx, req)
	}
	ctx, cancel := context.WithCancel(ctx)
	// Cancel the request when the stream ends. This goroutine exits when the
	// request has returned, at the latest.
	go func() {
		select {
		case <-h.closed:
			cancel()
		case <-ctx.Done():
		}
	}()
	h.requests.Go(func() {
		defer cancel()
		result, err := h.inner.Handle(ctx, req)
		if errors.Is(err, jsonrpc2.ErrNotHandled) {
			// deliverMessages would normally do this conversion; the async
			// path bypasses it.
			result, err = nil, fmt.Errorf("%w: %q", jsonrpc2.ErrMethodNotFound, req.Method)
		}
		// Respond only fails if the result cannot be marshalled or the
		// stream is broken, and there is nobody left to report that to.
		_ = h.conn.Respond(req.ID, result, err)
	})
	return nil, jsonrpc2.ErrAsyncResponse
}

// Call invokes method on conn and waits for the result. If ctx is done before
// the response arrives, the request is cancelled with a '$/cancelRequest'
// notification.
func Call(ctx context.Context, conn *jsonrpc2.Connection, method string, params any, result any) error {
	call := conn.Call(ctx, method, params)
	err := call.Await(ctx, result)
	if ctx.Err() != nil {
		// Best effort: the call has already failed with the context error.
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
