package lsp_test

import (
	"context"
	"encoding/json"
	"sync/atomic"
	"testing"
	"time"

	"golang.org/x/exp/jsonrpc2"
	"typefox.dev/lsp"
)

// lifecycleServer is a minimal lsp.Server that records which of its methods
// were invoked. Methods that are not overridden panic through the nil
// embedded interface, which is fine since the tests never reach them.
type lifecycleServer struct {
	lsp.Server

	shutdowns atomic.Int32
	hovers    atomic.Int32
	didOpens  atomic.Int32
	exited    chan struct{}
}

func (s *lifecycleServer) Shutdown(context.Context) error {
	s.shutdowns.Add(1)
	return nil
}

func (s *lifecycleServer) Exit(context.Context) error {
	close(s.exited)
	return nil
}

func (s *lifecycleServer) Hover(context.Context, *lsp.HoverParams) (*lsp.Hover, error) {
	s.hovers.Add(1)
	return &lsp.Hover{}, nil
}

func (s *lifecycleServer) DidOpen(context.Context, *lsp.DidOpenTextDocumentParams) error {
	s.didOpens.Add(1)
	return nil
}

// connectLifecycleServer serves impl over an in-memory pipe and returns a
// dispatcher that talks to it like a language client would.
func connectLifecycleServer(t *testing.T, impl lsp.Server) lsp.Server {
	t.Helper()
	ctx := context.Background()

	listener, err := jsonrpc2.NetPipe(ctx)
	if err != nil {
		t.Fatal(err)
	}
	rpcServer, err := jsonrpc2.Serve(ctx, listener, jsonrpc2.ConnectionOptions{
		Handler: lsp.ServerHandler(impl),
	})
	if err != nil {
		t.Fatal(err)
	}
	conn, err := jsonrpc2.Dial(ctx, listener.Dialer(), jsonrpc2.ConnectionOptions{})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = conn.Close()
		_ = listener.Close()
		_ = rpcServer.Wait()
	})
	return lsp.ServerDispatcher(conn)
}

// errorCode extracts the JSON-RPC error code from an error returned by a call.
func errorCode(t *testing.T, err error) lsp.ErrorCodes {
	t.Helper()
	data, marshalErr := json.Marshal(err)
	if marshalErr != nil {
		t.Fatalf("error %v is not a JSON-RPC error: %v", err, marshalErr)
	}
	var wire struct {
		Code *lsp.ErrorCodes `json:"code"`
	}
	if unmarshalErr := json.Unmarshal(data, &wire); unmarshalErr != nil || wire.Code == nil {
		t.Fatalf("error %v (%s) is not a JSON-RPC error", err, data)
	}
	return *wire.Code
}

// TestShutdown tests the shutdown/exit sequence of the LSP lifecycle.
// See https://microsoft.github.io/language-server-protocol/specifications/lsp/3.18/specification#shutdown
func TestShutdown(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	t.Run("ShutdownSucceeds", func(t *testing.T) {
		impl := &lifecycleServer{exited: make(chan struct{})}
		client := connectLifecycleServer(t, impl)

		if _, err := client.Hover(ctx, &lsp.HoverParams{}); err != nil {
			t.Fatalf("Hover before shutdown failed: %v", err)
		}
		if err := client.Shutdown(ctx); err != nil {
			t.Fatalf("Shutdown failed: %v", err)
		}
		if got := impl.shutdowns.Load(); got != 1 {
			t.Errorf("Expected Shutdown to be handled once, got %d", got)
		}
	})

	t.Run("RequestsAfterShutdownAreInvalid", func(t *testing.T) {
		impl := &lifecycleServer{exited: make(chan struct{})}
		client := connectLifecycleServer(t, impl)

		if err := client.Shutdown(ctx); err != nil {
			t.Fatalf("Shutdown failed: %v", err)
		}

		_, err := client.Hover(ctx, &lsp.HoverParams{})
		if err == nil {
			t.Fatal("Expected Hover after shutdown to fail")
		}
		if code := errorCode(t, err); code != lsp.InvalidRequest {
			t.Errorf("Expected error code %d (InvalidRequest), got %d: %v", lsp.InvalidRequest, code, err)
		}
		if got := impl.hovers.Load(); got != 0 {
			t.Errorf("Expected Hover not to reach the server after shutdown, got %d calls", got)
		}

		// A second shutdown is a request like any other.
		err = client.Shutdown(ctx)
		if err == nil {
			t.Fatal("Expected second Shutdown to fail")
		}
		if code := errorCode(t, err); code != lsp.InvalidRequest {
			t.Errorf("Expected error code %d (InvalidRequest), got %d: %v", lsp.InvalidRequest, code, err)
		}
		if got := impl.shutdowns.Load(); got != 1 {
			t.Errorf("Expected Shutdown to be handled once, got %d", got)
		}
	})

	t.Run("NotificationsAfterShutdownAreDropped", func(t *testing.T) {
		impl := &lifecycleServer{exited: make(chan struct{})}
		client := connectLifecycleServer(t, impl)

		if err := client.Shutdown(ctx); err != nil {
			t.Fatalf("Shutdown failed: %v", err)
		}
		if err := client.DidOpen(ctx, &lsp.DidOpenTextDocumentParams{}); err != nil {
			t.Fatalf("DidOpen failed: %v", err)
		}
		// Messages are handled in order, so once exit has arrived the
		// preceding didOpen has been processed as well.
		if err := client.Exit(ctx); err != nil {
			t.Fatalf("Exit failed: %v", err)
		}
		select {
		case <-impl.exited:
		case <-ctx.Done():
			t.Fatal("Expected exit notification to reach the server after shutdown")
		}
		if got := impl.didOpens.Load(); got != 0 {
			t.Errorf("Expected DidOpen to be dropped after shutdown, got %d calls", got)
		}
	})

	t.Run("ExitWithoutShutdown", func(t *testing.T) {
		impl := &lifecycleServer{exited: make(chan struct{})}
		client := connectLifecycleServer(t, impl)

		if err := client.Exit(ctx); err != nil {
			t.Fatalf("Exit failed: %v", err)
		}
		select {
		case <-impl.exited:
		case <-ctx.Done():
			t.Fatal("Expected exit notification to reach the server")
		}
	})
}
