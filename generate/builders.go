// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"sort"
)

// genBuilders generates ergonomic helpers into the builders section:
//   - NewX constructors for every struct (required fields only) and array type
//   - WithF methods for every field of every struct (fluent optional building)
//   - union wrappers: a method M.U() (or package-level func UFromF) for every
//     member M of every union type U, returning U with that case set.
func genBuilders(model *Model) {
	genUnionWrappers()
}

func genUnionWrappers() {
	for _, nt := range genTypes {
		if nt.kind != "or" {
			continue
		}
		if _, collapsed := goplsType[typeNames[nt.typ]]; collapsed {
			continue
		}
		un := goplsName(nt.typ)
		names := []string{}
		for _, t := range nt.items {
			if notNil(t) {
				names = append(names, goplsName(t))
			}
		}
		sort.Strings(names)
		fields := orFields(names)
		for i, field := range fields {
			member := names[i]
			var b bytes.Buffer
			fmt.Fprintf(&b, "// %sFrom%s wraps a %s value as a %s union.\n", un, field, member, un)
			fmt.Fprintf(&b, "func %sFrom%s(v %s) %s {\n\treturn %s{%s: &v}\n}\n\n", un, field, member, un, un, field)
			builders[un+"\x00from\x00"+field] = b.String()
		}
	}
}
