// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated for LSP. DO NOT EDIT.

package lsp

// Code generated from protocol/metaModel.json at ref release/protocol/3.18.1 (hash bb5ee9298f3b0881df78c35e5762512d2c922484).
// https://github.com/microsoft/vscode-languageserver-node/blob/release/protocol/3.18.1/protocol/metaModel.json
// LSP metaData.version = 3.18.0.

import "encoding/json"

// UnmarshalError indicates that a JSON value did not conform to
// one of the expected cases of an LSP union type.
type UnmarshalError struct {
	msg string
}

func (e UnmarshalError) Error() string {
	return e.msg
}
func (t CancelParamsId) MarshalJSON() ([]byte, error) {
	switch {
	case t.Int32 != nil:
		return json.Marshal(*t.Int32)
	case t.String != nil:
		return json.Marshal(*t.String)
	}
	return []byte("null"), nil
}

func (t *CancelParamsId) UnmarshalJSON(x []byte) error {
	*t = CancelParamsId{}
	if string(x) == "null" {
		return nil
	}
	var h0 int32
	if err := json.Unmarshal(x, &h0); err == nil {
		t.Int32 = &h0
		return nil
	}
	var h1 string
	if err := json.Unmarshal(x, &h1); err == nil {
		t.String = &h1
		return nil
	}
	return &UnmarshalError{"unmarshal failed to match one of [int32 string]"}
}

func (t ClientSemanticTokensRequestOptionsFull) MarshalJSON() ([]byte, error) {
	switch {
	case t.ClientSemanticTokensRequestFullDelta != nil:
		return json.Marshal(*t.ClientSemanticTokensRequestFullDelta)
	case t.Bool != nil:
		return json.Marshal(*t.Bool)
	}
	return []byte("null"), nil
}

func (t *ClientSemanticTokensRequestOptionsFull) UnmarshalJSON(x []byte) error {
	*t = ClientSemanticTokensRequestOptionsFull{}
	if string(x) == "null" {
		return nil
	}
	var h0 ClientSemanticTokensRequestFullDelta
	if err := json.Unmarshal(x, &h0); err == nil {
		t.ClientSemanticTokensRequestFullDelta = &h0
		return nil
	}
	var h1 bool
	if err := json.Unmarshal(x, &h1); err == nil {
		t.Bool = &h1
		return nil
	}
	return &UnmarshalError{"unmarshal failed to match one of [ClientSemanticTokensRequestFullDelta bool]"}
}

func (t ClientSemanticTokensRequestOptionsRange) MarshalJSON() ([]byte, error) {
	switch {
	case t.Lit_ClientSemanticTokensRequestOptions_range_Item1 != nil:
		return json.Marshal(*t.Lit_ClientSemanticTokensRequestOptions_range_Item1)
	case t.Bool != nil:
		return json.Marshal(*t.Bool)
	}
	return []byte("null"), nil
}

func (t *ClientSemanticTokensRequestOptionsRange) UnmarshalJSON(x []byte) error {
	*t = ClientSemanticTokensRequestOptionsRange{}
	if string(x) == "null" {
		return nil
	}
	var h0 Lit_ClientSemanticTokensRequestOptions_range_Item1
	if err := json.Unmarshal(x, &h0); err == nil {
		t.Lit_ClientSemanticTokensRequestOptions_range_Item1 = &h0
		return nil
	}
	var h1 bool
	if err := json.Unmarshal(x, &h1); err == nil {
		t.Bool = &h1
		return nil
	}
	return &UnmarshalError{"unmarshal failed to match one of [Lit_ClientSemanticTokensRequestOptions_range_Item1 bool]"}
}

func (t CompletionItemDefaultsEditRange) MarshalJSON() ([]byte, error) {
	switch {
	case t.EditRangeWithInsertReplace != nil:
		return json.Marshal(*t.EditRangeWithInsertReplace)
	case t.Range != nil:
		return json.Marshal(*t.Range)
	}
	return []byte("null"), nil
}

