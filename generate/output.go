// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"log"
	"regexp"
	"slices"
	"sort"
	"strings"
)

var (
	// tsclient.go has 3 sections
	cdecls = make(sortedMap[string])
	ccases = make(sortedMap[string])
	cfuncs = make(sortedMap[string])
	// tsserver.go has 3 sections
	sdecls = make(sortedMap[string])
	scases = make(sortedMap[string])
	sfuncs = make(sortedMap[string])
	// tsprotocol.go has 2 sections
	types  = make(sortedMap[string])
	consts = make(sortedMap[string])
	// tsjson has 1 section
	jsons = make(sortedMap[string])
	// tsbuilders has 1 section (constructors, WithX methods, union wrappers)
	builders = make(sortedMap[string])
)

func generateOutput(model *Model) {
	for _, r := range model.Requests {
		genDecl(model, r.Method, r.Params, r.Result, r.Direction)
		genCase(model, r.Method, r.Params, r.Result, r.Direction)
		genFunc(model, r.Method, r.Params, r.Result, r.Direction, false)
	}
	for _, n := range model.Notifications {
		if n.Method == "$/cancelRequest" {
			continue // handled internally by jsonrpc2
		}
		genDecl(model, n.Method, n.Params, nil, n.Direction)
		genCase(model, n.Method, n.Params, nil, n.Direction)
		genFunc(model, n.Method, n.Params, nil, n.Direction, true)
	}
	genStructs(model)
	genAliases(model)
	genGenTypes() // generate the unnamed types
	genConsts(model)
	genMarshal()
	genBuilders(model)
}

func genDecl(model *Model, method string, param, result *Type, dir string) {
	fname := methodName(method)
	p := ""
	if notNil(param) {
		p = ", *" + goplsName(param)
	}
	ret := "error"
	if notNil(result) {
		tp := goplsName(result)
		if !hasNilValue(tp) {
			tp = "*" + tp
		}
		ret = fmt.Sprintf("(%s, error)", tp)
	}
	// special gopls compatibility case (PJW: still needed?)
	switch method {
	case "workspace/configuration":
		// was And_Param_workspace_configuration, but the type substitution doesn't work,
		// as ParamConfiguration is embedded in And_Param_workspace_configuration
		p = ", *ParamConfiguration"
		ret = "([]LSPAny, error)"
	}
	fragment := strings.ReplaceAll(strings.TrimPrefix(method, "$/"), "/", "_")
	msg := fmt.Sprintf("\t%s\t%s(context.Context%s) %s\n", lspLink(model, fragment), fname, p, ret)
	switch dir {
	case "clientToServer":
		sdecls[method] = msg
	case "serverToClient":
		cdecls[method] = msg
	case "both":
		sdecls[method] = msg
		cdecls[method] = msg
	default:
		log.Fatalf("impossible direction %q", dir)
	}
}

func genCase(_ *Model, method string, param, result *Type, dir string) {
	out := new(bytes.Buffer)
	fmt.Fprintf(out, "\tcase %q:\n", method)
	var p string
	fname := methodName(method)
	if notNil(param) {
		nm := goplsName(param)
		if method == "workspace/configuration" { // gopls compatibility
			// was And_Param_workspace_configuration, which contains ParamConfiguration
			// so renaming the type leads to circular definitions
			nm = "ParamConfiguration" // gopls compatibility
		}
		fmt.Fprintf(out, "\t\tvar params %s\n", nm)
		fmt.Fprintf(out, "\t\tif err := UnmarshalJSON(req.Params, &params); err != nil {\n")
		fmt.Fprintf(out, "\t\t\treturn nil, err\n\t\t}\n")
		p = ", &params"
	}
	if notNil(result) {
		fmt.Fprintf(out, "\t\tresp, err := %%s.%s(ctx%s)\n", fname, p)
		out.WriteString("\t\tif err != nil {\n")
		out.WriteString("\t\t\treturn nil, err\n")
		out.WriteString("\t\t}\n")
		out.WriteString("\t\treturn resp, nil\n")
	} else {
		fmt.Fprintf(out, "\t\terr := %%s.%s(ctx%s)\n", fname, p)
		out.WriteString("\t\treturn nil, err\n")
	}
	out.WriteString("\n")
	msg := out.String()
	switch dir {
	case "clientToServer":
		scases[method] = fmt.Sprintf(msg, "server")
	case "serverToClient":
		ccases[method] = fmt.Sprintf(msg, "client")
	case "both":
		scases[method] = fmt.Sprintf(msg, "server")
		ccases[method] = fmt.Sprintf(msg, "client")
	default:
		log.Fatalf("impossible direction %q", dir)
	}
}

