// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lsp

// This file declares URI and DocumentURI.
//
// For the LSP definition of these types, see
// https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#uri

// A DocumentURI is the URI of a client editor document.
//
// According to the LSP specification:
//
//	Care should be taken to handle encoding in URIs. For
//	example, some clients (such as VS Code) may encode colons
//	in drive letters while others do not. The URIs below are
//	both valid, but clients and servers should be consistent
//	with the form they use themselves to ensure the other party
//	doesn’t interpret them as distinct URIs. Clients and
//	servers should not assume that each other are encoding the
//	same way (for example a client encoding colons in drive
//	letters cannot assume server responses will have encoded
//	colons). The same applies to casing of drive letters - one
//	party should not assume the other party will return paths
//	with drive letters cased the same as it.
//
//	file:///c:/project/readme.md
//	file:///C%3A/project/readme.md
type DocumentURI = string

// A URI is an arbitrary URL (e.g. https), not necessarily a file.
type URI = string