func (t *CompletionItemDefaultsEditRange) UnmarshalJSON(x []byte) error {
	*t = CompletionItemDefaultsEditRange{}
	if string(x) == "null" {
		return nil
	}
	var h0 EditRangeWithInsertReplace
	if err := json.Unmarshal(x, &h0); err == nil {
		t.EditRangeWithInsertReplace = &h0
		return nil
	}
	var h1 Range
	if err := json.Unmarshal(x, &h1); err == nil {
		t.Range = &h1
		return nil
	}
	return &UnmarshalError{"unmarshal failed to match one of [EditRangeWithInsertReplace Range]"}
}

func (t CompletionItemDocumentation) MarshalJSON() ([]byte, error) {
	switch {
	case t.MarkupContent != nil:
		return json.Marshal(*t.MarkupContent)
	case t.String != nil:
		return json.Marshal(*t.String)
	}
	return []byte("null"), nil
}

func (t *CompletionItemDocumentation) UnmarshalJSON(x []byte) error {
	*t = CompletionItemDocumentation{}
	if string(x) == "null" {
		return nil
	}
	var h0 MarkupContent
	if err := json.Unmarshal(x, &h0); err == nil {
		t.MarkupContent = &h0
		return nil
	}
	var h1 string
	if err := json.Unmarshal(x, &h1); err == nil {
		t.String = &h1
		return nil
	}
	return &UnmarshalError{"unmarshal failed to match one of [MarkupContent string]"}
}

func (t CompletionItemTextEdit) MarshalJSON() ([]byte, error) {
	switch {
	case t.InsertReplaceEdit != nil:
		return json.Marshal(*t.InsertReplaceEdit)
	case t.TextEdit != nil:
		return json.Marshal(*t.TextEdit)
	}
	return []byte("null"), nil
}

func (t *CompletionItemTextEdit) UnmarshalJSON(x []byte) error {
	*t = CompletionItemTextEdit{}
	if string(x) == "null" {
		return nil
	}
	var h0 InsertReplaceEdit
	if err := json.Unmarshal(x, &h0); err == nil {
		t.InsertReplaceEdit = &h0
		return nil
	}
	var h1 TextEdit
	if err := json.Unmarshal(x, &h1); err == nil {
		t.TextEdit = &h1
		return nil
	}
	return &UnmarshalError{"unmarshal failed to match one of [InsertReplaceEdit TextEdit]"}
}

func (t Definition) MarshalJSON() ([]byte, error) {
	switch {
	case t.Location != nil:
		return json.Marshal(*t.Location)
	case t.Locations != nil:
		return json.Marshal(*t.Locations)
	}
	return []byte("null"), nil
}

func (t *Definition) UnmarshalJSON(x []byte) error {
	*t = Definition{}
	if string(x) == "null" {
		return nil
	}
	var h0 Location
	if err := json.Unmarshal(x, &h0); err == nil {
		t.Location = &h0
		return nil
	}
	var h1 []Location
	if err := json.Unmarshal(x, &h1); err == nil {
		t.Locations = &h1
		return nil
	}
	return &UnmarshalError{"unmarshal failed to match one of [Location []Location]"}
}

func (t DiagnosticCode) MarshalJSON() ([]byte, error) {
	switch {
	case t.Int32 != nil:
		return json.Marshal(*t.Int32)
	case t.String != nil:
		return json.Marshal(*t.String)
	}
	return []byte("null"), nil
}

func (t *DiagnosticCode) UnmarshalJSON(x []byte) error {
	*t = DiagnosticCode{}
	if string(x) == "null" {
		return nil
	}
	var h0 int32
	if err := json.Unmarshal(x, &h0); err == nil {
		t.Int32 = &h0
		return nil
	}
	var h1 string
	if err := json.Unmarshal(x, &h1); err == nil {
		t.String = &h1
		return nil
	}
	return &UnmarshalError{"unmarshal failed to match one of [int32 string]"}
}

func (t DiagnosticMessage) MarshalJSON() ([]byte, error) {
	switch {
	case t.MarkupContent != nil:
		return json.Marshal(*t.MarkupContent)
	case t.String != nil:
		return json.Marshal(*t.String)
	}
	return []byte("null"), nil
}

