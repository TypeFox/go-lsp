// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package protocol_test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/TypeFox/go-lsp/internal/event/export/eventtest"
	"github.com/TypeFox/go-lsp/internal/jsonrpc2"
	"github.com/TypeFox/go-lsp/internal/jsonrpc2/stack/stacktest"
	"github.com/TypeFox/go-lsp/protocol"
)

// TestClientServerIntegration tests actual client-server communication
func TestClientServerIntegration(t *testing.T) {
	stacktest.NoLeak(t)
	ctx := eventtest.NewContext(context.Background(), t)

	t.Run("InitializeSequence", func(t *testing.T) {
		client, _, cleanup := setupLSPConnection(ctx, t)
		defer cleanup()

		// Test initialize request
		initParams := protocol.ParamInitialize{
			XInitializeParams: protocol.XInitializeParams{
				ProcessID: 12345,
				ClientInfo: &protocol.ClientInfo{
					Name:    "test-client",
					Version: "1.0.0",
				},
				RootURI: "file:///workspace",
				Capabilities: protocol.ClientCapabilities{
					TextDocument: protocol.TextDocumentClientCapabilities{
						Hover: &protocol.HoverClientCapabilities{
							DynamicRegistration: true,
							ContentFormat:       []protocol.MarkupKind{protocol.PlainText, protocol.Markdown},
						},
					},
				},
			},
		}

		var initResult protocol.InitializeResult
		_, err := client.Call(ctx, "initialize", initParams, &initResult)
		if err != nil {
			t.Fatalf("Initialize request failed: %v", err)
		}

		// Send initialized notification
		err = client.Notify(ctx, "initialized", protocol.InitializedParams{})
		if err != nil {
			t.Fatalf("Initialized notification failed: %v", err)
		}

		// Give some time for processing
		time.Sleep(10 * time.Millisecond)
	})

	t.Run("TextDocumentLifecycle", func(t *testing.T) {
		client, _, cleanup := setupLSPConnection(ctx, t)
		defer cleanup()

		// Initialize first
		initParams := protocol.ParamInitialize{
			XInitializeParams: protocol.XInitializeParams{
				ProcessID: 12345,
				Capabilities: protocol.ClientCapabilities{
					TextDocument: protocol.TextDocumentClientCapabilities{
						Synchronization: &protocol.TextDocumentSyncClientCapabilities{
							DynamicRegistration: true,
						},
					},
				},
			},
		}

		var initResult protocol.InitializeResult
		_, err := client.Call(ctx, "initialize", initParams, &initResult)
		if err != nil {
			t.Fatalf("Initialize request failed: %v", err)
		}

		// Send didOpen notification
		didOpenParams := protocol.DidOpenTextDocumentParams{
			TextDocument: protocol.TextDocumentItem{
				URI:        "file:///test.go",
				LanguageID: "go",
				Version:    1,
				Text:       "package main\n\nfunc main() {}\n",
			},
		}

		err = client.Notify(ctx, "textDocument/didOpen", didOpenParams)
		if err != nil {
			t.Fatalf("didOpen notification failed: %v", err)
		}

		// Send didChange notification
		didChangeParams := protocol.DidChangeTextDocumentParams{
			TextDocument: protocol.VersionedTextDocumentIdentifier{
				TextDocumentIdentifier: protocol.TextDocumentIdentifier{
					URI: "file:///test.go",
				},
				Version: 2,
			},
			ContentChanges: []protocol.TextDocumentContentChangeEvent{
				{
					Text: "package main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"Hello, World!\")\n}\n",
				},
			},
		}

		err = client.Notify(ctx, "textDocument/didChange", didChangeParams)
		if err != nil {
			t.Fatalf("didChange notification failed: %v", err)
		}

		// Send didSave notification
		didSaveParams := protocol.DidSaveTextDocumentParams{
			TextDocument: protocol.TextDocumentIdentifier{
				URI: "file:///test.go",
			},
		}

		err = client.Notify(ctx, "textDocument/didSave", didSaveParams)
		if err != nil {
			t.Fatalf("didSave notification failed: %v", err)
		}

		// Send didClose notification
		didCloseParams := protocol.DidCloseTextDocumentParams{
			TextDocument: protocol.TextDocumentIdentifier{
				URI: "file:///test.go",
			},
		}

		err = client.Notify(ctx, "textDocument/didClose", didCloseParams)
		if err != nil {
			t.Fatalf("didClose notification failed: %v", err)
		}

		// Give some time for processing
		time.Sleep(10 * time.Millisecond)
	})

	t.Run("LanguageFeatures", func(t *testing.T) {
		client, _, cleanup := setupLSPConnection(ctx, t)
		defer cleanup()

		// Initialize with language feature capabilities
		initParams := protocol.ParamInitialize{
			XInitializeParams: protocol.XInitializeParams{
				ProcessID: 12345,
				Capabilities: protocol.ClientCapabilities{
					TextDocument: protocol.TextDocumentClientCapabilities{
						Hover: &protocol.HoverClientCapabilities{
							DynamicRegistration: true,
							ContentFormat:       []protocol.MarkupKind{protocol.PlainText},
						},
						Completion: protocol.CompletionClientCapabilities{
							DynamicRegistration: true,
						},
					},
				},
			},
		}

		var initResult protocol.InitializeResult
		_, err := client.Call(ctx, "initialize", initParams, &initResult)
		if err != nil {
			t.Fatalf("Initialize request failed: %v", err)
		}

		// Test hover request
		hoverParams := protocol.HoverParams{
			TextDocumentPositionParams: protocol.TextDocumentPositionParams{
				TextDocument: protocol.TextDocumentIdentifier{
					URI: "file:///test.go",
				},
				Position: protocol.Position{
					Line:      2,
					Character: 5,
				},
			},
		}

		var hoverResult protocol.Hover
		_, err = client.Call(ctx, "textDocument/hover", hoverParams, &hoverResult)
		if err != nil {
			t.Fatalf("Hover request failed: %v", err)
		}

		// Test completion request
		completionParams := protocol.CompletionParams{
			TextDocumentPositionParams: protocol.TextDocumentPositionParams{
				TextDocument: protocol.TextDocumentIdentifier{
					URI: "file:///test.go",
				},
				Position: protocol.Position{
					Line:      3,
					Character: 1,
				},
			},
		}

		var completionResult protocol.CompletionList
		_, err = client.Call(ctx, "textDocument/completion", completionParams, &completionResult)
		if err != nil {
			t.Fatalf("Completion request failed: %v", err)
		}

		// Give some time for processing
		time.Sleep(10 * time.Millisecond)
	})
}

