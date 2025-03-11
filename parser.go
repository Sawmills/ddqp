package ddqp

import "github.com/alecthomas/participle/v2/lexer"

// This is the primary lexer for all parsers
// nolint:govet
var lex = lexer.MustSimple([]lexer.SimpleRule{
	{"Comment", `(?i)rem[^\n]*`},
	{"String", `"(\\"|[^"])*"`},
	{"SpaceAggregatorCondition", `v: v[<>=]*([0-9]*[.])?[0-9]+`},
	{"Operator", `[<>=*\/+-]+`},
	{"OperatorValue", `[<>=]+([0-9]*[.])?[0-9]+`},
	{"RangeValue", `\[[0-9]+ TO [0-9]+\]`},
	{"Ident", `[a-zA-Z0-9_\*][\w\d-\*\./\?]*`},
	{"Float", `[+-]?([0-9]*[.])?[0-9]+`},
	{"Int", `\d+`},
	{"Attribute", `@[a-zA-Z0-9_\*][\w\d-\*\./\?]*`},
	{"BooleanOp", `(?i)(AND|OR|NOT)`},
	{"CIDR", `\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}/\d{1,2}`},
	{"RangeOp", `(?i)TO`},
	{"Punct", `[-[!@#$%^&*()+_={}\|:;"'<,>.?\/]|]`},
	{"EOL", `[\n\r]+`},
	{"whitespace", `[ \t]+`},
})

// This is the primary lexer for all parsers
// nolint:govet
var logQueryLex = lexer.MustSimple([]lexer.SimpleRule{
	{"String", `"(\\"|[^"])*"`},
	{"BoolLikeIdent", `(?i)(not|and|or)[\-\w\d\*\./\?]+`},
	{"EscapedBoolean", `\\(?i)(AND|OR|NOT)\b`},
	{"EscapedIdent", `[/a-zA-Z0-9_\*][\w\d\-\*\./\?]*(\\[ +\-=&|><!()\[\]{}^"~*?:\\#])[\w\d\-\*\./\?]*(\\[ +\-=&|><!()\[\]{}^"~*?:\\#][\w\d\-\*\./\?]*)*`},
	{"CIDR", `\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}/\d{1,2}`},
	{"IPAddress", `\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}`},        // IP address pattern
	{"NegativeFloat", `-([0-9]*[.])?[0-9]+([eE][+-]?[0-9]+)?`}, // Add specific pattern for negative floating point numbers
	{"NegativeInt", `-\d+`},                                    // Add specific pattern for negative integers
	{"RangeValue", `\[[0-9]+ TO [0-9]+\]`},
	{"OperatorValue", `[<>]=?[+-]?([0-9]*[.])?[0-9]+([eE][+-]?[0-9]+)?`},
	{"SpaceAggregatorCondition", `v: v[<>=]*([0-9]*[.])?[0-9]+`},
	{"NegatedIdent", `-[a-zA-Z0-9_\*][\w\d-\*\./\?]*`},
	{"Float", `[+]?([0-9]*[.])?[0-9]+([eE][+-]?[0-9]+)?`}, // Changed to only match positive numbers with optional +
	{"Int", `\d+`},
	{"Attribute", `@[a-zA-Z0-9_\*][\w\d-\*\./\?]*`},
	{"ComplexWildcardPath", `\*/[\w\d\-\*\./\?]+(/\*)+(/[\w\d\-\*\./\?]+)*`},
	{"BooleanOp", `(?i)\b(AND|OR|NOT)\b`},
	{"RangeOp", `(?i)TO`},
	{"Operator", `[<>=*\/+-]+`},
	{"Ident", `[a-zA-Z0-9_\*][\w\d\-\*\./\?]*`},
	{"Punct", `[-[!@#$%^&*()+_={}\|:;"'<,>.?\/]|]`},
	{"EOL", `[\n\r]+`},
	{"whitespace", `[ \t]+`},
})
