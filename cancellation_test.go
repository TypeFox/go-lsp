// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lsp_test

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"golang.org/x/exp/jsonrpc2"
	"typefox.dev/lsp"
)

// testConn serves binder on an in-memory pipe and returns a client connection.
func testConn(t *testing.T, binder jsonrpc2.Binder) *jsonrpc2.Connection {
	t.Helper()
	ctx := context.Background()
	listener, err := jsonrpc2.NetPipe(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := jsonrpc2.Serve(ctx, listener, binder); err != nil {
		t.Fatal(err)
	}
	conn, err := jsonrpc2.Dial(ctx, listener.Dialer(), jsonrpc2.ConnectionOptions{})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		err = conn.Close()
		if err != nil {
			t.Fatal(err)
		}
		err = listener.Close()
		if err != nil {
			t.Fatal(err)
		}
	})
	return conn
}

// testTimeout bounds every wait in these tests, so that a regression fails
// fast instead of hanging until the go test timeout.
const testTimeout = 5 * time.Second

// testContext returns a context that expires after testTimeout.
func testContext(t *testing.T) context.Context {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	t.Cleanup(cancel)
	return ctx
}

// waitFor blocks until ch is closed and fails the test after testTimeout.
func waitFor(t *testing.T, ch <-chan struct{}, what string) {
	t.Helper()
	select {
	case <-ch:
	case <-time.After(testTimeout):
		t.Fatalf("timed out waiting for %s", what)
	}
}

// TestCancelRequest verifies that a $/cancelRequest notification cancels the
// context of an in-flight request, which then responds with -32800 even though
// the handler is otherwise blocking the connection's delivery loop.
func TestCancelRequest(t *testing.T) {
	started := make(chan struct{})
	slow := jsonrpc2.HandlerFunc(func(ctx context.Context, req *jsonrpc2.Request) (any, error) {
		close(started)
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(10 * time.Second):
			return nil, errors.New("handler was not cancelled")
		}
	})
	conn := testConn(t, lsp.CancelBinder(jsonrpc2.ConnectionOptions{Handler: slow}))

	ctx := testContext(t)
	call := conn.Call(ctx, "test/slow", nil)
	waitFor(t, started, "the handler to start")
	id := int32(call.ID().Raw().(int64))
	err := conn.Notify(ctx, "$/cancelRequest", &lsp.CancelParams{
		ID: lsp.CancelParamsIdFromInt32(id),
	})
	if err != nil {
		t.Fatal(err)
	}
	err = call.Await(ctx, nil)
	if err == nil || err.Error() != lsp.RequestCancelledError.Error() {
		t.Fatalf("expected RequestCancelledError, got %v", err)
	}
}