func (t *DiagnosticMessage) UnmarshalJSON(x []byte) error {
	*t = DiagnosticMessage{}
	if string(x) == "null" {
		return nil
	}
	var h0 MarkupContent
	if err := json.Unmarshal(x, &h0); err == nil {
		t.MarkupContent = &h0
		return nil
	}
	var h1 string
	if err := json.Unmarshal(x, &h1); err == nil {
		t.String = &h1
		return nil
	}
	return &UnmarshalError{"unmarshal failed to match one of [MarkupContent string]"}
}

func (t DidChangeConfigurationRegistrationOptionsSection) MarshalJSON() ([]byte, error) {
	switch {
	case t.Strings != nil:
		return json.Marshal(*t.Strings)
	case t.String != nil:
		return json.Marshal(*t.String)
	}
	return []byte("null"), nil
}

func (t *DidChangeConfigurationRegistrationOptionsSection) UnmarshalJSON(x []byte) error {
	*t = DidChangeConfigurationRegistrationOptionsSection{}
	if string(x) == "null" {
		return nil
	}
	var h0 []string
	if err := json.Unmarshal(x, &h0); err == nil {
		t.Strings = &h0
		return nil
	}
	var h1 string
	if err := json.Unmarshal(x, &h1); err == nil {
		t.String = &h1
		return nil
	}
	return &UnmarshalError{"unmarshal failed to match one of [[]string string]"}
}

func (t DocumentDiagnosticReport) MarshalJSON() ([]byte, error) {
	switch {
	case t.RelatedFullDocumentDiagnosticReport != nil:
		return json.Marshal(*t.RelatedFullDocumentDiagnosticReport)
	case t.RelatedUnchangedDocumentDiagnosticReport != nil:
		return json.Marshal(*t.RelatedUnchangedDocumentDiagnosticReport)
	}
	return []byte("null"), nil
}

func (t *DocumentDiagnosticReport) UnmarshalJSON(x []byte) error {
	*t = DocumentDiagnosticReport{}
	if string(x) == "null" {
		return nil
	}
	var h0 RelatedFullDocumentDiagnosticReport
	if err := json.Unmarshal(x, &h0); err == nil {
		t.RelatedFullDocumentDiagnosticReport = &h0
		return nil
	}
	var h1 RelatedUnchangedDocumentDiagnosticReport
	if err := json.Unmarshal(x, &h1); err == nil {
		t.RelatedUnchangedDocumentDiagnosticReport = &h1
		return nil
	}
	return &UnmarshalError{"unmarshal failed to match one of [RelatedFullDocumentDiagnosticReport RelatedUnchangedDocumentDiagnosticReport]"}
}

func (t DocumentFilter) MarshalJSON() ([]byte, error) {
	switch {
	case t.NotebookCellTextDocumentFilter != nil:
		return json.Marshal(*t.NotebookCellTextDocumentFilter)
	case t.TextDocumentFilter != nil:
		return json.Marshal(*t.TextDocumentFilter)
	}
	return []byte("null"), nil
}

func (t *DocumentFilter) UnmarshalJSON(x []byte) error {
	*t = DocumentFilter{}
	if string(x) == "null" {
		return nil
	}
	var h0 NotebookCellTextDocumentFilter
	if err := json.Unmarshal(x, &h0); err == nil {
		t.NotebookCellTextDocumentFilter = &h0
		return nil
	}
	var h1 TextDocumentFilter
	if err := json.Unmarshal(x, &h1); err == nil {
		t.TextDocumentFilter = &h1
		return nil
	}
	return &UnmarshalError{"unmarshal failed to match one of [NotebookCellTextDocumentFilter TextDocumentFilter]"}
}

func (t GlobPattern) MarshalJSON() ([]byte, error) {
	switch {
	case t.Pattern != nil:
		return json.Marshal(*t.Pattern)
	case t.RelativePattern != nil:
		return json.Marshal(*t.RelativePattern)
	}
	return []byte("null"), nil
}

