package types

import (
	"github.com/julelang/jule/ast"
	"github.com/julelang/jule/build"
	"github.com/julelang/jule/lex"
)

// I8CompatibleWith reports i8 is compatible or not with data-type specified type.
func I8CompatibleWith(t uint8) bool {
	t = GetRealCode(t)
	return t == I8
}

// I16CompatibleWith reports i16 is compatible or not with data-type specified type.
func I16CompatibleWith(t uint8) bool {
	t = GetRealCode(t)
	return t == I8 || t == I16 || t == U8
}

// I32CompatibleWith reports i32 is compatible or not with data-type specified type.
func I32CompatibleWith(t uint8) bool {
	t = GetRealCode(t)
	return t == I8 || t == I16 || t == I32 || t == U8 || t == U16
}

// I64CompatibleWith reports i64 is compatible or not with data-type specified type.
func I64CompatibleWith(t uint8) bool {
	t = GetRealCode(t)
	switch t {
	case I8, I16, I32, I64, U8, U16, U32:
		return true
	default:
		return false
	}
}

// U8CompatibleWith reports u8 is compatible or not with data-type specified type.
func U8CompatibleWith(t uint8) bool {
	t = GetRealCode(t)
	return t == U8
}

// U16CompatibleWith reports u16 is compatible or not with data-type specified type.
func U16CompatibleWith(t uint8) bool {
	t = GetRealCode(t)
	return t == U8 || t == U16
}

// U32CompatibleWith reports u32 is compatible or not with data-type specified type.
func U32CompatibleWith(t uint8) bool {
	t = GetRealCode(t)
	return t == U8 || t == U16 || t == U32
}

// U16CompatibleWith reports u64 is compatible or not with data-type specified type.
func U64CompatibleWith(t uint8) bool {
	t = GetRealCode(t)
	return t == U8 || t == U16 || t == U32 || t == U64
}

// F32CompatibleWith reports f32 is compatible or not with data-type specified type.
func F32CompatibleWith(t uint8) bool {
	t = GetRealCode(t)
	switch t {
	case F32, I8, I16, I32, I64, U8, U16, U32, U64:
		return true
	default:
		return false
	}
}

// F64CompatibleWith reports f64 is compatible or not with data-type specified type.
func F64CompatibleWith(t uint8) bool {
	t = GetRealCode(t)
	switch t {
	case F64, F32, I8, I16, I32, I64, U8, U16, U32, U64:
		return true
	default:
		return false
	}
}

// TypeAreCompatible reports type one and type two is compatible or not.
func TypesAreCompatible(t1, t2 uint8, ignoreany bool) bool {
	t1 = GetRealCode(t1)
	switch t1 {
	case ANY:
		return !ignoreany
	case I8:
		return I8CompatibleWith(t2)
	case I16:
		return I16CompatibleWith(t2)
	case I32:
		return I32CompatibleWith(t2)
	case I64:
		return I64CompatibleWith(t2)
	case U8:
		return U8CompatibleWith(t2)
	case U16:
		return U16CompatibleWith(t2)
	case U32:
		return U32CompatibleWith(t2)
	case U64:
		return U64CompatibleWith(t2)
	case BOOL:
		return t2 == BOOL
	case STR:
		return t2 == STR
	case F32:
		return F32CompatibleWith(t2)
	case F64:
		return F64CompatibleWith(t2)
	case NIL:
		return t2 == NIL
	}
	return false
}

// Checker is type checker.
type Checker struct {
	ErrTok      lex.Token
	L           Type
	R           Type
	ErrorLogged bool
	IgnoreAny   bool
	AllowAssign bool
	Errors      []build.Log
}

// pusherrtok appends new error by token.
func (c *Checker) pusherrtok(tok lex.Token, key string, args ...any) {
	c.Errors = append(c.Errors, build.Log{
		Type:    build.ERR,
		Row:     tok.Row,
		Column:  tok.Column,
		Path:    tok.File.Path(),
		Message: build.Errorf(key, args...),
	})
}

func (c *Checker) check_ref() bool {
	if c.L.Kind == c.R.Kind {
		return true
	} else if !c.AllowAssign {
		return false
	}
	c.L = DerefPtrOrRef(c.L)
	return c.Check()
}

