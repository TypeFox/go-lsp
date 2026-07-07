// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated for LSP. DO NOT EDIT.

package lsp

// Code generated from protocol/metaModel.json at ref release/protocol/3.18.1 (hash bb5ee9298f3b0881df78c35e5762512d2c922484).
// https://github.com/microsoft/vscode-languageserver-node/blob/release/protocol/3.18.1/protocol/metaModel.json
// LSP metaData.version = 3.18.0.

// CancelParamsIdFromInt32 wraps a int32 value as a CancelParamsId union.
func CancelParamsIdFromInt32(v int32) CancelParamsId {
	return CancelParamsId{Int32: &v}
}

// CancelParamsIdFromString wraps a string value as a CancelParamsId union.
func CancelParamsIdFromString(v string) CancelParamsId {
	return CancelParamsId{String: &v}
}

// ClientSemanticTokensRequestOptionsFullFromBool wraps a bool value as a ClientSemanticTokensRequestOptionsFull union.
func ClientSemanticTokensRequestOptionsFullFromBool(v bool) ClientSemanticTokensRequestOptionsFull {
	return ClientSemanticTokensRequestOptionsFull{Bool: &v}
}

// ClientSemanticTokensRequestOptionsFullFromClientSemanticTokensRequestFullDelta wraps a ClientSemanticTokensRequestFullDelta value as a ClientSemanticTokensRequestOptionsFull union.
func ClientSemanticTokensRequestOptionsFullFromClientSemanticTokensRequestFullDelta(v ClientSemanticTokensRequestFullDelta) ClientSemanticTokensRequestOptionsFull {
	return ClientSemanticTokensRequestOptionsFull{ClientSemanticTokensRequestFullDelta: &v}
}

// ClientSemanticTokensRequestOptionsRangeFromBool wraps a bool value as a ClientSemanticTokensRequestOptionsRange union.
func ClientSemanticTokensRequestOptionsRangeFromBool(v bool) ClientSemanticTokensRequestOptionsRange {
	return ClientSemanticTokensRequestOptionsRange{Bool: &v}
}

// ClientSemanticTokensRequestOptionsRangeFromLit_ClientSemanticTokensRequestOptions_range_Item1 wraps a Lit_ClientSemanticTokensRequestOptions_range_Item1 value as a ClientSemanticTokensRequestOptionsRange union.
func ClientSemanticTokensRequestOptionsRangeFromLit_ClientSemanticTokensRequestOptions_range_Item1(v Lit_ClientSemanticTokensRequestOptions_range_Item1) ClientSemanticTokensRequestOptionsRange {
	return ClientSemanticTokensRequestOptionsRange{Lit_ClientSemanticTokensRequestOptions_range_Item1: &v}
}

// CompletionItemDefaultsEditRangeFromEditRangeWithInsertReplace wraps a EditRangeWithInsertReplace value as a CompletionItemDefaultsEditRange union.
func CompletionItemDefaultsEditRangeFromEditRangeWithInsertReplace(v EditRangeWithInsertReplace) CompletionItemDefaultsEditRange {
	return CompletionItemDefaultsEditRange{EditRangeWithInsertReplace: &v}
}

// CompletionItemDefaultsEditRangeFromRange wraps a Range value as a CompletionItemDefaultsEditRange union.
func CompletionItemDefaultsEditRangeFromRange(v Range) CompletionItemDefaultsEditRange {
	return CompletionItemDefaultsEditRange{Range: &v}
}

// CompletionItemDocumentationFromMarkupContent wraps a MarkupContent value as a CompletionItemDocumentation union.
func CompletionItemDocumentationFromMarkupContent(v MarkupContent) CompletionItemDocumentation {
	return CompletionItemDocumentation{MarkupContent: &v}
}

// CompletionItemDocumentationFromString wraps a string value as a CompletionItemDocumentation union.
func CompletionItemDocumentationFromString(v string) CompletionItemDocumentation {
	return CompletionItemDocumentation{String: &v}
}