func (t *GlobPattern) UnmarshalJSON(x []byte) error {
	*t = GlobPattern{}
	if string(x) == "null" {
		return nil
	}
	var h0 Pattern
	if err := json.Unmarshal(x, &h0); err == nil {
		t.Pattern = &h0
		return nil
	}
	var h1 RelativePattern
	if err := json.Unmarshal(x, &h1); err == nil {
		t.RelativePattern = &h1
		return nil
	}
	return &UnmarshalError{"unmarshal failed to match one of [Pattern RelativePattern]"}
}

func (t InlayHintLabelPartTooltip) MarshalJSON() ([]byte, error) {
	switch {
	case t.MarkupContent != nil:
		return json.Marshal(*t.MarkupContent)
	case t.String != nil:
		return json.Marshal(*t.String)
	}
	return []byte("null"), nil
}

func (t *InlayHintLabelPartTooltip) UnmarshalJSON(x []byte) error {
	*t = InlayHintLabelPartTooltip{}
	if string(x) == "null" {
		return nil
	}
	var h0 MarkupContent
	if err := json.Unmarshal(x, &h0); err == nil {
		t.MarkupContent = &h0
		return nil
	}
	var h1 string
	if err := json.Unmarshal(x, &h1); err == nil {
		t.String = &h1
		return nil
	}
	return &UnmarshalError{"unmarshal failed to match one of [MarkupContent string]"}
}

func (t InlayHintTooltip) MarshalJSON() ([]byte, error) {
	switch {
	case t.MarkupContent != nil:
		return json.Marshal(*t.MarkupContent)
	case t.String != nil:
		return json.Marshal(*t.String)
	}
	return []byte("null"), nil
}

func (t *InlayHintTooltip) UnmarshalJSON(x []byte) error {
	*t = InlayHintTooltip{}
	if string(x) == "null" {
		return nil
	}
	var h0 MarkupContent
	if err := json.Unmarshal(x, &h0); err == nil {
		t.MarkupContent = &h0
		return nil
	}
	var h1 string
	if err := json.Unmarshal(x, &h1); err == nil {
		t.String = &h1
		return nil
	}
	return &UnmarshalError{"unmarshal failed to match one of [MarkupContent string]"}
}

func (t InlineCompletionItemInsertText) MarshalJSON() ([]byte, error) {
	switch {
	case t.StringValue != nil:
		return json.Marshal(*t.StringValue)
	case t.String != nil:
		return json.Marshal(*t.String)
	}
	return []byte("null"), nil
}

func (t *InlineCompletionItemInsertText) UnmarshalJSON(x []byte) error {
	*t = InlineCompletionItemInsertText{}
	if string(x) == "null" {
		return nil
	}
	var h0 StringValue
	if err := json.Unmarshal(x, &h0); err == nil {
		t.StringValue = &h0
		return nil
	}
	var h1 string
	if err := json.Unmarshal(x, &h1); err == nil {
		t.String = &h1
		return nil
	}
	return &UnmarshalError{"unmarshal failed to match one of [StringValue string]"}
}

func (t InlineValue) MarshalJSON() ([]byte, error) {
	switch {
	case t.InlineValueEvaluatableExpression != nil:
		return json.Marshal(*t.InlineValueEvaluatableExpression)
	case t.InlineValueText != nil:
		return json.Marshal(*t.InlineValueText)
	case t.InlineValueVariableLookup != nil:
		return json.Marshal(*t.InlineValueVariableLookup)
	}
	return []byte("null"), nil
}

func (t *InlineValue) UnmarshalJSON(x []byte) error {
	*t = InlineValue{}
	if string(x) == "null" {
		return nil
	}
	var h0 InlineValueEvaluatableExpression
	if err := json.Unmarshal(x, &h0); err == nil {
		t.InlineValueEvaluatableExpression = &h0
		return nil
	}
	var h1 InlineValueText
	if err := json.Unmarshal(x, &h1); err == nil {
		t.InlineValueText = &h1
		return nil
	}
	var h2 InlineValueVariableLookup
	if err := json.Unmarshal(x, &h2); err == nil {
		t.InlineValueVariableLookup = &h2
		return nil
	}
	return &UnmarshalError{"unmarshal failed to match one of [InlineValueEvaluatableExpression InlineValueText InlineValueVariableLookup]"}
}

