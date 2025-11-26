// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lsp

import (
	"fmt"
	"slices"
	"sort"
)

// Edit describes the replacement of a portion of a text file.
type Edit struct {
	Start, End int    // byte offsets of the region to replace
	New        string // the replacement
}

// apply applies a sequence of edits to the src buffer and returns the
// result. Edits are applied in order of start offset; edits with the
// same start offset are applied in they order they were provided.
func apply(src string, edits []Edit) (string, error) {
	edits, size, err := validate(src, edits)
	if err != nil {
		return "", err
	}

	// Apply edits.
	out := make([]byte, 0, size)
	lastEnd := 0
	for _, edit := range edits {
		if lastEnd < edit.Start {
			out = append(out, src[lastEnd:edit.Start]...)
		}
		out = append(out, edit.New...)
		lastEnd = edit.End
	}
	out = append(out, src[lastEnd:]...)

	if len(out) != size {
		panic("wrong size")
	}
	return string(out), nil
}

// applyBytes is like apply, but it accepts a byte slice.
func applyBytes(src []byte, edits []Edit) ([]byte, error) {
	res, err := apply(string(src), edits)
	return []byte(res), err
}

// validate checks that edits are consistent with src,
// and returns the size of the patched output.
func validate(src string, edits []Edit) ([]Edit, int, error) {
	if !sort.IsSorted(editsSort(edits)) {
		edits = slices.Clone(edits)
		sortEdits(edits)
	}

	// Check validity of edits and compute final size.
	size := len(src)
	lastEnd := 0
	for _, edit := range edits {
		if edit.Start < 0 || edit.Start > edit.End || edit.End > len(src) {
			return nil, 0, fmt.Errorf("diff has out-of-bounds edits")
		}
		if edit.Start < lastEnd {
			return nil, 0, fmt.Errorf("diff has overlapping edits")
		}
		size += len(edit.New) + edit.Start - edit.End
		lastEnd = edit.End
	}

	return edits, size, nil
}

// sortEdits orders a slice of Edits by (start, end) offset.
func sortEdits(edits []Edit) {
	sort.Stable(editsSort(edits))
}

type editsSort []Edit

func (a editsSort) Len() int { return len(a) }
func (a editsSort) Less(i, j int) bool {
	if cmp := a[i].Start - a[j].Start; cmp != 0 {
		return cmp < 0
	}
	return a[i].End < a[j].End
}
func (a editsSort) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// EditsFromDiffEdits converts Edits to a non-nil slice of LSP TextEdits.
// See https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#textEditArray
func EditsFromDiffEdits(m *Mapper, edits []Edit) ([]TextEdit, error) {
	// LSP doesn't require TextEditArray to be sorted:
	// this is the receiver's concern. But govim, and perhaps
	// other clients have historically relied on the order.
	edits = slices.Clone(edits)
	sortEdits(edits)

	result := make([]TextEdit, len(edits))
	for i, edit := range edits {
		rng, err := m.OffsetRange(edit.Start, edit.End)
		if err != nil {
			return nil, err
		}
		result[i] = TextEdit{
			Range:   rng,
			NewText: edit.New,
		}
	}
	return result, nil
}

// EditsToDiffEdits converts LSP TextEdits to Edits.
// See https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#textEditArray
func EditsToDiffEdits(m *Mapper, edits []TextEdit) ([]Edit, error) {
	if edits == nil {
		return nil, nil
	}
	result := make([]Edit, len(edits))
	for i, edit := range edits {
		start, end, err := m.RangeOffsets(edit.Range)
		if err != nil {
			return nil, err
		}
		result[i] = Edit{
			Start: start,
			End:   end,
			New:   edit.NewText,
		}
	}
	return result, nil
}