// CompletionItemTextEditFromInsertReplaceEdit wraps a InsertReplaceEdit value as a CompletionItemTextEdit union.
func CompletionItemTextEditFromInsertReplaceEdit(v InsertReplaceEdit) CompletionItemTextEdit {
	return CompletionItemTextEdit{InsertReplaceEdit: &v}
}

// CompletionItemTextEditFromTextEdit wraps a TextEdit value as a CompletionItemTextEdit union.
func CompletionItemTextEditFromTextEdit(v TextEdit) CompletionItemTextEdit {
	return CompletionItemTextEdit{TextEdit: &v}
}

// DefinitionFromLocation wraps a Location value as a Definition union.
func DefinitionFromLocation(v Location) Definition {
	return Definition{Location: &v}
}

// DefinitionFromLocations wraps a []Location value as a Definition union.
func DefinitionFromLocations(v []Location) Definition {
	return Definition{Locations: &v}
}

// DiagnosticCodeFromInt32 wraps a int32 value as a DiagnosticCode union.
func DiagnosticCodeFromInt32(v int32) DiagnosticCode {
	return DiagnosticCode{Int32: &v}
}

// DiagnosticCodeFromString wraps a string value as a DiagnosticCode union.
func DiagnosticCodeFromString(v string) DiagnosticCode {
	return DiagnosticCode{String: &v}
}

// DiagnosticMessageFromMarkupContent wraps a MarkupContent value as a DiagnosticMessage union.
func DiagnosticMessageFromMarkupContent(v MarkupContent) DiagnosticMessage {
	return DiagnosticMessage{MarkupContent: &v}
}

// DiagnosticMessageFromString wraps a string value as a DiagnosticMessage union.
func DiagnosticMessageFromString(v string) DiagnosticMessage {
	return DiagnosticMessage{String: &v}
}

// DidChangeConfigurationRegistrationOptionsSectionFromString wraps a string value as a DidChangeConfigurationRegistrationOptionsSection union.
func DidChangeConfigurationRegistrationOptionsSectionFromString(v string) DidChangeConfigurationRegistrationOptionsSection {
	return DidChangeConfigurationRegistrationOptionsSection{String: &v}
}

// DidChangeConfigurationRegistrationOptionsSectionFromStrings wraps a []string value as a DidChangeConfigurationRegistrationOptionsSection union.
func DidChangeConfigurationRegistrationOptionsSectionFromStrings(v []string) DidChangeConfigurationRegistrationOptionsSection {
	return DidChangeConfigurationRegistrationOptionsSection{Strings: &v}
}

// DocumentDiagnosticReportFromRelatedFullDocumentDiagnosticReport wraps a RelatedFullDocumentDiagnosticReport value as a DocumentDiagnosticReport union.
func DocumentDiagnosticReportFromRelatedFullDocumentDiagnosticReport(v RelatedFullDocumentDiagnosticReport) DocumentDiagnosticReport {
	return DocumentDiagnosticReport{RelatedFullDocumentDiagnosticReport: &v}
}

// DocumentDiagnosticReportFromRelatedUnchangedDocumentDiagnosticReport wraps a RelatedUnchangedDocumentDiagnosticReport value as a DocumentDiagnosticReport union.
func DocumentDiagnosticReportFromRelatedUnchangedDocumentDiagnosticReport(v RelatedUnchangedDocumentDiagnosticReport) DocumentDiagnosticReport {
	return DocumentDiagnosticReport{RelatedUnchangedDocumentDiagnosticReport: &v}
}

// DocumentFilterFromNotebookCellTextDocumentFilter wraps a NotebookCellTextDocumentFilter value as a DocumentFilter union.
func DocumentFilterFromNotebookCellTextDocumentFilter(v NotebookCellTextDocumentFilter) DocumentFilter {
	return DocumentFilter{NotebookCellTextDocumentFilter: &v}
}

// DocumentFilterFromTextDocumentFilter wraps a TextDocumentFilter value as a DocumentFilter union.
func DocumentFilterFromTextDocumentFilter(v TextDocumentFilter) DocumentFilter {
	return DocumentFilter{TextDocumentFilter: &v}
}

