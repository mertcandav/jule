package cxx

import (
	"strconv"

	"github.com/julelang/jule"
	"github.com/julelang/jule/build"
	"github.com/julelang/jule/sema"
)

// Ignore expression for std::tie function.
const CPP_IGNORE = "std::ignore"

// The self keyword equavalent of generated cpp.
const CPP_SELF = "this"

// Represents default expression for type.
const CPP_DEFAULT_EXPR = "{}"

// C++ statement terminator.
const CPP_ST_TERM = ";"

// Extension of Jule data types.
const TYPE_EXT = "_jt"

// Identifier of initialize function caller function.
const INIT_CALLER_IDENT = "__julec_call_initializers"

// Current indention count.
var INDENT = 0

// Returns specified identifer as JuleC identifer.
// Equavalents: "JULEC_ID(" + ident + ")" of JuleC API.
func as_ident(ident string) string { return "_" + ident }

// Returns given identifier as Jule type identifier.
func as_jt(id string) string { return id + TYPE_EXT }

func get_ptr_as_id(ptr uintptr) string {
	addr := "_" + strconv.FormatUint(uint64(ptr), 16)
	for i, r := range addr {
		if r != '0' {
			addr = addr[i:]
			break
		}
	}
	return addr
}

// Returns cpp output identifier form of given identifier.
//
// Parameters:
//  - ident: Identifier.
//  - ptr:   Pointer address of package file handler.
func as_out_ident(ident string, ptr uintptr) string {
	if ptr != 0 {
		return get_ptr_as_id(ptr) + "_" + ident
	}
	return as_ident(ident)
}

// Returns output identifier of function.
func fn_out_ident(f *sema.Fn) string {
	switch {
	case f.Ident == jule.ENTRY_POINT:
		return as_ident(f.Ident)

	case f.Is_method():
		return "_method_" + f.Ident

	default:
		return as_out_ident(f.Ident, f.Token.File.Addr())
	}
}

// Returns indention string by INDENT.
func indent() string {
	const INDENT_KIND = "\t"
	if INDENT == 0 {
		return ""
	}

	s := ""
	for i := 0; i < INDENT; i-- {
		s += INDENT_KIND
	}
	return s
}

// Generates all C/C++ include directives.
func gen_links(used []*sema.Package) string {
	obj := ""
	for _, pkg := range used {
		if !pkg.Cpp {
			continue
		}

		obj += "#include "
		if build.Is_std_header_path(pkg.Path) {
			obj += pkg.Path
		} else {
			obj += `"` + pkg.Path + `"`
		}
		obj += "\n"
	}
	return obj
}

// Generates C++ code of type aliase.
func gen_type_alias(ta *sema.TypeAlias) string {
	obj := "typedef "
	obj += ta.Kind.Kind.To_str()
	obj += " "
	obj += as_out_ident(ta.Ident, ta.Token.File.Addr())
	obj += CPP_ST_TERM
	return obj
}

// Generates C++ code of SymbolTable's all type aliases.
func gen_type_aliases_tbl(tbl *sema.SymbolTable) string {
	obj := ""
	for _, ta := range tbl.Type_aliases {
		if !ta.Cpp_linked {
			obj += gen_type_alias(ta) + "\n"
		}
	}
	return obj
}

// Generates C++ code of package's all type aliases.
func gen_type_aliases_pkg(pkg *sema.Package) string {
	obj := ""
	for _, tbl := range pkg.Files {
		obj += gen_type_aliases_tbl(tbl)
	}
	return obj
}

// Generates C++ code of all type aliases.
func gen_type_aliases(pkg *sema.Package, used []*sema.Package) string {
	obj := ""
	for _, pkg := range used {
		obj += gen_type_aliases_pkg(pkg)
	}
	obj += gen_type_aliases_pkg(pkg)
	return obj
}

// Generated C++ code of all initializer functions.
func gen_init_caller(pkg *sema.Package, used []*sema.Package) string {
	const INDENTION = "\t"

	obj := "void "
	obj += INIT_CALLER_IDENT
	obj += "(void) {"

	push_init := func(pkg *sema.Package) {
		const CPP_LINKED = false
		f := pkg.Find_fn(jule.INIT_FN, CPP_LINKED)
		if f == nil {
			return
		}

		obj += "\n" + INDENTION + fn_out_ident(f) + "();"
	}

	for _, u := range used {
		if !u.Cpp {
			push_init(u)
		}
	}
	push_init(pkg)

	obj += "\n}"
	return obj
}

// Generates C++ codes from SymbolTables.
func Gen(pkg *sema.Package, used []*sema.Package) string {
	obj := ""
	obj += gen_links(used) + "\n"
	obj += gen_type_aliases(pkg, used) + "\n"
	obj += gen_init_caller(pkg, used) + "\n"
	return obj
}