// TestCancelQueuedRequest verifies that a $/cancelRequest for a request that
// is still queued behind a running one is not lost: the queued request responds
// with -32800 and is never passed to the handler.
func TestCancelQueuedRequest(t *testing.T) {
	started := make(chan struct{})
	release := make(chan struct{})
	queuedHandled := make(chan struct{}, 1)
	inner := jsonrpc2.HandlerFunc(func(ctx context.Context, req *jsonrpc2.Request) (any, error) {
		switch req.Method {
		case "test/blocking":
			close(started)
			<-release
			return "done", nil
		case "test/queued":
			queuedHandled <- struct{}{}
			return "queued", nil
		}
		return nil, jsonrpc2.ErrNotHandled
	})
	// A marker notification sent right after the cancellation is preempted in
	// arrival order, so once it is seen the cancellation has been processed
	// and the blocking request can be released. This also checks that a
	// preempter supplied in the options keeps seeing other messages.
	cancelSeen := make(chan struct{})
	marker := preempterFunc(func(ctx context.Context, req *jsonrpc2.Request) (any, error) {
		if req.Method == "test/marker" {
			close(cancelSeen)
			return nil, nil
		}
		return nil, jsonrpc2.ErrNotHandled
	})
	conn := testConn(t, lsp.CancelBinder(jsonrpc2.ConnectionOptions{Preempter: marker, Handler: inner}))

	ctx := testContext(t)
	blocking := conn.Call(ctx, "test/blocking", nil)
	waitFor(t, started, "the blocking handler to start")
	queued := conn.Call(ctx, "test/queued", nil)
	queuedID := int32(queued.ID().Raw().(int64))
	err := conn.Notify(ctx, "$/cancelRequest", &lsp.CancelParams{
		ID: lsp.CancelParamsIdFromInt32(queuedID),
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := conn.Notify(ctx, "test/marker", nil); err != nil {
		t.Fatal(err)
	}
	waitFor(t, cancelSeen, "the cancellation to be processed")
	close(release)

	var result string
	if err := blocking.Await(ctx, &result); err != nil {
		t.Fatal(err)
	}
	if result != "done" {
		t.Errorf("expected result 'done', got %q", result)
	}
	err = queued.Await(ctx, nil)
	if err == nil || err.Error() != lsp.RequestCancelledError.Error() {
		t.Fatalf("expected RequestCancelledError, got %v", err)
	}
	select {
	case <-queuedHandled:
		t.Error("cancelled queued request was passed to the handler")
	default:
	}
}

type preempterFunc func(ctx context.Context, req *jsonrpc2.Request) (any, error)

func (f preempterFunc) Preempt(ctx context.Context, req *jsonrpc2.Request) (any, error) {
	return f(ctx, req)
}

// TestAsyncCancelRequest verifies the combined stack (CancelBinder around
// AsyncBinder): a $/cancelRequest cancels the targeted concurrent request with
// -32800 while another in-flight request is unaffected and completes normally.
func TestAsyncCancelRequest(t *testing.T) {
	started := make(chan struct{})
	release := make(chan struct{})
	inner := jsonrpc2.HandlerFunc(func(ctx context.Context, req *jsonrpc2.Request) (any, error) {
		switch req.Method {
		case "test/slow":
			close(started)
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(10 * time.Second):
				return nil, errors.New("handler was not cancelled")
			}
		case "test/fast":
			<-release
			return "fast", nil
		}
		return nil, jsonrpc2.ErrNotHandled
	})
	conn := testConn(t, lsp.CancelBinder(lsp.AsyncBinder(jsonrpc2.ConnectionOptions{Handler: inner})))

	ctx := testContext(t)
	slow := conn.Call(ctx, "test/slow", nil)
	waitFor(t, started, "the slow handler to start")
	fast := conn.Call(ctx, "test/fast", nil)

	id := int32(slow.ID().Raw().(int64))
	err := conn.Notify(ctx, "$/cancelRequest", &lsp.CancelParams{
		ID: lsp.CancelParamsIdFromInt32(id),
	})
	if err != nil {
		t.Fatal(err)
	}
	// The cancelled response must arrive while test/fast is still blocked,
	// proving the requests ran concurrently and only the target was cancelled.
	err = slow.Await(ctx, nil)
	if err == nil || err.Error() != lsp.RequestCancelledError.Error() {
		t.Fatalf("expected RequestCancelledError, got %v", err)
	}
	close(release)
	var result string
	if err := fast.Await(ctx, &result); err != nil {
		t.Fatal(err)
	}
	if result != "fast" {
		t.Errorf("expected result 'fast', got %q", result)
	}
}

// TestAsyncCloseCancelsRequests verifies that with AsyncBinder, a request that
// is still running when the connection shuts down has its context cancelled
// instead of being leaked.
func TestAsyncCloseCancelsRequests(t *testing.T) {
	started := make(chan struct{})
	cancelled := make(chan struct{})
	handler := jsonrpc2.HandlerFunc(func(ctx context.Context, req *jsonrpc2.Request) (any, error) {
		close(started)
		select {
		case <-ctx.Done():
			close(cancelled)
			return nil, ctx.Err()
		case <-time.After(2 * testTimeout):
			return nil, errors.New("handler was not cancelled")
		}
	})
	conn := testConn(t, lsp.AsyncBinder(jsonrpc2.ConnectionOptions{Handler: handler}))

	conn.Call(testContext(t), "test/slow", nil)
	waitFor(t, started, "the handler to start")
	if err := conn.Close(); err != nil {
		t.Fatal(err)
	}
	waitFor(t, cancelled, "the handler context to be cancelled")
}

// TestAsyncNotificationOrdering verifies that with AsyncBinder, requests run
// concurrently with each other, while a notification is only handled after
// every earlier request has responded and before any later request starts.
func TestAsyncNotificationOrdering(t *testing.T) {
	var mu sync.Mutex
	var log []string
	record := func(s string) {
		mu.Lock()
		defer mu.Unlock()
		log = append(log, s)
	}
	firstRunning := make(chan struct{})
	secondRunning := make(chan struct{})
	release := make(chan struct{})
	handler := jsonrpc2.HandlerFunc(func(ctx context.Context, req *jsonrpc2.Request) (any, error) {
		switch req.Method {
		case "test/first", "test/second":
			if req.Method == "test/first" {
				close(firstRunning)
			} else {
				close(secondRunning)
			}
			<-release
			record("request")
			return "done", nil
		case "test/after":
			record("after")
			return "done", nil
		case "test/note":
			var s string
			if err := json.Unmarshal(req.Params, &s); err != nil {
				return nil, err
			}
			record(s)
			return nil, nil
		}
		return nil, jsonrpc2.ErrNotHandled
	})
	conn := testConn(t, lsp.AsyncBinder(jsonrpc2.ConnectionOptions{Handler: handler}))

	ctx := testContext(t)
	first := conn.Call(ctx, "test/first", nil)
	second := conn.Call(ctx, "test/second", nil)
	// Both requests must be running at the same time.
	waitFor(t, firstRunning, "the first request to start")
	waitFor(t, secondRunning, "the second request to start")

	for _, s := range []string{"a", "b"} {
		if err := conn.Notify(ctx, "test/note", s); err != nil {
			t.Fatal(err)
		}
	}
	after := conn.Call(ctx, "test/after", nil)
	// Give the notifications a chance to be handled too early.
	time.Sleep(50 * time.Millisecond)
	close(release)

	for _, call := range []*jsonrpc2.AsyncCall{first, second, after} {
		var result string
		if err := call.Await(ctx, &result); err != nil {
			t.Fatal(err)
		}
		if result != "done" {
			t.Errorf("expected result 'done', got %q", result)
		}
	}
	mu.Lock()
	defer mu.Unlock()
	want := []string{"request", "request", "a", "b", "after"}
	if diff := cmp.Diff(want, log); diff != "" {
		t.Errorf("handling order mismatch (-want +got):\n%s", diff)
	}
}

// TestAsyncCloseWithQueuedNotification verifies that the connection can shut
// down while a notification is waiting for a request that only returns once
// its context is cancelled.
func TestAsyncCloseWithQueuedNotification(t *testing.T) {
	started := make(chan struct{})
	noted := make(chan struct{})
	handler := jsonrpc2.HandlerFunc(func(ctx context.Context, req *jsonrpc2.Request) (any, error) {
		switch req.Method {
		case "test/slow":
			close(started)
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(2 * testTimeout):
				return nil, errors.New("handler was not cancelled")
			}
		case "test/note":
			close(noted)
			return nil, nil
		}
		return nil, jsonrpc2.ErrNotHandled
	})
	conn := testConn(t, lsp.AsyncBinder(jsonrpc2.ConnectionOptions{Handler: handler}))

	ctx := testContext(t)
	conn.Call(ctx, "test/slow", nil)
	waitFor(t, started, "the handler to start")
	if err := conn.Notify(ctx, "test/note", nil); err != nil {
		t.Fatal(err)
	}
	if err := conn.Close(); err != nil {
		t.Fatal(err)
	}
	waitFor(t, noted, "the queued notification to be handled")
}