func (t NotebookCellTextDocumentFilterNotebook) MarshalJSON() ([]byte, error) {
	switch {
	case t.NotebookDocumentFilter != nil:
		return json.Marshal(*t.NotebookDocumentFilter)
	case t.String != nil:
		return json.Marshal(*t.String)
	}
	return []byte("null"), nil
}

func (t *NotebookCellTextDocumentFilterNotebook) UnmarshalJSON(x []byte) error {
	*t = NotebookCellTextDocumentFilterNotebook{}
	if string(x) == "null" {
		return nil
	}
	var h0 NotebookDocumentFilter
	if err := json.Unmarshal(x, &h0); err == nil {
		t.NotebookDocumentFilter = &h0
		return nil
	}
	var h1 string
	if err := json.Unmarshal(x, &h1); err == nil {
		t.String = &h1
		return nil
	}
	return &UnmarshalError{"unmarshal failed to match one of [NotebookDocumentFilter string]"}
}

func (t NotebookDocumentFilter) MarshalJSON() ([]byte, error) {
	switch {
	case t.NotebookDocumentFilterNotebookType != nil:
		return json.Marshal(*t.NotebookDocumentFilterNotebookType)
	case t.NotebookDocumentFilterPattern != nil:
		return json.Marshal(*t.NotebookDocumentFilterPattern)
	case t.NotebookDocumentFilterScheme != nil:
		return json.Marshal(*t.NotebookDocumentFilterScheme)
	}
	return []byte("null"), nil
}

func (t *NotebookDocumentFilter) UnmarshalJSON(x []byte) error {
	*t = NotebookDocumentFilter{}
	if string(x) == "null" {
		return nil
	}
	var h0 NotebookDocumentFilterNotebookType
	if err := json.Unmarshal(x, &h0); err == nil {
		t.NotebookDocumentFilterNotebookType = &h0
		return nil
	}
	var h1 NotebookDocumentFilterPattern
	if err := json.Unmarshal(x, &h1); err == nil {
		t.NotebookDocumentFilterPattern = &h1
		return nil
	}
	var h2 NotebookDocumentFilterScheme
	if err := json.Unmarshal(x, &h2); err == nil {
		t.NotebookDocumentFilterScheme = &h2
		return nil
	}
	return &UnmarshalError{"unmarshal failed to match one of [NotebookDocumentFilterNotebookType NotebookDocumentFilterPattern NotebookDocumentFilterScheme]"}
}

func (t NotebookDocumentFilterWithCellsNotebook) MarshalJSON() ([]byte, error) {
	switch {
	case t.NotebookDocumentFilter != nil:
		return json.Marshal(*t.NotebookDocumentFilter)
	case t.String != nil:
		return json.Marshal(*t.String)
	}
	return []byte("null"), nil
}

func (t *NotebookDocumentFilterWithCellsNotebook) UnmarshalJSON(x []byte) error {
	*t = NotebookDocumentFilterWithCellsNotebook{}
	if string(x) == "null" {
		return nil
	}
	var h0 NotebookDocumentFilter
	if err := json.Unmarshal(x, &h0); err == nil {
		t.NotebookDocumentFilter = &h0
		return nil
	}
	var h1 string
	if err := json.Unmarshal(x, &h1); err == nil {
		t.String = &h1
		return nil
	}
	return &UnmarshalError{"unmarshal failed to match one of [NotebookDocumentFilter string]"}
}

func (t NotebookDocumentFilterWithNotebookNotebook) MarshalJSON() ([]byte, error) {
	switch {
	case t.NotebookDocumentFilter != nil:
		return json.Marshal(*t.NotebookDocumentFilter)
	case t.String != nil:
		return json.Marshal(*t.String)
	}
	return []byte("null"), nil
}

func (t *NotebookDocumentFilterWithNotebookNotebook) UnmarshalJSON(x []byte) error {
	*t = NotebookDocumentFilterWithNotebookNotebook{}
	if string(x) == "null" {
		return nil
	}
	var h0 NotebookDocumentFilter
	if err := json.Unmarshal(x, &h0); err == nil {
		t.NotebookDocumentFilter = &h0
		return nil
	}
	var h1 string
	if err := json.Unmarshal(x, &h1); err == nil {
		t.String = &h1
		return nil
	}
	return &UnmarshalError{"unmarshal failed to match one of [NotebookDocumentFilter string]"}
}