// GlobPatternFromPattern wraps a Pattern value as a GlobPattern union.
func GlobPatternFromPattern(v Pattern) GlobPattern {
	return GlobPattern{Pattern: &v}
}

// GlobPatternFromRelativePattern wraps a RelativePattern value as a GlobPattern union.
func GlobPatternFromRelativePattern(v RelativePattern) GlobPattern {
	return GlobPattern{RelativePattern: &v}
}

// InlayHintLabelPartTooltipFromMarkupContent wraps a MarkupContent value as a InlayHintLabelPartTooltip union.
func InlayHintLabelPartTooltipFromMarkupContent(v MarkupContent) InlayHintLabelPartTooltip {
	return InlayHintLabelPartTooltip{MarkupContent: &v}
}

// InlayHintLabelPartTooltipFromString wraps a string value as a InlayHintLabelPartTooltip union.
func InlayHintLabelPartTooltipFromString(v string) InlayHintLabelPartTooltip {
	return InlayHintLabelPartTooltip{String: &v}
}

// InlayHintTooltipFromMarkupContent wraps a MarkupContent value as a InlayHintTooltip union.
func InlayHintTooltipFromMarkupContent(v MarkupContent) InlayHintTooltip {
	return InlayHintTooltip{MarkupContent: &v}
}

// InlayHintTooltipFromString wraps a string value as a InlayHintTooltip union.
func InlayHintTooltipFromString(v string) InlayHintTooltip {
	return InlayHintTooltip{String: &v}
}

// InlineCompletionItemInsertTextFromString wraps a string value as a InlineCompletionItemInsertText union.
func InlineCompletionItemInsertTextFromString(v string) InlineCompletionItemInsertText {
	return InlineCompletionItemInsertText{String: &v}
}

// InlineCompletionItemInsertTextFromStringValue wraps a StringValue value as a InlineCompletionItemInsertText union.
func InlineCompletionItemInsertTextFromStringValue(v StringValue) InlineCompletionItemInsertText {
	return InlineCompletionItemInsertText{StringValue: &v}
}

// InlineValueFromInlineValueEvaluatableExpression wraps a InlineValueEvaluatableExpression value as a InlineValue union.
func InlineValueFromInlineValueEvaluatableExpression(v InlineValueEvaluatableExpression) InlineValue {
	return InlineValue{InlineValueEvaluatableExpression: &v}
}

// InlineValueFromInlineValueText wraps a InlineValueText value as a InlineValue union.
func InlineValueFromInlineValueText(v InlineValueText) InlineValue {
	return InlineValue{InlineValueText: &v}
}

// InlineValueFromInlineValueVariableLookup wraps a InlineValueVariableLookup value as a InlineValue union.
func InlineValueFromInlineValueVariableLookup(v InlineValueVariableLookup) InlineValue {
	return InlineValue{InlineValueVariableLookup: &v}
}

// NotebookCellTextDocumentFilterNotebookFromNotebookDocumentFilter wraps a NotebookDocumentFilter value as a NotebookCellTextDocumentFilterNotebook union.
func NotebookCellTextDocumentFilterNotebookFromNotebookDocumentFilter(v NotebookDocumentFilter) NotebookCellTextDocumentFilterNotebook {
	return NotebookCellTextDocumentFilterNotebook{NotebookDocumentFilter: &v}
}

// NotebookCellTextDocumentFilterNotebookFromString wraps a string value as a NotebookCellTextDocumentFilterNotebook union.
func NotebookCellTextDocumentFilterNotebookFromString(v string) NotebookCellTextDocumentFilterNotebook {
	return NotebookCellTextDocumentFilterNotebook{String: &v}
}

// NotebookDocumentFilterFromNotebookDocumentFilterNotebookType wraps a NotebookDocumentFilterNotebookType value as a NotebookDocumentFilter union.
func NotebookDocumentFilterFromNotebookDocumentFilterNotebookType(v NotebookDocumentFilterNotebookType) NotebookDocumentFilter {
	return NotebookDocumentFilter{NotebookDocumentFilterNotebookType: &v}
}