// setupLSPConnection creates a mock LSP client-server connection for testing
func setupLSPConnection(ctx context.Context, t *testing.T) (jsonrpc2.Conn, jsonrpc2.Conn, func()) {
	clientPipe, serverPipe := net.Pipe()

	// Mock server handler that responds to LSP requests
	serverHandler := func(ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
		switch req.Method() {
		case "initialize":
			result := protocol.InitializeResult{
				Capabilities: protocol.ServerCapabilities{
					TextDocumentSync: &protocol.TextDocumentSyncOptions{
						OpenClose: true,
						Change:    protocol.Incremental,
					},
					CompletionProvider: &protocol.CompletionOptions{},
				},
			}
			return reply(ctx, result, nil)

		case "textDocument/hover":
			result := protocol.Hover{
				Contents: protocol.MarkupContent{
					Kind:  protocol.PlainText,
					Value: "func main()",
				},
			}
			return reply(ctx, result, nil)

		case "textDocument/completion":
			result := protocol.CompletionList{
				IsIncomplete: false,
				Items: []protocol.CompletionItem{
					{
						Label:  "println",
						Kind:   protocol.FunctionCompletion,
						Detail: "func println(a ...interface{})",
					},
				},
			}
			return reply(ctx, result, nil)

		default:
			// For notifications and other requests, just acknowledge
			return reply(ctx, nil, nil)
		}
	}

	// Mock client handler (for server-to-client requests)
	clientHandler := func(ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
		return reply(ctx, nil, nil)
	}

	// Create connections using the existing jsonrpc2 test pattern
	serverStream := jsonrpc2.NewRawStream(serverPipe)
	clientStream := jsonrpc2.NewRawStream(clientPipe)

	server := jsonrpc2.NewConn(serverStream)
	client := jsonrpc2.NewConn(clientStream)

	// Start the handlers
	server.Go(ctx, serverHandler)
	client.Go(ctx, clientHandler)

	cleanup := func() {
		server.Close()
		client.Close()
		<-server.Done()
		<-client.Done()
	}

	return client, server, cleanup
}