func (t NotebookDocumentSyncOptionsNotebookSelectorElem) MarshalJSON() ([]byte, error) {
	switch {
	case t.NotebookDocumentFilterWithCells != nil:
		return json.Marshal(*t.NotebookDocumentFilterWithCells)
	case t.NotebookDocumentFilterWithNotebook != nil:
		return json.Marshal(*t.NotebookDocumentFilterWithNotebook)
	}
	return []byte("null"), nil
}

func (t *NotebookDocumentSyncOptionsNotebookSelectorElem) UnmarshalJSON(x []byte) error {
	*t = NotebookDocumentSyncOptionsNotebookSelectorElem{}
	if string(x) == "null" {
		return nil
	}
	var h0 NotebookDocumentFilterWithCells
	if err := json.Unmarshal(x, &h0); err == nil {
		t.NotebookDocumentFilterWithCells = &h0
		return nil
	}
	var h1 NotebookDocumentFilterWithNotebook
	if err := json.Unmarshal(x, &h1); err == nil {
		t.NotebookDocumentFilterWithNotebook = &h1
		return nil
	}
	return &UnmarshalError{"unmarshal failed to match one of [NotebookDocumentFilterWithCells NotebookDocumentFilterWithNotebook]"}
}

func (t ResultTextDocumentInlineCompletion) MarshalJSON() ([]byte, error) {
	switch {
	case t.InlineCompletionList != nil:
		return json.Marshal(*t.InlineCompletionList)
	case t.InlineCompletionItems != nil:
		return json.Marshal(*t.InlineCompletionItems)
	}
	return []byte("null"), nil
}

func (t *ResultTextDocumentInlineCompletion) UnmarshalJSON(x []byte) error {
	*t = ResultTextDocumentInlineCompletion{}
	if string(x) == "null" {
		return nil
	}
	var h0 InlineCompletionList
	if err := json.Unmarshal(x, &h0); err == nil {
		t.InlineCompletionList = &h0
		return nil
	}
	var h1 []InlineCompletionItem
	if err := json.Unmarshal(x, &h1); err == nil {
		t.InlineCompletionItems = &h1
		return nil
	}
	return &UnmarshalError{"unmarshal failed to match one of [InlineCompletionList []InlineCompletionItem]"}
}

func (t SemanticTokensOptionsFull) MarshalJSON() ([]byte, error) {
	switch {
	case t.SemanticTokensFullDelta != nil:
		return json.Marshal(*t.SemanticTokensFullDelta)
	case t.Bool != nil:
		return json.Marshal(*t.Bool)
	}
	return []byte("null"), nil
}

func (t *SemanticTokensOptionsFull) UnmarshalJSON(x []byte) error {
	*t = SemanticTokensOptionsFull{}
	if string(x) == "null" {
		return nil
	}
	var h0 SemanticTokensFullDelta
	if err := json.Unmarshal(x, &h0); err == nil {
		t.SemanticTokensFullDelta = &h0
		return nil
	}
	var h1 bool
	if err := json.Unmarshal(x, &h1); err == nil {
		t.Bool = &h1
		return nil
	}
	return &UnmarshalError{"unmarshal failed to match one of [SemanticTokensFullDelta bool]"}
}

func (t SemanticTokensOptionsRange) MarshalJSON() ([]byte, error) {
	switch {
	case t.PRangeESemanticTokensOptions != nil:
		return json.Marshal(*t.PRangeESemanticTokensOptions)
	case t.Bool != nil:
		return json.Marshal(*t.Bool)
	}
	return []byte("null"), nil
}

func (t *SemanticTokensOptionsRange) UnmarshalJSON(x []byte) error {
	*t = SemanticTokensOptionsRange{}
	if string(x) == "null" {
		return nil
	}
	var h0 PRangeESemanticTokensOptions
	if err := json.Unmarshal(x, &h0); err == nil {
		t.PRangeESemanticTokensOptions = &h0
		return nil
	}
	var h1 bool
	if err := json.Unmarshal(x, &h1); err == nil {
		t.Bool = &h1
		return nil
	}
	return &UnmarshalError{"unmarshal failed to match one of [PRangeESemanticTokensOptions bool]"}
}