// NotebookDocumentFilterFromNotebookDocumentFilterPattern wraps a NotebookDocumentFilterPattern value as a NotebookDocumentFilter union.
func NotebookDocumentFilterFromNotebookDocumentFilterPattern(v NotebookDocumentFilterPattern) NotebookDocumentFilter {
	return NotebookDocumentFilter{NotebookDocumentFilterPattern: &v}
}

// NotebookDocumentFilterFromNotebookDocumentFilterScheme wraps a NotebookDocumentFilterScheme value as a NotebookDocumentFilter union.
func NotebookDocumentFilterFromNotebookDocumentFilterScheme(v NotebookDocumentFilterScheme) NotebookDocumentFilter {
	return NotebookDocumentFilter{NotebookDocumentFilterScheme: &v}
}

// NotebookDocumentFilterWithCellsNotebookFromNotebookDocumentFilter wraps a NotebookDocumentFilter value as a NotebookDocumentFilterWithCellsNotebook union.
func NotebookDocumentFilterWithCellsNotebookFromNotebookDocumentFilter(v NotebookDocumentFilter) NotebookDocumentFilterWithCellsNotebook {
	return NotebookDocumentFilterWithCellsNotebook{NotebookDocumentFilter: &v}
}

// NotebookDocumentFilterWithCellsNotebookFromString wraps a string value as a NotebookDocumentFilterWithCellsNotebook union.
func NotebookDocumentFilterWithCellsNotebookFromString(v string) NotebookDocumentFilterWithCellsNotebook {
	return NotebookDocumentFilterWithCellsNotebook{String: &v}
}

// NotebookDocumentFilterWithNotebookNotebookFromNotebookDocumentFilter wraps a NotebookDocumentFilter value as a NotebookDocumentFilterWithNotebookNotebook union.
func NotebookDocumentFilterWithNotebookNotebookFromNotebookDocumentFilter(v NotebookDocumentFilter) NotebookDocumentFilterWithNotebookNotebook {
	return NotebookDocumentFilterWithNotebookNotebook{NotebookDocumentFilter: &v}
}

// NotebookDocumentFilterWithNotebookNotebookFromString wraps a string value as a NotebookDocumentFilterWithNotebookNotebook union.
func NotebookDocumentFilterWithNotebookNotebookFromString(v string) NotebookDocumentFilterWithNotebookNotebook {
	return NotebookDocumentFilterWithNotebookNotebook{String: &v}
}

// NotebookDocumentSyncOptionsNotebookSelectorElemFromNotebookDocumentFilterWithCells wraps a NotebookDocumentFilterWithCells value as a NotebookDocumentSyncOptionsNotebookSelectorElem union.
func NotebookDocumentSyncOptionsNotebookSelectorElemFromNotebookDocumentFilterWithCells(v NotebookDocumentFilterWithCells) NotebookDocumentSyncOptionsNotebookSelectorElem {
	return NotebookDocumentSyncOptionsNotebookSelectorElem{NotebookDocumentFilterWithCells: &v}
}

// NotebookDocumentSyncOptionsNotebookSelectorElemFromNotebookDocumentFilterWithNotebook wraps a NotebookDocumentFilterWithNotebook value as a NotebookDocumentSyncOptionsNotebookSelectorElem union.
func NotebookDocumentSyncOptionsNotebookSelectorElemFromNotebookDocumentFilterWithNotebook(v NotebookDocumentFilterWithNotebook) NotebookDocumentSyncOptionsNotebookSelectorElem {
	return NotebookDocumentSyncOptionsNotebookSelectorElem{NotebookDocumentFilterWithNotebook: &v}
}

// ResultTextDocumentInlineCompletionFromInlineCompletionItems wraps a []InlineCompletionItem value as a ResultTextDocumentInlineCompletion union.
func ResultTextDocumentInlineCompletionFromInlineCompletionItems(v []InlineCompletionItem) ResultTextDocumentInlineCompletion {
	return ResultTextDocumentInlineCompletion{InlineCompletionItems: &v}
}