func genFunc(_ *Model, method string, param, result *Type, dir string, isnotify bool) {
	out := new(bytes.Buffer)
	var p, r string
	var goResult string
	if notNil(param) {
		p = ", params *" + goplsName(param)
	}
	if notNil(result) {
		goResult = goplsName(result)
		if !hasNilValue(goResult) {
			goResult = "*" + goResult
		}
		r = fmt.Sprintf("(%s, error)", goResult)
	} else {
		r = "error"
	}
	// special gopls compatibility case
	switch method {
	case "workspace/configuration":
		// was And_Param_workspace_configuration, but the type substitution doesn't work,
		// as ParamConfiguration is embedded in And_Param_workspace_configuration
		p = ", params *ParamConfiguration"
		r = "([]LSPAny, error)"
		goResult = "[]LSPAny"
	}
	fname := methodName(method)
	fmt.Fprintf(out, "func (s *%%sDispatcher) %s(ctx context.Context%s) %s {\n",
		fname, p, r)

	if !notNil(result) {
		if isnotify {
			if notNil(param) {
				fmt.Fprintf(out, "\treturn s.sender.Notify(ctx, %q, params)\n", method)
			} else {
				fmt.Fprintf(out, "\treturn s.sender.Notify(ctx, %q, nil)\n", method)
			}
		} else {
			if notNil(param) {
				fmt.Fprintf(out, "\treturn s.sender.Call(ctx, %q, params, nil)\n", method)
			} else {
				fmt.Fprintf(out, "\treturn s.sender.Call(ctx, %q, nil, nil)\n", method)
			}
		}
	} else {
		fmt.Fprintf(out, "\tvar result %s\n", goResult)
		if isnotify {
			if notNil(param) {
				fmt.Fprintf(out, "\ts.sender.Notify(ctx, %q, params)\n", method)
			} else {
				fmt.Fprintf(out, "\t\tif err := s.sender.Notify(ctx, %q, nil); err != nil {\n", method)
			}
		} else {
			if notNil(param) {
				fmt.Fprintf(out, "\t\tif err := s.sender.Call(ctx, %q, params, &result); err != nil {\n", method)
			} else {
				fmt.Fprintf(out, "\t\tif err := s.sender.Call(ctx, %q, nil, &result); err != nil {\n", method)
			}
		}
		fmt.Fprintf(out, "\t\treturn nil, err\n\t}\n\treturn result, nil\n")
	}
	out.WriteString("}\n")
	msg := out.String()
	switch dir {
	case "clientToServer":
		sfuncs[method] = fmt.Sprintf(msg, "server")
	case "serverToClient":
		cfuncs[method] = fmt.Sprintf(msg, "client")
	case "both":
		sfuncs[method] = fmt.Sprintf(msg, "server")
		cfuncs[method] = fmt.Sprintf(msg, "client")
	default:
		log.Fatalf("impossible direction %q", dir)
	}
}

func genStructs(model *Model) {
	structures := make(map[string]*Structure) // for expanding Extends
	for _, s := range model.Structures {
		structures[s.Name] = s
	}
	for _, s := range model.Structures {
		out := new(bytes.Buffer)
		generateDoc(out, s.Documentation)
		nm := goName(s.Name)
		fmt.Fprintf(out, "//\n")
		out.WriteString(lspLink(model, camelCase(s.Name)))
		fmt.Fprintf(out, "type %s struct {%s\n", nm, linex(s.Line))
		// for gopls compatibility, embed most extensions, but expand the rest some day
		props := slices.Clone(s.Properties)
		genProps(out, props, nm)
		for _, ex := range s.Extends {
			fmt.Fprintf(out, "\t%s\n", goName(ex.Name))
		}
		for _, ex := range s.Mixins {
			fmt.Fprintf(out, "\t%s\n", goName(ex.Name))
		}
		out.WriteString("}\n")
		declare(nm)
		types[nm] = out.String()
	}

	// base types
	// (For URI and DocumentURI, see ../uri.go.)
	declare("LSPAny")
	types["LSPAny"] = "type LSPAny = any\n"
}

// "FooBar" -> "fooBar"
func camelCase(TitleCased string) string {
	return strings.ToLower(TitleCased[:1]) + TitleCased[1:]
}