func (t ServerCapabilitiesDiagnosticProvider) MarshalJSON() ([]byte, error) {
	switch {
	case t.DiagnosticOptions != nil:
		return json.Marshal(*t.DiagnosticOptions)
	case t.DiagnosticRegistrationOptions != nil:
		return json.Marshal(*t.DiagnosticRegistrationOptions)
	}
	return []byte("null"), nil
}

func (t *ServerCapabilitiesDiagnosticProvider) UnmarshalJSON(x []byte) error {
	*t = ServerCapabilitiesDiagnosticProvider{}
	if string(x) == "null" {
		return nil
	}
	var h0 DiagnosticOptions
	if err := json.Unmarshal(x, &h0); err == nil {
		t.DiagnosticOptions = &h0
		return nil
	}
	var h1 DiagnosticRegistrationOptions
	if err := json.Unmarshal(x, &h1); err == nil {
		t.DiagnosticRegistrationOptions = &h1
		return nil
	}
	return &UnmarshalError{"unmarshal failed to match one of [DiagnosticOptions DiagnosticRegistrationOptions]"}
}

func (t SignatureInformationDocumentation) MarshalJSON() ([]byte, error) {
	switch {
	case t.MarkupContent != nil:
		return json.Marshal(*t.MarkupContent)
	case t.String != nil:
		return json.Marshal(*t.String)
	}
	return []byte("null"), nil
}

func (t *SignatureInformationDocumentation) UnmarshalJSON(x []byte) error {
	*t = SignatureInformationDocumentation{}
	if string(x) == "null" {
		return nil
	}
	var h0 MarkupContent
	if err := json.Unmarshal(x, &h0); err == nil {
		t.MarkupContent = &h0
		return nil
	}
	var h1 string
	if err := json.Unmarshal(x, &h1); err == nil {
		t.String = &h1
		return nil
	}
	return &UnmarshalError{"unmarshal failed to match one of [MarkupContent string]"}
}

func (t TextDocumentEditEditsElem) MarshalJSON() ([]byte, error) {
	switch {
	case t.AnnotatedTextEdit != nil:
		return json.Marshal(*t.AnnotatedTextEdit)
	case t.SnippetTextEdit != nil:
		return json.Marshal(*t.SnippetTextEdit)
	case t.TextEdit != nil:
		return json.Marshal(*t.TextEdit)
	}
	return []byte("null"), nil
}

func (t *TextDocumentEditEditsElem) UnmarshalJSON(x []byte) error {
	*t = TextDocumentEditEditsElem{}
	if string(x) == "null" {
		return nil
	}
	var h0 AnnotatedTextEdit
	if err := json.Unmarshal(x, &h0); err == nil {
		t.AnnotatedTextEdit = &h0
		return nil
	}
	var h1 SnippetTextEdit
	if err := json.Unmarshal(x, &h1); err == nil {
		t.SnippetTextEdit = &h1
		return nil
	}
	var h2 TextEdit
	if err := json.Unmarshal(x, &h2); err == nil {
		t.TextEdit = &h2
		return nil
	}
	return &UnmarshalError{"unmarshal failed to match one of [AnnotatedTextEdit SnippetTextEdit TextEdit]"}
}

func (t TextDocumentFilter) MarshalJSON() ([]byte, error) {
	switch {
	case t.TextDocumentFilterLanguage != nil:
		return json.Marshal(*t.TextDocumentFilterLanguage)
	case t.TextDocumentFilterPattern != nil:
		return json.Marshal(*t.TextDocumentFilterPattern)
	case t.TextDocumentFilterScheme != nil:
		return json.Marshal(*t.TextDocumentFilterScheme)
	}
	return []byte("null"), nil
}