// ResultTextDocumentInlineCompletionFromInlineCompletionList wraps a InlineCompletionList value as a ResultTextDocumentInlineCompletion union.
func ResultTextDocumentInlineCompletionFromInlineCompletionList(v InlineCompletionList) ResultTextDocumentInlineCompletion {
	return ResultTextDocumentInlineCompletion{InlineCompletionList: &v}
}

// SemanticTokensOptionsFullFromBool wraps a bool value as a SemanticTokensOptionsFull union.
func SemanticTokensOptionsFullFromBool(v bool) SemanticTokensOptionsFull {
	return SemanticTokensOptionsFull{Bool: &v}
}

// SemanticTokensOptionsFullFromSemanticTokensFullDelta wraps a SemanticTokensFullDelta value as a SemanticTokensOptionsFull union.
func SemanticTokensOptionsFullFromSemanticTokensFullDelta(v SemanticTokensFullDelta) SemanticTokensOptionsFull {
	return SemanticTokensOptionsFull{SemanticTokensFullDelta: &v}
}

// SemanticTokensOptionsRangeFromBool wraps a bool value as a SemanticTokensOptionsRange union.
func SemanticTokensOptionsRangeFromBool(v bool) SemanticTokensOptionsRange {
	return SemanticTokensOptionsRange{Bool: &v}
}

// SemanticTokensOptionsRangeFromPRangeESemanticTokensOptions wraps a PRangeESemanticTokensOptions value as a SemanticTokensOptionsRange union.
func SemanticTokensOptionsRangeFromPRangeESemanticTokensOptions(v PRangeESemanticTokensOptions) SemanticTokensOptionsRange {
	return SemanticTokensOptionsRange{PRangeESemanticTokensOptions: &v}
}

// ServerCapabilitiesDiagnosticProviderFromDiagnosticOptions wraps a DiagnosticOptions value as a ServerCapabilitiesDiagnosticProvider union.
func ServerCapabilitiesDiagnosticProviderFromDiagnosticOptions(v DiagnosticOptions) ServerCapabilitiesDiagnosticProvider {
	return ServerCapabilitiesDiagnosticProvider{DiagnosticOptions: &v}
}

// ServerCapabilitiesDiagnosticProviderFromDiagnosticRegistrationOptions wraps a DiagnosticRegistrationOptions value as a ServerCapabilitiesDiagnosticProvider union.
func ServerCapabilitiesDiagnosticProviderFromDiagnosticRegistrationOptions(v DiagnosticRegistrationOptions) ServerCapabilitiesDiagnosticProvider {
	return ServerCapabilitiesDiagnosticProvider{DiagnosticRegistrationOptions: &v}
}

// SignatureInformationDocumentationFromMarkupContent wraps a MarkupContent value as a SignatureInformationDocumentation union.
func SignatureInformationDocumentationFromMarkupContent(v MarkupContent) SignatureInformationDocumentation {
	return SignatureInformationDocumentation{MarkupContent: &v}
}

// SignatureInformationDocumentationFromString wraps a string value as a SignatureInformationDocumentation union.
func SignatureInformationDocumentationFromString(v string) SignatureInformationDocumentation {
	return SignatureInformationDocumentation{String: &v}
}

// TextDocumentEditEditsElemFromAnnotatedTextEdit wraps a AnnotatedTextEdit value as a TextDocumentEditEditsElem union.
func TextDocumentEditEditsElemFromAnnotatedTextEdit(v AnnotatedTextEdit) TextDocumentEditEditsElem {
	return TextDocumentEditEditsElem{AnnotatedTextEdit: &v}
}

// TextDocumentEditEditsElemFromSnippetTextEdit wraps a SnippetTextEdit value as a TextDocumentEditEditsElem union.
func TextDocumentEditEditsElemFromSnippetTextEdit(v SnippetTextEdit) TextDocumentEditEditsElem {
	return TextDocumentEditEditsElem{SnippetTextEdit: &v}
}