// ApplyEdits applies the patch (edits) to m.Content and returns the result.
// It also returns the edits converted to Edit form.
func ApplyEdits(m *Mapper, edits []TextEdit) ([]byte, []Edit, error) {
	diffEdits, err := EditsToDiffEdits(m, edits)
	if err != nil {
		return nil, nil, err
	}
	out, err := applyBytes(m.Content, diffEdits)
	return out, diffEdits, err
}

// AsTextEdits converts a slice possibly containing AnnotatedTextEdits
// to a slice of TextEdits.
func AsTextEdits(edits []Or_TextDocumentEdit_edits_Elem) []TextEdit {
	var result []TextEdit
	for _, e := range edits {
		var te TextEdit
		if x, ok := e.Value.(AnnotatedTextEdit); ok {
			te = x.TextEdit
		} else if x, ok := e.Value.(TextEdit); ok {
			te = x
		} else {
			panic(fmt.Sprintf("unexpected type %T, expected AnnotatedTextEdit or TextEdit", e.Value))
		}
		result = append(result, te)
	}
	return result
}

// AsAnnotatedTextEdits converts a slice of TextEdits
// to a slice of Or_TextDocumentEdit_edits_Elem.
// (returning a typed nil is required in server: in code_action.go and command.go))
func AsAnnotatedTextEdits(edits []TextEdit) []Or_TextDocumentEdit_edits_Elem {
	if edits == nil {
		return []Or_TextDocumentEdit_edits_Elem{}
	}
	var result []Or_TextDocumentEdit_edits_Elem
	for _, e := range edits {
		result = append(result, Or_TextDocumentEdit_edits_Elem{
			Value: TextEdit{
				Range:   e.Range,
				NewText: e.NewText,
			},
		})
	}
	return result
}

// fileHandle abstracts file.Handle to avoid a cycle.
type fileHandle interface {
	URI() DocumentURI
	Version() int32
}

// NewWorkspaceEdit constructs a WorkspaceEdit from a list of document changes.
//
// Any ChangeAnnotations must be added after.
func NewWorkspaceEdit(changes ...DocumentChange) *WorkspaceEdit {
	return &WorkspaceEdit{DocumentChanges: changes}
}

// DocumentChangeEdit constructs a DocumentChange containing a
// TextDocumentEdit from a file.Handle and a list of TextEdits.
func DocumentChangeEdit(fh fileHandle, textedits []TextEdit) DocumentChange {
	return DocumentChange{
		TextDocumentEdit: &TextDocumentEdit{
			TextDocument: OptionalVersionedTextDocumentIdentifier{
				Version:                fh.Version(),
				TextDocumentIdentifier: TextDocumentIdentifier{URI: fh.URI()},
			},
			Edits: AsAnnotatedTextEdits(textedits),
		},
	}
}

// DocumentChangeCreate constructs a DocumentChange that creates a file.
func DocumentChangeCreate(uri DocumentURI) DocumentChange {
	return DocumentChange{
		CreateFile: &CreateFile{
			Kind: "create",
			URI:  uri,
		},
	}
}

// DocumentChangeRename constructs a DocumentChange that renames a file.
func DocumentChangeRename(src, dst DocumentURI) DocumentChange {
	return DocumentChange{
		RenameFile: &RenameFile{
			Kind:   "rename",
			OldURI: src,
			NewURI: dst,
		},
	}
}

// SelectCompletionTextEdit returns insert or replace mode TextEdit
// included in the completion item.
func SelectCompletionTextEdit(item CompletionItem, useReplaceMode bool) (TextEdit, error) {
	var edit TextEdit
	switch typ := item.TextEdit.Value.(type) {
	case TextEdit: // old style completion item.
		return typ, nil
	case InsertReplaceEdit:
		if useReplaceMode {
			return TextEdit{
				NewText: typ.NewText,
				Range:   typ.Replace,
			}, nil
		} else {
			return TextEdit{
				NewText: typ.NewText,
				Range:   typ.Insert,
			}, nil
		}
	default:
		return edit, fmt.Errorf("unsupported edit type %T", typ)
	}
}
