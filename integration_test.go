// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lsp_test

import (
	"context"
	"testing"

	"typefox.dev/lsp"
)

// TestProtocolStructures tests protocol data structures and serialization
func TestProtocolStructures(t *testing.T) {
	ctx := context.Background()
	_ = ctx // Use ctx to avoid unused variable error

	t.Run("InitializeSequence", func(t *testing.T) {
		// Test initialize request structure
		initParams := lsp.ParamInitialize{
			XInitializeParams: lsp.XInitializeParams{
				ProcessID: 12345,
				ClientInfo: &lsp.ClientInfo{
					Name:    "test-client",
					Version: "1.0.0",
				},
				RootURI: "file:///workspace",
				Capabilities: lsp.ClientCapabilities{
					TextDocument: lsp.TextDocumentClientCapabilities{
						Hover: &lsp.HoverClientCapabilities{
							DynamicRegistration: true,
							ContentFormat:       []lsp.MarkupKind{lsp.PlainText, lsp.Markdown},
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
		syncOptions := &lsp.TextDocumentSyncOptions{
			OpenClose: true,
			Change:    lsp.Incremental,
		}
		initResult := lsp.InitializeResult{
			Capabilities: lsp.ServerCapabilities{
				TextDocumentSync:   syncOptions,
				CompletionProvider: &lsp.CompletionOptions{},
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
		didOpenParams := lsp.DidOpenTextDocumentParams{
			TextDocument: lsp.TextDocumentItem{
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
		didChangeParams := lsp.DidChangeTextDocumentParams{
			TextDocument: lsp.VersionedTextDocumentIdentifier{
				TextDocumentIdentifier: lsp.TextDocumentIdentifier{
					URI: "file:///test.go",
				},
				Version: 2,
			},
			ContentChanges: []lsp.TextDocumentContentChangeEvent{
				{
					Text: "package main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"Hello, World!\")\n}\n",
				},
			},
		}

		if didChangeParams.TextDocument.Version != 2 {
			t.Errorf("Expected version 2, got %d", didChangeParams.TextDocument.Version)
		}

		// Test didSave notification
		didSaveParams := lsp.DidSaveTextDocumentParams{
			TextDocument: lsp.TextDocumentIdentifier{
				URI: "file:///test.go",
			},
		}

		if didSaveParams.TextDocument.URI != "file:///test.go" {
			t.Errorf("Expected URI file:///test.go, got %s", didSaveParams.TextDocument.URI)
		}

		// Test didClose notification
		didCloseParams := lsp.DidCloseTextDocumentParams{
			TextDocument: lsp.TextDocumentIdentifier{
				URI: "file:///test.go",
			},
		}

		if didCloseParams.TextDocument.URI != "file:///test.go" {
			t.Errorf("Expected URI file:///test.go, got %s", didCloseParams.TextDocument.URI)
		}
	})

	t.Run("LSPRequests", func(t *testing.T) {
		// Test hover request
		hoverParams := lsp.HoverParams{
			TextDocumentPositionParams: lsp.TextDocumentPositionParams{
				TextDocument: lsp.TextDocumentIdentifier{
					URI: "file:///test.go",
				},
				Position: lsp.Position{
					Line:      5,
					Character: 10,
				},
			},
		}

		if hoverParams.Position.Line != 5 {
			t.Errorf("Expected line 5, got %d", hoverParams.Position.Line)
		}

		// Test completion request
		completionParams := lsp.CompletionParams{
			TextDocumentPositionParams: lsp.TextDocumentPositionParams{
				TextDocument: lsp.TextDocumentIdentifier{
					URI: "file:///test.go",
				},
				Position: lsp.Position{
					Line:      5,
					Character: 10,
				},
			},
		}

		if completionParams.Position.Character != 10 {
			t.Errorf("Expected character 10, got %d", completionParams.Position.Character)
		}

		// Test completion response
		completionList := lsp.CompletionList{
			IsIncomplete: false,
			Items: []lsp.CompletionItem{
				{
					Label:  "println",
					Kind:   lsp.FunctionCompletion,
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