func (c *Checker) check_ptr() bool {
	if c.R.Id == NIL {
		return true
	} else if IsUnsafePtr(c.L) {
		return true
	}
	return c.L.Kind == c.R.Kind
}

func trait_has_reference_receiver(t *ast.Trait) bool {
	for _, f := range t.Defines.Fns {
		if IsRef(f.Receiver.Type) {
			return true
		}
	}
	return false
}

func (c *Checker) check_trait() bool {
	if c.R.Id == NIL {
		return true
	}
	t := c.L.Tag.(*ast.Trait)
	lm := c.L.Modifiers()
	ref := false
	switch {
	case IsRef(c.R):
		ref = true
		c.R = DerefPtrOrRef(c.R)
		if !IsStruct(c.R) {
			break
		}
		fallthrough
	case IsStruct(c.R):
		if lm != "" {
			return false
		}
		rm := c.R.Modifiers()
		if rm != "" {
			return false
		}
		s := c.R.Tag.(*ast.Struct)
		if !s.HasTrait(t) {
			return false
		}
		if trait_has_reference_receiver(t) && !ref {
			c.ErrorLogged = true
			c.pusherrtok(c.ErrTok, "trait_has_reference_parametered_function")
			return false
		}
		return true
	case IsTrait(c.R):
		return t == c.R.Tag.(*ast.Trait) && lm == c.R.Modifiers()
	}
	return false
}

func (c *Checker) check_struct() bool {
	if c.R.Tag == nil {
		return false
	}
	s1, s2 := c.L.Tag.(*ast.Struct), c.R.Tag.(*ast.Struct)
	switch {
	case s1.Id != s2.Id,
		s1.Token.File != s2.Token.File:
		return false
	}
	if len(s1.Generics) == 0 {
		return true
	}
	n1, n2 := len(s1.GetGenerics()), len(s2.GetGenerics())
	if n1 != n2 {
		return false
	}
	for i, g1 := range s1.GetGenerics() {
		g2 := s2.GetGenerics()[i]
		if !Equals(g1, g2) {
			return false
		}
	}
	return true
}

func (c *Checker) check_slice() bool {
	if c.R.Id == NIL {
		return true
	}
	return c.L.Kind == c.R.Kind
}

func (c *Checker) check_array() bool {
	if !IsArray(c.R) {
		return false
	}
	return c.L.Size.N == c.R.Size.N
}

func (c *Checker) check_map() bool {
	if c.R.Id == NIL {
		return true
	}
	return c.L.Kind == c.R.Kind
}

// Check checks type compatilility and reports.
func (c *Checker) Check() bool {
	switch {
	case IsTrait(c.L), IsTrait(c.R):
		if IsTrait(c.R) {
			c.L, c.R = c.R, c.L
		}
		return c.check_trait()
	case IsRef(c.L), IsRef(c.R):
		if IsRef(c.R) {
			c.L, c.R = c.R, c.L
		}
		return c.check_ref()
	case IsPtr(c.L), IsPtr(c.R):
		if !IsPtr(c.L) && IsPtr(c.R) {
			c.L, c.R = c.R, c.L
		}
		return c.check_ptr()
	case IsSlice(c.L), IsSlice(c.R):
		if IsSlice(c.R) {
			c.L, c.R = c.R, c.L
		}
		return c.check_slice()
	case IsArray(c.L), IsArray(c.R):
		if IsArray(c.R) {
			c.L, c.R = c.R, c.L
		}
		return c.check_array()
	case IsMap(c.L), IsMap(c.R):
		if IsMap(c.R) {
			c.L, c.R = c.R, c.L
		}
		return c.check_map()
	case IsNilCompatible(c.L):
		return c.R.Id == NIL
	case IsNilCompatible(c.R):
		return c.L.Id == NIL
	case IsEnum(c.L), IsEnum(c.R):
		return c.L.Id == c.R.Id && c.L.Kind == c.R.Kind
	case IsStruct(c.L), IsStruct(c.R):
		if c.R.Id == STRUCT {
			c.L, c.R = c.R, c.L
		}
		return c.check_struct()
	}
	return TypesAreCompatible(c.L.Id, c.R.Id, c.IgnoreAny)
}
