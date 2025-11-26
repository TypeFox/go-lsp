// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lsp_test

import (
	"encoding/json"
	"testing"

	"typefox.dev/lsp"
)

// TestCancellationSupport tests the LSP cancellation protocol
func TestCancellationSupport(t *testing.T) {

	t.Run("CancelParams", func(t *testing.T) {
		// Test CancelParams structure
		cancelParams := lsp.CancelParams{
			ID: "test-request-123",
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

		if cancelParams.ID != unmarshaled.ID {
			t.Errorf("Expected ID %v, got %v", cancelParams.ID, unmarshaled.ID)
		}
	})

	t.Run("CancelNotificationStructure", func(t *testing.T) {
		// Test that cancel notification structure is correct
		cancelParams := lsp.CancelParams{
			ID: "request-to-cancel",
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

		if cancelParams.ID != unmarshaled.ID {
			t.Errorf("Expected ID %v, got %v", cancelParams.ID, unmarshaled.ID)
		}
	})
}

// TestProgressSupport tests the LSP progress protocol
func TestProgressSupport(t *testing.T) {

	t.Run("ProgressParams", func(t *testing.T) {
		// Test ProgressParams structure
		progressParams := lsp.ProgressParams{
			Token: "progress-token-123",
			Value: map[string]interface{}{
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

		if progressParams.Token != unmarshaled.Token {
			t.Errorf("Expected Token %v, got %v", progressParams.Token, unmarshaled.Token)
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

		if begin.Title != unmarshaled.Title {
			t.Errorf("Expected Title %v, got %v", begin.Title, unmarshaled.Title)
		}
		if begin.Cancellable != unmarshaled.Cancellable {
			t.Errorf("Expected Cancellable %v, got %v", begin.Cancellable, unmarshaled.Cancellable)
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

		if *report.Percentage != *unmarshaled.Percentage {
			t.Errorf("Expected Percentage %v, got %v", *report.Percentage, *unmarshaled.Percentage)
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
