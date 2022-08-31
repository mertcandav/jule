package tokens

// lex.Tokenen kinds.
const (
	DOUBLE_COLON        = "::"
	COLON               = ":"
	SEMICOLON           = ";"
	COMMA               = ","
	TRIPLE_DOT          = "..."
	DOT                 = "."
	PLUS_EQUAL          = "+="
	MINUS_EQUAL         = "-="
	STAR_EQUAL          = "*="
	SLASH_EQUAL         = "/="
	PERCENT_EQUAL       = "%="
	LSHIFT_EQUAL        = "<<="
	RSHIFT_EQUAL        = ">>="
	CARET_EQUAL         = "^="
	AMPER_EQUAL         = "&="
	VLINE_EQUAL         = "|="
	EQUALS              = "=="
	NOT_EQUALS          = "!="
	GREAT_EQUAL         = ">="
	LESS_EQUAL          = "<="
	DOUBLE_AMPER        = "&&"
	DOUBLE_VLINE        = "||"
	LSHIFT              = "<<"
	RSHIFT              = ">>"
	DOUBLE_PLUS         = "++"
	DOUBLE_MINUS        = "--"
	PLUS                = "+"
	MINUS               = "-"
	STAR                = "*"
	SOLIDUS             = "/"
	PERCENT             = "%"
	AMPER               = "&"
	VLINE               = "|"
	CARET               = "^"
	EXCLAMATION         = "!"
	LESS                = "<"
	GREAT               = ">"
	EQUAL               = "="
	LINE_COMMENT        = "//"
	RANGE_COMMENT_OPEN  = "/*"
	RANGE_COMMENT_CLOSE = "*/"
	LPARENTHESES        = "("
	RPARENTHESES        = ")"
	LBRACKET            = "["
	RBRACKET            = "]"
	LBRACE              = "{"
	RBRACE              = "}"
	I8                  = "i8"
	I16                 = "i16"
	I32                 = "i32"
	I64                 = "i64"
	U8                  = "u8"
	U16                 = "u16"
	U32                 = "u32"
	U64                 = "u64"
	F32                 = "f32"
	F64                 = "f64"
	UINT                = "uint"
	INT                 = "int"
	UINTPTR             = "uintptr"
	BOOL                = "bool"
	STR                 = "str"
	ANY                 = "any"
	TRUE                = "true"
	FALSE               = "false"
	NIL                 = "nil"
	CONST               = "const"
	RET                 = "ret"
	TYPE                = "type"
	FOR                 = "for"
	BREAK               = "break"
	CONTINUE            = "continue"
	IN                  = "in"
	IF                  = "if"
	ELSE                = "else"
	USE                 = "use"
	PUB                 = "pub"
	DEFER               = "defer"
	GOTO                = "goto"
	ENUM                = "enum"
	STRUCT              = "struct"
	CO                  = "co"
	MATCH               = "match"
	CASE                = "case"
	DEFAULT             = "default"
	SELF                = "self"
	TRAIT               = "trait"
	IMPL                = "impl"
	CPP                 = "cpp"
	FALLTHROUGH         = "fallthrough"
	FN                  = "fn"
	LET                 = "let"
	UNSAFE              = "unsafe"
	MUT                 = "mut"
)
