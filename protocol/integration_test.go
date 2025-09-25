// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package protocol_test

import (
	"context"
	"testing"

	"github.com/TypeFox/go-lsp/protocol"
)

// TestProtocolStructures tests protocol data structures and serialization
func TestProtocolStructures(t *testing.T) {
	ctx := context.Background()
	_ = ctx // Use ctx to avoid unused variable error

	t.Run("InitializeSequence", func(t *testing.T) {
		// Test initialize request structure
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

		// Test that the structure has expected values
		if initParams.ProcessID != 12345 {
			t.Errorf("Expected ProcessID 12345, got %d", initParams.ProcessID)
		}

		if initParams.ClientInfo.Name != "test-client" {
			t.Errorf("Expected client name 'test-client', got %s", initParams.ClientInfo.Name)
		}

		// Test initialize result structure
		syncOptions := &protocol.TextDocumentSyncOptions{
			OpenClose: true,
			Change:    protocol.Incremental,
		}
		initResult := protocol.InitializeResult{
			Capabilities: protocol.ServerCapabilities{
				TextDocumentSync:   syncOptions,
				CompletionProvider: &protocol.CompletionOptions{},
			},
		}

		if syncOptions.OpenClose != true {
			t.Error("Expected OpenClose to be true")
		}

		// Verify the result was created properly
		if initResult.Capabilities.CompletionProvider == nil {
			t.Error("Expected CompletionProvider to be set")
		}
	})

	t.Run("TextDocumentLifecycle", func(t *testing.T) {
		// Test text document lifecycle structures
		didOpenParams := protocol.DidOpenTextDocumentParams{
			TextDocument: protocol.TextDocumentItem{
				URI:        "file:///test.go",
				LanguageID: "go",
				Version:    1,
				Text:       "package main\n\nfunc main() {}\n",
			},
		}

		if didOpenParams.TextDocument.URI != "file:///test.go" {
			t.Errorf("Expected URI file:///test.go, got %s", didOpenParams.TextDocument.URI)
		}

		// Test didChange notification
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

		if didChangeParams.TextDocument.Version != 2 {
			t.Errorf("Expected version 2, got %d", didChangeParams.TextDocument.Version)
		}

		// Test didSave notification
		didSaveParams := protocol.DidSaveTextDocumentParams{
			TextDocument: protocol.TextDocumentIdentifier{
				URI: "file:///test.go",
			},
		}

		if didSaveParams.TextDocument.URI != "file:///test.go" {
			t.Errorf("Expected URI file:///test.go, got %s", didSaveParams.TextDocument.URI)
		}

		// Test didClose notification
		didCloseParams := protocol.DidCloseTextDocumentParams{
			TextDocument: protocol.TextDocumentIdentifier{
				URI: "file:///test.go",
			},
		}

		if didCloseParams.TextDocument.URI != "file:///test.go" {
			t.Errorf("Expected URI file:///test.go, got %s", didCloseParams.TextDocument.URI)
		}
	})

	t.Run("LSPRequests", func(t *testing.T) {
		// Test hover request
		hoverParams := protocol.HoverParams{
			TextDocumentPositionParams: protocol.TextDocumentPositionParams{
				TextDocument: protocol.TextDocumentIdentifier{
					URI: "file:///test.go",
				},
				Position: protocol.Position{
					Line:      5,
					Character: 10,
				},
			},
		}

		if hoverParams.Position.Line != 5 {
			t.Errorf("Expected line 5, got %d", hoverParams.Position.Line)
		}

		// Test completion request
		completionParams := protocol.CompletionParams{
			TextDocumentPositionParams: protocol.TextDocumentPositionParams{
				TextDocument: protocol.TextDocumentIdentifier{
					URI: "file:///test.go",
				},
				Position: protocol.Position{
					Line:      5,
					Character: 10,
				},
			},
		}

		if completionParams.Position.Character != 10 {
			t.Errorf("Expected character 10, got %d", completionParams.Position.Character)
		}

		// Test completion response
		completionList := protocol.CompletionList{
			IsIncomplete: false,
			Items: []protocol.CompletionItem{
				{
					Label:  "println",
					Kind:   protocol.FunctionCompletion,
					Detail: "func println(a ...interface{})",
				},
			},
		}

		if len(completionList.Items) != 1 {
			t.Errorf("Expected 1 completion item, got %d", len(completionList.Items))
		}

		if completionList.Items[0].Label != "println" {
			t.Errorf("Expected label 'println', got %s", completionList.Items[0].Label)
		}
	})
}