func (t *TextDocumentFilter) UnmarshalJSON(x []byte) error {
	*t = TextDocumentFilter{}
	if string(x) == "null" {
		return nil
	}
	var h0 TextDocumentFilterLanguage
	if err := json.Unmarshal(x, &h0); err == nil {
		t.TextDocumentFilterLanguage = &h0
		return nil
	}
	var h1 TextDocumentFilterPattern
	if err := json.Unmarshal(x, &h1); err == nil {
		t.TextDocumentFilterPattern = &h1
		return nil
	}
	var h2 TextDocumentFilterScheme
	if err := json.Unmarshal(x, &h2); err == nil {
		t.TextDocumentFilterScheme = &h2
		return nil
	}
	return &UnmarshalError{"unmarshal failed to match one of [TextDocumentFilterLanguage TextDocumentFilterPattern TextDocumentFilterScheme]"}
}

func (t WorkspaceDocumentDiagnosticReport) MarshalJSON() ([]byte, error) {
	switch {
	case t.WorkspaceFullDocumentDiagnosticReport != nil:
		return json.Marshal(*t.WorkspaceFullDocumentDiagnosticReport)
	case t.WorkspaceUnchangedDocumentDiagnosticReport != nil:
		return json.Marshal(*t.WorkspaceUnchangedDocumentDiagnosticReport)
	}
	return []byte("null"), nil
}

func (t *WorkspaceDocumentDiagnosticReport) UnmarshalJSON(x []byte) error {
	*t = WorkspaceDocumentDiagnosticReport{}
	if string(x) == "null" {
		return nil
	}
	var h0 WorkspaceFullDocumentDiagnosticReport
	if err := json.Unmarshal(x, &h0); err == nil {
		t.WorkspaceFullDocumentDiagnosticReport = &h0
		return nil
	}
	var h1 WorkspaceUnchangedDocumentDiagnosticReport
	if err := json.Unmarshal(x, &h1); err == nil {
		t.WorkspaceUnchangedDocumentDiagnosticReport = &h1
		return nil
	}
	return &UnmarshalError{"unmarshal failed to match one of [WorkspaceFullDocumentDiagnosticReport WorkspaceUnchangedDocumentDiagnosticReport]"}
}

func (t WorkspaceOptionsTextDocumentContent) MarshalJSON() ([]byte, error) {
	switch {
	case t.TextDocumentContentOptions != nil:
		return json.Marshal(*t.TextDocumentContentOptions)
	case t.TextDocumentContentRegistrationOptions != nil:
		return json.Marshal(*t.TextDocumentContentRegistrationOptions)
	}
	return []byte("null"), nil
}

func (t *WorkspaceOptionsTextDocumentContent) UnmarshalJSON(x []byte) error {
	*t = WorkspaceOptionsTextDocumentContent{}
	if string(x) == "null" {
		return nil
	}
	var h0 TextDocumentContentOptions
	if err := json.Unmarshal(x, &h0); err == nil {
		t.TextDocumentContentOptions = &h0
		return nil
	}
	var h1 TextDocumentContentRegistrationOptions
	if err := json.Unmarshal(x, &h1); err == nil {
		t.TextDocumentContentRegistrationOptions = &h1
		return nil
	}
	return &UnmarshalError{"unmarshal failed to match one of [TextDocumentContentOptions TextDocumentContentRegistrationOptions]"}
}

func (t WorkspaceSymbolLocation) MarshalJSON() ([]byte, error) {
	switch {
	case t.Location != nil:
		return json.Marshal(*t.Location)
	case t.LocationUriOnly != nil:
		return json.Marshal(*t.LocationUriOnly)
	}
	return []byte("null"), nil
}

func (t *WorkspaceSymbolLocation) UnmarshalJSON(x []byte) error {
	*t = WorkspaceSymbolLocation{}
	if string(x) == "null" {
		return nil
	}
	var h0 Location
	if err := json.Unmarshal(x, &h0); err == nil {
		t.Location = &h0
		return nil
	}
	var h1 LocationUriOnly
	if err := json.Unmarshal(x, &h1); err == nil {
		t.LocationUriOnly = &h1
		return nil
	}
	return &UnmarshalError{"unmarshal failed to match one of [Location LocationUriOnly]"}
}