func lspLink(model *Model, fragment string) string {
	// Derive URL version from metaData.version in JSON file.
	parts := strings.Split(model.Version.Version, ".") // e.g. "3.17.0"
	return fmt.Sprintf("// See https://microsoft.github.io/language-server-protocol/specifications/lsp/%s.%s/specification#%s\n",
		parts[0], parts[1], // major.minor
		fragment)
}

func genProps(out *bytes.Buffer, props []NameType, name string) {
	for _, p := range props {
		tp, omit, star := propType(name, p)
		json := fmt.Sprintf(" `json:\"%s\"`", p.Name)
		if omit {
			json = fmt.Sprintf(" `json:\"%s,omitempty\"`", p.Name)
		}
		generateDoc(out, p.Documentation)
		if star {
			fmt.Fprintf(out, "\t%s *%s %s\n", goName(p.Name), tp, json)
		} else {
			fmt.Fprintf(out, "\t%s %s %s\n", goName(p.Name), tp, json)
		}
	}
}

// propType returns the Go type, omitempty-ness, and pointer-ness of a struct
// property, applying the renameProp and propStar gopls-compatibility tables.
func propType(structName string, p NameType) (tp string, omit, star bool) {
	tp = goplsName(p.Type)
	if newNm, ok := renameProp[prop{structName, p.Name}]; ok {
		usedRenameProp[prop{structName, p.Name}] = true
		tp = newNm
	}
	omit, star = propStar(structName, p, tp)
	return
}

func genAliases(model *Model) {
	for _, ta := range model.TypeAliases {
		out := new(bytes.Buffer)
		generateDoc(out, ta.Documentation)
		nm := goName(ta.Name)
		if nm != ta.Name {
			continue // renamed the type, e.g., "DocumentDiagnosticReport", an or-type to "string"
		}
		tp := goplsName(ta.Type)
		if tp == nm {
			// the alias target now has the same clean name (e.g. Definition is the
			// clean name of Or_Definition); the union type carries the name directly.
			continue
		}
		fmt.Fprintf(out, "//\n")
		out.WriteString(lspLink(model, camelCase(ta.Name)))
		fmt.Fprintf(out, "type %s = %s // (alias)\n", nm, tp)
		declare(nm)
		aliasNames[nm] = true
		types[nm] = out.String()
	}
}