// TestCancellationSupport tests the LSP cancellation protocol
func TestCancellationSupport(t *testing.T) {

	t.Run("CancelParams", func(t *testing.T) {
		// Test CancelParams structure
		cancelParams := lsp.CancelParams{
			ID: lsp.CancelParamsIdFromString("test-request-123"),
		}

		data, err := json.Marshal(cancelParams)
		if err != nil {
			t.Fatalf("Failed to marshal CancelParams: %v", err)
		}

		var unmarshaled lsp.CancelParams
		err = json.Unmarshal(data, &unmarshaled)
		if err != nil {
			t.Fatalf("Failed to unmarshal CancelParams: %v", err)
		}

		if diff := cmp.Diff(cancelParams, unmarshaled); diff != "" {
			t.Errorf("CancelParams mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("CancelNotificationStructure", func(t *testing.T) {
		// Test that cancel notification structure is correct
		cancelParams := lsp.CancelParams{
			ID: lsp.CancelParamsIdFromString("request-to-cancel"),
		}

		// Test JSON marshaling
		data, err := json.Marshal(cancelParams)
		if err != nil {
			t.Fatalf("Failed to marshal CancelParams: %v", err)
		}

		var unmarshaled lsp.CancelParams
		err = json.Unmarshal(data, &unmarshaled)
		if err != nil {
			t.Fatalf("Failed to unmarshal CancelParams: %v", err)
		}

		if diff := cmp.Diff(cancelParams, unmarshaled); diff != "" {
			t.Errorf("CancelParams mismatch (-want +got):\n%s", diff)
		}
	})
}

// TestProgressSupport tests the LSP progress protocol
func TestProgressSupport(t *testing.T) {

	t.Run("ProgressParams", func(t *testing.T) {
		// Test ProgressParams structure
		progressParams := lsp.ProgressParams{
			Token: "progress-token-123",
			Value: map[string]any{
				"kind":    "begin",
				"title":   "Processing",
				"message": "Starting work",
			},
		}

		data, err := json.Marshal(progressParams)
		if err != nil {
			t.Fatalf("Failed to marshal ProgressParams: %v", err)
		}

		var unmarshaled lsp.ProgressParams
		err = json.Unmarshal(data, &unmarshaled)
		if err != nil {
			t.Fatalf("Failed to unmarshal ProgressParams: %v", err)
		}

		if diff := cmp.Diff(progressParams, unmarshaled); diff != "" {
			t.Errorf("ProgressParams mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("WorkDoneProgressBegin", func(t *testing.T) {
		// Test WorkDoneProgressBegin structure
		begin := lsp.WorkDoneProgressBegin{
			Kind:        "begin",
			Title:       "Indexing",
			Cancellable: true,
			Message:     "Processing files",
			Percentage:  uintPtr(0),
		}

		data, err := json.Marshal(begin)
		if err != nil {
			t.Fatalf("Failed to marshal WorkDoneProgressBegin: %v", err)
		}

		var unmarshaled lsp.WorkDoneProgressBegin
		err = json.Unmarshal(data, &unmarshaled)
		if err != nil {
			t.Fatalf("Failed to unmarshal WorkDoneProgressBegin: %v", err)
		}

		if diff := cmp.Diff(begin, unmarshaled); diff != "" {
			t.Errorf("WorkDoneProgressBegin mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("WorkDoneProgressReport", func(t *testing.T) {
		// Test WorkDoneProgressReport structure
		report := lsp.WorkDoneProgressReport{
			Kind:        "report",
			Cancellable: true,
			Message:     "50% complete",
			Percentage:  uintPtr(50),
		}

		data, err := json.Marshal(report)
		if err != nil {
			t.Fatalf("Failed to marshal WorkDoneProgressReport: %v", err)
		}

		var unmarshaled lsp.WorkDoneProgressReport
		err = json.Unmarshal(data, &unmarshaled)
		if err != nil {
			t.Fatalf("Failed to unmarshal WorkDoneProgressReport: %v", err)
		}

		if diff := cmp.Diff(report, unmarshaled); diff != "" {
			t.Errorf("WorkDoneProgressReport mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("WorkDoneProgressEnd", func(t *testing.T) {
		// Test WorkDoneProgressEnd structure
		end := lsp.WorkDoneProgressEnd{
			Kind:    "end",
			Message: "Indexing complete",
		}

		data, err := json.Marshal(end)
		if err != nil {
			t.Fatalf("Failed to marshal WorkDoneProgressEnd: %v", err)
		}

		var unmarshaled lsp.WorkDoneProgressEnd
		err = json.Unmarshal(data, &unmarshaled)
		if err != nil {
			t.Fatalf("Failed to unmarshal WorkDoneProgressEnd: %v", err)
		}

		if end.Message != unmarshaled.Message {
			t.Errorf("Expected Message %v, got %v", end.Message, unmarshaled.Message)
		}
	})

	t.Run("ProgressNotificationStructure", func(t *testing.T) {
		// Test progress notification structures
		beginParams := lsp.ProgressParams{
			Token: "test-progress",
			Value: lsp.WorkDoneProgressBegin{
				Kind:        "begin",
				Title:       "Test Operation",
				Cancellable: false,
				Percentage:  uintPtr(0),
			},
		}

		// Test JSON marshaling of begin params
		data, err := json.Marshal(beginParams)
		if err != nil {
			t.Fatalf("Failed to marshal progress begin params: %v", err)
		}

		var unmarshaled lsp.ProgressParams
		err = json.Unmarshal(data, &unmarshaled)
		if err != nil {
			t.Fatalf("Failed to unmarshal progress begin params: %v", err)
		}

		if beginParams.Token != unmarshaled.Token {
			t.Errorf("Expected Token %v, got %v", beginParams.Token, unmarshaled.Token)
		}
	})
}

// TestErrorHandling tests LSP error handling scenarios
func TestErrorHandling(t *testing.T) {
	t.Run("RequestCancelledError", func(t *testing.T) {
		// Test that RequestCancelledError is available
		err := lsp.RequestCancelledError
		if err == nil {
			t.Error("RequestCancelledError should not be nil")
		}

		// Test that it has the expected message
		expectedMsg := "JSON RPC cancelled"
		if err.Error() != expectedMsg {
			t.Errorf("Expected error message %q, got %q", expectedMsg, err.Error())
		}
	})

	t.Run("ErrorCodes", func(t *testing.T) {
		// Test that RequestCancelledError exists and can be used
		err := lsp.RequestCancelledError
		if err == nil {
			t.Error("RequestCancelledError should be defined")
		}

		// Test that it's a proper error
		if err.Error() == "" {
			t.Error("RequestCancelledError should have a non-empty error message")
		}
	})
}

// Helper function for creating uint32 pointers
func uintPtr(u uint32) *uint32 {
	return &u
}
