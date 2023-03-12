// Copyright 2023 The Jule Programming Language.
// Use of this source code is governed by a BSD 3-Clause
// license that can be found in the LICENSE file.

package sema

import (
	"github.com/julelang/jule/ast"
	"github.com/julelang/jule/lex"
)

// Type alias.
type TypeAlias struct {
	Public     bool
	Cpp_linked bool
	Token      lex.Token
	Ident      string
	Kind       *Type
	Doc        string
	Refers     []*ast.IdentType
}

// Type's kind's type.
type TypeKind struct { kind any }
// Returns reference type if kind is reference, nil if not.
func (tk *TypeKind) Ref() *Ref {
	switch tk.kind.(type) {
	case *Ref:
		return tk.kind.(*Ref)

	default:
		return nil
	}
}
// Returns pointer type if kind is pointer, nil if not.
func (tk *TypeKind) Ptr() *Ptr {
	switch tk.kind.(type) {
	case *Ptr:
		return tk.kind.(*Ptr)

	default:
		return nil
	}
}

// Type.
type Type struct {
	Decl *ast.TypeDecl // Never changed by semantic analyzer.
	Kind *TypeKind
}

// Reports whether type is checked already.
func (t *Type) checked() bool { return t.Kind != nil }
// Removes kind and ready to check.
// checked() reports false after this function.
func (t *Type) remove_kind() { t.Kind = nil }

// Primitive type.
type Prim struct { kind string }
// Reports whether type is primitive i8.
func (p *Prim) Is_i8() bool { return p.kind == lex.KND_I8 }
// Reports whether type is primitive i16.
func (p *Prim) Is_i16() bool { return p.kind == lex.KND_I16 }
// Reports whether type is primitive i32.
func (p *Prim) Is_i32() bool { return p.kind == lex.KND_I32 }
// Reports whether type is primitive i64.
func (p *Prim) Is_i64() bool { return p.kind == lex.KND_I64 }
// Reports whether type is primitive u8.
func (p *Prim) Is_u8() bool { return p.kind == lex.KND_U8 }
// Reports whether type is primitive u16.
func (p *Prim) Is_u16() bool { return p.kind == lex.KND_U16 }
// Reports whether type is primitive u32.
func (p *Prim) Is_u32() bool { return p.kind == lex.KND_U32 }
// Reports whether type is primitive u64.
func (p *Prim) Is_u64() bool { return p.kind == lex.KND_U64 }
// Reports whether type is primitive f32.
func (p *Prim) Is_f32() bool { return p.kind == lex.KND_F32 }
// Reports whether type is primitive f64.
func (p *Prim) Is_f64() bool { return p.kind == lex.KND_F64 }
// Reports whether type is primitive int.
func (p *Prim) Is_int() bool { return p.kind == lex.KND_INT }
// Reports whether type is primitive uint.
func (p *Prim) Is_uint() bool { return p.kind == lex.KND_UINT }
// Reports whether type is primitive uintptr.
func (p *Prim) Is_uintptr() bool { return p.kind == lex.KND_UINTPTR }
// Reports whether type is primitive bool.
func (p *Prim) Is_bool() bool { return p.kind == lex.KND_BOOL }
// Reports whether type is primitive str.
func (p *Prim) Is_str() bool { return p.kind == lex.KND_STR }
// Reports whether type is primitive any.
func (p *Prim) Is_any() bool { return p.kind == lex.KND_ANY }

// Reference type.
type Ref struct { Elem *TypeKind }
// Pointer type.
type Ptr struct { Elem *TypeKind }
// Slice type.
type Slc struct { Elem *TypeKind }
// Tuple type.
type Tuple struct { Types []*TypeKind }
// Map type.
type Map struct {
	Key *TypeKind
	Val *TypeKind
}

// Checks type and builds result as kind.
// Removes kind if error occurs,
// so type is not reports true for checked state.
type _TypeChecker struct {
	// Uses Sema for:
	//  - Push errors.
	s *_Sema

	// Uses Lookup for:
	//  - Lookup symbol tables.
	lookup _Lookup

	error_token lex.Token
}

func (tc *_TypeChecker) push_err(token lex.Token, key string, args ...any) {
	tc.s.push_err(token, key, args...)
}