func genGenTypes() {
	for _, nt := range genTypes {
		out := new(bytes.Buffer)
		nm := goplsName(nt.typ)
		switch nt.kind {
		case "literal":
			fmt.Fprintf(out, "// created for Literal (%s)\n", nt.name)
			fmt.Fprintf(out, "type %s struct {%s\n", nm, linex(nt.line+1))
			genProps(out, nt.properties, nt.name) // systematic name, not gopls name; is this a good choice?
		case "or":
			if _, collapsed := goplsType[typeNames[nt.typ]]; collapsed {
				// It was replaced by a narrower type defined elsewhere
				continue
			}
			names := []string{}
			for _, t := range nt.items {
				if notNil(t) {
					names = append(names, goplsName(t))
				}
			}
			sort.Strings(names)
			fields := orFields(names)
			fmt.Fprintf(out, "// created for Or %v\n", names)
			fmt.Fprintf(out, "type %s struct {%s\n", nm, linex(nt.line+1))
			for i, fld := range fields {
				fmt.Fprintf(out, "\t%s *%s\n", fld, names[i])
			}
		case "and":
			fmt.Fprintf(out, "// created for And\n")
			fmt.Fprintf(out, "type %s struct {%s\n", nm, linex(nt.line+1))
			for _, x := range nt.items {
				nm := goplsName(x)
				fmt.Fprintf(out, "\t%s\n", nm)
			}
		case "tuple": // there's only this one
			nt.name = "UIntCommaUInt"
			fmt.Fprintf(out, "//created for Tuple\ntype %s struct {%s\n", nm, linex(nt.line+1))
			fmt.Fprintf(out, "\tFld0 uint32 `json:\"fld0\"`\n")
			fmt.Fprintf(out, "\tFld1 uint32 `json:\"fld1\"`\n")
		default:
			log.Fatalf("%s not handled", nt.kind)
		}
		out.WriteString("}\n")
		declare(nm)
		types[nm] = out.String()
	}
}
func genConsts(model *Model) {
	for _, e := range model.Enumerations {
		out := new(bytes.Buffer)
		generateDoc(out, e.Documentation)
		tp := goplsName(e.Type)
		nm := goName(e.Name)
		fmt.Fprintf(out, "type %s %s%s\n", nm, tp, linex(e.Line))
		declare(nm)
		types[nm] = out.String()
		vals := new(bytes.Buffer)
		generateDoc(vals, e.Documentation)
		for _, v := range e.Values {
			generateDoc(vals, v.Documentation)
			nm := goName(v.Name)
			more, ok := disambiguate[e.Name]
			if ok {
				usedDisambiguate[e.Name] = true
				nm = more.prefix + nm + more.suffix
				nm = goName(nm) // stringType
			}
			var val string
			switch v := v.Value.(type) {
			case string:
				val = fmt.Sprintf("%q", v)
			case float64:
				val = fmt.Sprintf("%d", int(v))
			default:
				log.Fatalf("impossible type %T", v)
			}
			fmt.Fprintf(vals, "\t%s %s = %s%s\n", nm, e.Name, val, linex(v.Line))
		}
		consts[nm] = vals.String()
	}
}
func genMarshal() {
	for _, nt := range genTypes {
		nm := goplsName(nt.typ)
		if nt.kind != "or" {
			continue
		}
		if _, collapsed := goplsType[typeNames[nt.typ]]; collapsed {
			continue
		}
		names := []string{}
		for _, t := range nt.items {
			if notNil(t) {
				names = append(names, goplsName(t))
			}
		}
		sort.Strings(names)
		fields := orFields(names)
		var buf bytes.Buffer
		fmt.Fprintf(&buf, "func (t %s) MarshalJSON() ([]byte, error) {\n", nm)
		buf.WriteString("\tswitch {\n")
		for _, fld := range fields {
			fmt.Fprintf(&buf, "\tcase t.%s != nil:\n", fld)
			fmt.Fprintf(&buf, "\t\treturn json.Marshal(*t.%s)\n", fld)
		}
		buf.WriteString("\t}\n")
		buf.WriteString("\treturn []byte(\"null\"), nil\n")
		buf.WriteString("}\n\n")

		fmt.Fprintf(&buf, "func (t *%s) UnmarshalJSON(x []byte) error {\n", nm)
		fmt.Fprintf(&buf, "\t*t = %s{}\n", nm)
		buf.WriteString("\tif string(x) == \"null\" {\n\t\treturn nil\n\t}\n")
		for i, fld := range fields {
			fmt.Fprintf(&buf, "\tvar h%d %s\n", i, names[i])
			fmt.Fprintf(&buf, "\tif err := json.Unmarshal(x, &h%d); err == nil {\n\t\tt.%s = &h%d\n\t\t\treturn nil\n\t\t}\n", i, fld, i)
		}
		fmt.Fprintf(&buf, "return &UnmarshalError{\"unmarshal failed to match one of %v\"}", names)
		buf.WriteString("}\n\n")
		jsons[nm] = buf.String()
	}
}

func linex(n int) string {
	if *lineNumbers {
		return fmt.Sprintf(" // line %d", n)
	}
	return ""
}

func goplsName(t *Type) string {
	nm := typeNames[t]
	// translate systematic name to gopls name
	if newNm, ok := goplsType[nm]; ok {
		usedGoplsType[nm] = true
		nm = newNm
	}
	// arrays: clean the element name, and turn arrays of exported named types
	// into defined types, e.g. []Location -> Locations.
	if elem, ok := strings.CutPrefix(nm, "[]"); ok {
		elem = cleanOrPrefix(elem)
		return "[]" + elem
	}
	// surviving "or" types get a clean, non-generated-sounding name
	return cleanOrPrefix(nm)
}

// cleanOrPrefix applies cleanOrName to systematic "Or_" names, leaving others alone.
func cleanOrPrefix(nm string) string {
	if strings.HasPrefix(nm, "Or_") {
		return cleanOrName(nm)
	}
	return nm
}

// cleanOrName turns a systematic "or" name like Or_WorkspaceSymbol_location into
// a clean exported name like WorkspaceSymbolLocation.
func cleanOrName(systematic string) string {
	rest := strings.TrimPrefix(systematic, "Or_")
	var b strings.Builder
	for _, seg := range strings.Split(rest, "_") {
		b.WriteString(title(seg))
	}
	return b.String()
}

// aliasNames records exported names that are emitted as type aliases (`type X = Y`).
// Methods cannot be attached to aliases of non-local types.
var aliasNames = map[string]bool{"LSPAny": true, "WatchKind": true}

// declared records every top-level type name emitted, to detect conflicts.
var declared = map[string]bool{}

