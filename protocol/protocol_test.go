// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package protocol_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/TypeFox/go-lsp/protocol"
	"github.com/google/go-cmp/cmp"
)

// TestJSONMarshaling tests JSON marshaling/unmarshaling of basic LSP types
func TestJSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  interface{}
	}{
		{
			name: "Position",
			obj: protocol.Position{
				Line:      10,
				Character: 5,
			},
		},
		{
			name: "Range",
			obj: protocol.Range{
				Start: protocol.Position{Line: 1, Character: 0},
				End:   protocol.Position{Line: 1, Character: 10},
			},
		},
		{
			name: "TextDocumentIdentifier",
			obj: protocol.TextDocumentIdentifier{
				URI: "file:///path/to/file.go",
			},
		},
		{
			name: "TextEdit",
			obj: protocol.TextEdit{
				Range: protocol.Range{
					Start: protocol.Position{Line: 0, Character: 0},
					End:   protocol.Position{Line: 0, Character: 5},
				},
				NewText: "hello",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal to JSON
			data, err := json.Marshal(tt.obj)
			if err != nil {
				t.Fatalf("Failed to marshal %s: %v", tt.name, err)
			}

			// Unmarshal back
			result := reflect.New(reflect.TypeOf(tt.obj)).Interface()
			err = json.Unmarshal(data, result)
			if err != nil {
				t.Fatalf("Failed to unmarshal %s: %v", tt.name, err)
			}

			// Compare
			resultValue := reflect.ValueOf(result).Elem().Interface()
			if diff := cmp.Diff(tt.obj, resultValue); diff != "" {
				t.Errorf("Marshaling/unmarshaling mismatch for %s (-want +got):\n%s", tt.name, diff)
			}
		})
	}
}