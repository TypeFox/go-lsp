// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lsp_test

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
	"typefox.dev/lsp"
)

// TestWorkspaceEdit tests complex workspace edit scenarios
func TestWorkspaceEdit(t *testing.T) {
	t.Run("BasicWorkspaceEdit", func(t *testing.T) {
		edit := lsp.WorkspaceEdit{
			Changes: map[lsp.DocumentURI][]lsp.TextEdit{
				"file:///test.go": {
					{
						Range: lsp.Range{
							Start: lsp.Position{Line: 0, Character: 0},
							End:   lsp.Position{Line: 0, Character: 7},
						},
						NewText: "package main",
					},
				},
			},
		}

		data, err := json.Marshal(edit)
		if err != nil {
			t.Fatalf("Failed to marshal WorkspaceEdit: %v", err)
		}

		var unmarshaled lsp.WorkspaceEdit
		err = json.Unmarshal(data, &unmarshaled)
		if err != nil {
			t.Fatalf("Failed to unmarshal WorkspaceEdit: %v", err)
		}

		if diff := cmp.Diff(edit, unmarshaled); diff != "" {
			t.Errorf("WorkspaceEdit mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("TextDocumentEdit", func(t *testing.T) {
		// Test TextDocumentEdit structure
		edit := lsp.TextDocumentEdit{
			TextDocument: lsp.OptionalVersionedTextDocumentIdentifier{
				TextDocumentIdentifier: lsp.TextDocumentIdentifier{
					URI: "file:///test.go",
				},
				Version: 1,
			},
		}

		data, err := json.Marshal(edit)
		if err != nil {
			t.Fatalf("Failed to marshal TextDocumentEdit: %v", err)
		}

		var unmarshaled lsp.TextDocumentEdit
		err = json.Unmarshal(data, &unmarshaled)
		if err != nil {
			t.Fatalf("Failed to unmarshal TextDocumentEdit: %v", err)
		}

		if diff := cmp.Diff(edit, unmarshaled); diff != "" {
			t.Errorf("TextDocumentEdit mismatch (-want +got):\n%s", diff)
		}
	})
}

// TestDiagnostics tests diagnostic publishing scenarios
func TestDiagnostics(t *testing.T) {
	t.Run("PublishDiagnostics", func(t *testing.T) {
		params := lsp.PublishDiagnosticsParams{
			URI: "file:///test.go",
			Diagnostics: []lsp.Diagnostic{
				{
					Range: lsp.Range{
						Start: lsp.Position{Line: 5, Character: 10},
						End:   lsp.Position{Line: 5, Character: 15},
					},
					Severity: lsp.SeverityError,
					Code:     "unused",
					Source:   "go",
					Message:  "variable 'x' is unused",
					RelatedInformation: []lsp.DiagnosticRelatedInformation{
						{
							Location: lsp.Location{
								URI: "file:///test.go",
								Range: lsp.Range{
									Start: lsp.Position{Line: 3, Character: 5},
									End:   lsp.Position{Line: 3, Character: 6},
								},
							},
							Message: "variable 'x' declared here",
						},
					},
				},
				{
					Range: lsp.Range{
						Start: lsp.Position{Line: 10, Character: 0},
						End:   lsp.Position{Line: 10, Character: 20},
					},
					Severity: lsp.SeverityWarning,
					Code:     "deprecated",
					Source:   "go",
					Message:  "function is deprecated",
					Tags:     []lsp.DiagnosticTag{lsp.Deprecated},
				},
			},
		}

		data, err := json.Marshal(params)
		if err != nil {
			t.Fatalf("Failed to marshal PublishDiagnosticsParams: %v", err)
		}

		var unmarshaled lsp.PublishDiagnosticsParams
		err = json.Unmarshal(data, &unmarshaled)
		if err != nil {
			t.Fatalf("Failed to unmarshal PublishDiagnosticsParams: %v", err)
		}

		if diff := cmp.Diff(params, unmarshaled); diff != "" {
			t.Errorf("PublishDiagnosticsParams mismatch (-want +got):\n%s", diff)
		}
	})
}

// TestSemanticTokens tests semantic tokens scenarios
func TestSemanticTokens(t *testing.T) {
	t.Run("SemanticTokensParams", func(t *testing.T) {
		params := lsp.SemanticTokensParams{
			TextDocument: lsp.TextDocumentIdentifier{
				URI: "file:///test.go",
			},
		}

		data, err := json.Marshal(params)
		if err != nil {
			t.Fatalf("Failed to marshal SemanticTokensParams: %v", err)
		}

		var unmarshaled lsp.SemanticTokensParams
		err = json.Unmarshal(data, &unmarshaled)
		if err != nil {
			t.Fatalf("Failed to unmarshal SemanticTokensParams: %v", err)
		}

		if diff := cmp.Diff(params, unmarshaled); diff != "" {
			t.Errorf("SemanticTokensParams mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("SemanticTokens", func(t *testing.T) {
		tokens := lsp.SemanticTokens{
			ResultID: "version-1",
			Data:     []uint32{0, 0, 7, 0, 0, 1, 0, 6, 1, 0}, // Example token data
		}

		data, err := json.Marshal(tokens)
		if err != nil {
			t.Fatalf("Failed to marshal SemanticTokens: %v", err)
		}

		var unmarshaled lsp.SemanticTokens
		err = json.Unmarshal(data, &unmarshaled)
		if err != nil {
			t.Fatalf("Failed to unmarshal SemanticTokens: %v", err)
		}

		if diff := cmp.Diff(tokens, unmarshaled); diff != "" {
			t.Errorf("SemanticTokens mismatch (-want +got):\n%s", diff)
		}
	})
}

// TestCodeAction tests code action scenarios
func TestCodeAction(t *testing.T) {
	t.Run("CodeActionParams", func(t *testing.T) {
		params := lsp.CodeActionParams{
			TextDocument: lsp.TextDocumentIdentifier{
				URI: "file:///test.go",
			},
			Range: lsp.Range{
				Start: lsp.Position{Line: 5, Character: 0},
				End:   lsp.Position{Line: 5, Character: 10},
			},
			Context: lsp.CodeActionContext{
				Diagnostics: []lsp.Diagnostic{
					{
						Range: lsp.Range{
							Start: lsp.Position{Line: 5, Character: 5},
							End:   lsp.Position{Line: 5, Character: 10},
						},
						Severity: lsp.SeverityError,
						Code:     "unused",
						Message:  "unused variable",
					},
				},
				Only: []lsp.CodeActionKind{lsp.QuickFix},
			},
		}

		data, err := json.Marshal(params)
		if err != nil {
			t.Fatalf("Failed to marshal CodeActionParams: %v", err)
		}

		var unmarshaled lsp.CodeActionParams
		err = json.Unmarshal(data, &unmarshaled)
		if err != nil {
			t.Fatalf("Failed to unmarshal CodeActionParams: %v", err)
		}

		if diff := cmp.Diff(params, unmarshaled); diff != "" {
			t.Errorf("CodeActionParams mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("CodeAction", func(t *testing.T) {
		action := lsp.CodeAction{
			Title: "Remove unused variable",
			Kind:  lsp.QuickFix,
			Edit: &lsp.WorkspaceEdit{
				Changes: map[lsp.DocumentURI][]lsp.TextEdit{
					"file:///test.go": {
						{
							Range: lsp.Range{
								Start: lsp.Position{Line: 5, Character: 0},
								End:   lsp.Position{Line: 6, Character: 0},
							},
							NewText: "",
						},
					},
				},
			},
			IsPreferred: true,
		}

		data, err := json.Marshal(action)
		if err != nil {
			t.Fatalf("Failed to marshal CodeAction: %v", err)
		}

		var unmarshaled lsp.CodeAction
		err = json.Unmarshal(data, &unmarshaled)
		if err != nil {
			t.Fatalf("Failed to unmarshal CodeAction: %v", err)
		}

		if diff := cmp.Diff(action, unmarshaled); diff != "" {
			t.Errorf("CodeAction mismatch (-want +got):\n%s", diff)
		}
	})
}

// TestCallHierarchy tests call hierarchy scenarios
func TestCallHierarchy(t *testing.T) {
	t.Run("CallHierarchyPrepareParams", func(t *testing.T) {
		params := lsp.CallHierarchyPrepareParams{
			TextDocumentPositionParams: lsp.TextDocumentPositionParams{
				TextDocument: lsp.TextDocumentIdentifier{
					URI: "file:///test.go",
				},
				Position: lsp.Position{
					Line:      10,
					Character: 5,
				},
			},
		}

		data, err := json.Marshal(params)
		if err != nil {
			t.Fatalf("Failed to marshal CallHierarchyPrepareParams: %v", err)
		}

		var unmarshaled lsp.CallHierarchyPrepareParams
		err = json.Unmarshal(data, &unmarshaled)
		if err != nil {
			t.Fatalf("Failed to unmarshal CallHierarchyPrepareParams: %v", err)
		}

		if diff := cmp.Diff(params, unmarshaled); diff != "" {
			t.Errorf("CallHierarchyPrepareParams mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("CallHierarchyItem", func(t *testing.T) {
		item := lsp.CallHierarchyItem{
			Name:   "main",
			Kind:   lsp.Function,
			Detail: "func main()",
			URI:    "file:///test.go",
			Range: lsp.Range{
				Start: lsp.Position{Line: 10, Character: 0},
				End:   lsp.Position{Line: 15, Character: 1},
			},
			SelectionRange: lsp.Range{
				Start: lsp.Position{Line: 10, Character: 5},
				End:   lsp.Position{Line: 10, Character: 9},
			},
		}

		data, err := json.Marshal(item)
		if err != nil {
			t.Fatalf("Failed to marshal CallHierarchyItem: %v", err)
		}

		var unmarshaled lsp.CallHierarchyItem
		err = json.Unmarshal(data, &unmarshaled)
		if err != nil {
			t.Fatalf("Failed to unmarshal CallHierarchyItem: %v", err)
		}

		if diff := cmp.Diff(item, unmarshaled); diff != "" {
			t.Errorf("CallHierarchyItem mismatch (-want +got):\n%s", diff)
		}
	})
}

// TestInlayHints tests inlay hints scenarios
func TestInlayHints(t *testing.T) {
	t.Run("InlayHintParams", func(t *testing.T) {
		params := lsp.InlayHintParams{
			TextDocument: lsp.TextDocumentIdentifier{
				URI: "file:///test.go",
			},
			Range: lsp.Range{
				Start: lsp.Position{Line: 0, Character: 0},
				End:   lsp.Position{Line: 50, Character: 0},
			},
		}

		data, err := json.Marshal(params)
		if err != nil {
			t.Fatalf("Failed to marshal InlayHintParams: %v", err)
		}

		var unmarshaled lsp.InlayHintParams
		err = json.Unmarshal(data, &unmarshaled)
		if err != nil {
			t.Fatalf("Failed to unmarshal InlayHintParams: %v", err)
		}

		if diff := cmp.Diff(params, unmarshaled); diff != "" {
			t.Errorf("InlayHintParams mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("InlayHint", func(t *testing.T) {
		hint := lsp.InlayHint{
			Position: lsp.Position{
				Line:      5,
				Character: 10,
			},
			Label: []lsp.InlayHintLabelPart{
				{
					Value: "string",
				},
			},
			Kind: lsp.Type,
		}

		data, err := json.Marshal(hint)
		if err != nil {
			t.Fatalf("Failed to marshal InlayHint: %v", err)
		}

		var unmarshaled lsp.InlayHint
		err = json.Unmarshal(data, &unmarshaled)
		if err != nil {
			t.Fatalf("Failed to unmarshal InlayHint: %v", err)
		}

		if diff := cmp.Diff(hint, unmarshaled); diff != "" {
			t.Errorf("InlayHint mismatch (-want +got):\n%s", diff)
		}
	})
}

// TestDocumentSymbol tests document symbol scenarios
func TestDocumentSymbol(t *testing.T) {
	t.Run("DocumentSymbolParams", func(t *testing.T) {
		params := lsp.DocumentSymbolParams{
			TextDocument: lsp.TextDocumentIdentifier{
				URI: "file:///test.go",
			},
		}

		data, err := json.Marshal(params)
		if err != nil {
			t.Fatalf("Failed to marshal DocumentSymbolParams: %v", err)
		}

		var unmarshaled lsp.DocumentSymbolParams
		err = json.Unmarshal(data, &unmarshaled)
		if err != nil {
			t.Fatalf("Failed to unmarshal DocumentSymbolParams: %v", err)
		}

		if diff := cmp.Diff(params, unmarshaled); diff != "" {
			t.Errorf("DocumentSymbolParams mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("DocumentSymbol", func(t *testing.T) {
		symbol := lsp.DocumentSymbol{
			Name:   "TestFunction",
			Detail: "func TestFunction()",
			Kind:   lsp.Function,
			Range: lsp.Range{
				Start: lsp.Position{Line: 10, Character: 0},
				End:   lsp.Position{Line: 20, Character: 1},
			},
			SelectionRange: lsp.Range{
				Start: lsp.Position{Line: 10, Character: 5},
				End:   lsp.Position{Line: 10, Character: 17},
			},
			Children: []lsp.DocumentSymbol{
				{
					Name: "localVar",
					Kind: lsp.Variable,
					Range: lsp.Range{
						Start: lsp.Position{Line: 11, Character: 1},
						End:   lsp.Position{Line: 11, Character: 15},
					},
					SelectionRange: lsp.Range{
						Start: lsp.Position{Line: 11, Character: 1},
						End:   lsp.Position{Line: 11, Character: 9},
					},
				},
			},
		}

		data, err := json.Marshal(symbol)
		if err != nil {
			t.Fatalf("Failed to marshal DocumentSymbol: %v", err)
		}

		var unmarshaled lsp.DocumentSymbol
		err = json.Unmarshal(data, &unmarshaled)
		if err != nil {
			t.Fatalf("Failed to unmarshal DocumentSymbol: %v", err)
		}

		if diff := cmp.Diff(symbol, unmarshaled); diff != "" {
			t.Errorf("DocumentSymbol mismatch (-want +got):\n%s", diff)
		}
	})
}