func declare(nm string) {
	if declared[nm] {
		log.Fatalf("name conflict: %q is declared more than once", nm)
	}
	declared[nm] = true
}

var identRe = regexp.MustCompile(`[A-Za-z_][A-Za-z0-9_]*`)

// pruneUnusedTypes removes generated synthetic types (Or/And/Literal/Tuple) that
// nothing references. They arise when the field that would use them is overridden
// to a different type via the renameProp/goplsType tables, leaving the generated
// union orphaned. Real spec types (structs, enums, aliases) are always kept.
func pruneUnusedTypes(handwritten string) {
	// candidate set: every emitted synthetic type name.
	synthetic := map[string]bool{}
	for _, nt := range genTypes {
		if nt.kind == "or" {
			if _, collapsed := goplsType[typeNames[nt.typ]]; collapsed {
				continue
			}
		}
		synthetic[goplsName(nt.typ)] = true
	}

	// syntheticRefs returns the synthetic names mentioned in text.
	syntheticRefs := func(text string) []string {
		var out []string
		for _, w := range identRe.FindAllString(text, -1) {
			if synthetic[w] {
				out = append(out, w)
			}
		}
		return out
	}

	// text "owned" by each synthetic type (its own decl, marshal and converters),
	// which only contributes references if that type is itself reachable.
	owned := map[string]*strings.Builder{}
	for s := range synthetic {
		var b strings.Builder
		b.WriteString(types[s])
		b.WriteString(jsons[s])
		owned[s] = &b
	}
	for k, v := range builders {
		if i := strings.IndexByte(k, 0); i > 0 {
			if b, ok := owned[k[:i]]; ok {
				b.WriteString(v)
			}
		}
	}

	// root text: everything that is not owned by a synthetic type — non-synthetic
	// declarations, constants, the client/server interfaces, and handwritten code.
	var root strings.Builder
	for k, v := range types {
		if !synthetic[k] {
			root.WriteString(v)
		}
	}
	for _, v := range consts {
		root.WriteString(v)
	}
	for _, m := range []sortedMap[string]{cdecls, ccases, cfuncs, sdecls, scases, sfuncs} {
		for _, v := range m {
			root.WriteString(v)
		}
	}
	root.WriteString(handwritten)

	// breadth-first reachability from the roots through synthetic references.
	reached := map[string]bool{}
	var queue []string
	for _, s := range syntheticRefs(root.String()) {
		if !reached[s] {
			reached[s] = true
			queue = append(queue, s)
		}
	}
	for len(queue) > 0 {
		s := queue[0]
		queue = queue[1:]
		for _, r := range syntheticRefs(owned[s].String()) {
			if !reached[r] {
				reached[r] = true
				queue = append(queue, r)
			}
		}
	}

	// drop unreachable synthetics from every output section.
	for s := range synthetic {
		if reached[s] {
			continue
		}
		log.Printf("pruning unused type %s", s)
		delete(types, s)
		delete(jsons, s)
		for k := range builders {
			if i := strings.IndexByte(k, 0); i > 0 && k[:i] == s {
				delete(builders, k)
			}
		}
	}
}

// orFields maps the gopls type names of an Or type's members to valid, unique,
// exported Go field names, aligned by index with names.
func orFields(names []string) []string {
	fields := make([]string, len(names))
	seen := map[string]int{}
	for i, n := range names {
		f := fieldName(n)
		if k := seen[f]; k > 0 {
			f = fmt.Sprintf("%s%d", f, k)
		}
		seen[fieldName(n)]++
		fields[i] = f
	}
	return fields
}

// fieldName turns a single gopls type name into an exported Go identifier.
// e.g. "Location" -> "Location", "[]Location" -> "Locations",
// "string" -> "String", "*Foo" -> "Foo".
func fieldName(n string) string {
	n = strings.TrimPrefix(n, "*")
	if elem := strings.TrimPrefix(n, "[]"); elem != n {
		return fieldName(elem) + "s"
	}
	return title(n)
}

func notNil(t *Type) bool { // shutdwon is the special case that needs this
	return t != nil && (t.Kind != "base" || t.Name != "null")
}

func hasNilValue(t string) bool {
	// this may be unreliable, and need a supplementary table
	if strings.HasPrefix(t, "[]") || strings.HasPrefix(t, "*") {
		return true
	}
	if t == "interface{}" || t == "any" {
		return true
	}
	// that's all the cases that occur currently
	return false
}