func (tc *_TypeChecker) build_prim(decl *ast.IdentType) *Prim {
	switch decl.Ident {
	case lex.KND_I8,
		lex.KND_I16,
		lex.KND_I32,
		lex.KND_I64,
		lex.KND_U8,
		lex.KND_U16,
		lex.KND_U32,
		lex.KND_U64,
		lex.KND_F32,
		lex.KND_F64,
		lex.KND_INT,
		lex.KND_UINT,
		lex.KND_UINTPTR,
		lex.KND_BOOL,
		lex.KND_STR,
		lex.KND_ANY:
		// Ignore.
	default:
		tc.push_err(tc.error_token, "invalid_type")
		return nil 
	}

	if len(decl.Generics) > 0 {
		tc.push_err(decl.Token, "type_not_supports_generics", decl.Ident)
	}

	return &Prim{
		kind: decl.Ident,
	}
}

func (tc *_TypeChecker) build_ident_kind(decl *ast.IdentType) any {
	// TODO:
	//  - Implement traits.
	//  - Implement enums.
	//  - Implement functions.
	//  - Implement structures.
	//  - Implement type aliases.
	//  - Implement cpp-linked defines.
	switch {
	case decl.Is_prim():
		return tc.build_prim(decl)
	
	default:
		tc.push_err(tc.error_token, "invalid_type")
		return nil
	}
}

func (tc *_TypeChecker) build_ref(decl *ast.RefType) *Ref {
	elem := tc.check_decl(decl.Elem)

	// TODO: check cases:
	//         - ref_refs_array
	//         - ref_refs_enum
	// Check special cases.
	switch {
	case elem == nil:
		return nil

	case elem.Ref() != nil:
		tc.push_err(tc.error_token, "ref_refs_ref")
		return nil

	case elem.Ptr() != nil:
		tc.push_err(tc.error_token, "ref_refs_ptr")
		return nil
	}

	return &Ref{
		Elem: elem,
	}
}

func (tc *_TypeChecker) build_ptr(decl *ast.PtrType) *Ptr {
	elem := tc.check_decl(decl.Elem)

	// Check special cases.
	switch {
	case elem == nil:
		return nil

	case elem.Ref() != nil:
		tc.push_err(tc.error_token, "ptr_points_ref")
		return nil
	}

	return &Ptr{
		Elem: elem,
	}
}

func (tc *_TypeChecker) build_slice(decl *ast.SliceType) *Slc {
	elem := tc.check_decl(decl.Elem)

	// TODO: check cases:
	//         - Elem is Array with auto-sized is invalid element
	// Check special cases.
	switch {
	case elem == nil:
		return nil
	}

	return &Slc{
		Elem: elem,
	}
}

func (tc *_TypeChecker) build_map(decl *ast.MapType) *Map {
	key := tc.check_decl(decl.Key)
	if key == nil {
		return nil
	}

	val := tc.check_decl(decl.Key)
	if key == nil {
		return nil
	}

	return &Map{
		Key: key,
		Val: val,
	}
}

func (tc *_TypeChecker) build_tuple(decl *ast.TupleType) *Tuple {
	types := make([]*TypeKind, len(decl.Types))
	for i, t := range decl.Types {
		kind := tc.check_decl(t)
		if kind == nil {
			return nil
		}
		types[i] = kind
	}

	return &Tuple{
		Types: types,
	}
}

func (tc *_TypeChecker) build_kind(decl_kind ast.TypeDeclKind) *TypeKind {
	var kind any = nil

	// TODO:
	//  - Implement arrays.
	//  - Implement functions.
	switch decl_kind.(type) {
	case *ast.IdentType:
		kind = tc.build_ident_kind(decl_kind.(*ast.IdentType))

	case *ast.RefType:
		kind = tc.build_ref(decl_kind.(*ast.RefType))

	case *ast.PtrType:
		kind = tc.build_ptr(decl_kind.(*ast.PtrType))

	case *ast.SliceType:
		kind = tc.build_slice(decl_kind.(*ast.SliceType))

	case *ast.MapType:
		kind = tc.build_map(decl_kind.(*ast.MapType))

	case *ast.TupleType:
		kind = tc.build_tuple(decl_kind.(*ast.TupleType))

	default:
		tc.push_err(tc.error_token, "invalid_type")
		return nil
	}

	if kind == nil {
		return nil
	}
	return &TypeKind{
		kind: kind,
	}
}

func (tc *_TypeChecker) check_decl(decl *ast.TypeDecl) *TypeKind {
	// Save current token.
	error_token := tc.error_token

	tc.error_token = decl.Token
	kind := tc.build_kind(decl.Kind)
	tc.error_token = error_token
	return kind
}

func (tc *_TypeChecker) check(t *Type) {
	// TODO: Detect cycles.
	// TODO: Check type validity.
	kind := tc.check_decl(t.Decl)
	if kind == nil {
		t.remove_kind()
		return
	}
	t.Kind = kind
}