// TextDocumentEditEditsElemFromTextEdit wraps a TextEdit value as a TextDocumentEditEditsElem union.
func TextDocumentEditEditsElemFromTextEdit(v TextEdit) TextDocumentEditEditsElem {
	return TextDocumentEditEditsElem{TextEdit: &v}
}

// TextDocumentFilterFromTextDocumentFilterLanguage wraps a TextDocumentFilterLanguage value as a TextDocumentFilter union.
func TextDocumentFilterFromTextDocumentFilterLanguage(v TextDocumentFilterLanguage) TextDocumentFilter {
	return TextDocumentFilter{TextDocumentFilterLanguage: &v}
}

// TextDocumentFilterFromTextDocumentFilterPattern wraps a TextDocumentFilterPattern value as a TextDocumentFilter union.
func TextDocumentFilterFromTextDocumentFilterPattern(v TextDocumentFilterPattern) TextDocumentFilter {
	return TextDocumentFilter{TextDocumentFilterPattern: &v}
}

// TextDocumentFilterFromTextDocumentFilterScheme wraps a TextDocumentFilterScheme value as a TextDocumentFilter union.
func TextDocumentFilterFromTextDocumentFilterScheme(v TextDocumentFilterScheme) TextDocumentFilter {
	return TextDocumentFilter{TextDocumentFilterScheme: &v}
}

// WorkspaceDocumentDiagnosticReportFromWorkspaceFullDocumentDiagnosticReport wraps a WorkspaceFullDocumentDiagnosticReport value as a WorkspaceDocumentDiagnosticReport union.
func WorkspaceDocumentDiagnosticReportFromWorkspaceFullDocumentDiagnosticReport(v WorkspaceFullDocumentDiagnosticReport) WorkspaceDocumentDiagnosticReport {
	return WorkspaceDocumentDiagnosticReport{WorkspaceFullDocumentDiagnosticReport: &v}
}

// WorkspaceDocumentDiagnosticReportFromWorkspaceUnchangedDocumentDiagnosticReport wraps a WorkspaceUnchangedDocumentDiagnosticReport value as a WorkspaceDocumentDiagnosticReport union.
func WorkspaceDocumentDiagnosticReportFromWorkspaceUnchangedDocumentDiagnosticReport(v WorkspaceUnchangedDocumentDiagnosticReport) WorkspaceDocumentDiagnosticReport {
	return WorkspaceDocumentDiagnosticReport{WorkspaceUnchangedDocumentDiagnosticReport: &v}
}

// WorkspaceOptionsTextDocumentContentFromTextDocumentContentOptions wraps a TextDocumentContentOptions value as a WorkspaceOptionsTextDocumentContent union.
func WorkspaceOptionsTextDocumentContentFromTextDocumentContentOptions(v TextDocumentContentOptions) WorkspaceOptionsTextDocumentContent {
	return WorkspaceOptionsTextDocumentContent{TextDocumentContentOptions: &v}
}

// WorkspaceOptionsTextDocumentContentFromTextDocumentContentRegistrationOptions wraps a TextDocumentContentRegistrationOptions value as a WorkspaceOptionsTextDocumentContent union.
func WorkspaceOptionsTextDocumentContentFromTextDocumentContentRegistrationOptions(v TextDocumentContentRegistrationOptions) WorkspaceOptionsTextDocumentContent {
	return WorkspaceOptionsTextDocumentContent{TextDocumentContentRegistrationOptions: &v}
}

// WorkspaceSymbolLocationFromLocation wraps a Location value as a WorkspaceSymbolLocation union.
func WorkspaceSymbolLocationFromLocation(v Location) WorkspaceSymbolLocation {
	return WorkspaceSymbolLocation{Location: &v}
}

// WorkspaceSymbolLocationFromLocationUriOnly wraps a LocationUriOnly value as a WorkspaceSymbolLocation union.
func WorkspaceSymbolLocationFromLocationUriOnly(v LocationUriOnly) WorkspaceSymbolLocation {
	return WorkspaceSymbolLocation{LocationUriOnly: &v}
}
