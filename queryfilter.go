package ddqp

import (
	"fmt"
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

// QueryFilter represents a Datadog query filter
type QueryFilter struct {
	Pos lexer.Position

	Left  *FilterParam   `@@`
	Right []*FilterParam `@@*`
}

// NegatedFilter represents a negated attribute or tag filter
type NegatedFilter struct {
	Key   string            `@Ident`
	Value *QueryFilterValue `":" @@`
}

// String returns the string representation of a negated filter
func (nf *NegatedFilter) String() string {
	if nf == nil {
		return ""
	}
	return fmt.Sprintf("-%s:%s", nf.Key, nf.Value.String())
}

// NegatedTerm represents a simple negated term (no colon)
type NegatedTerm struct {
	Term string `@NegatedIdent`
}

// String returns the string representation of a negated term
func (nt *NegatedTerm) String() string {
	if nt == nil {
		return ""
	}
	return nt.Term
}

// NegatedAttributeFilter represents a negated attribute filter
type NegatedAttributeFilter struct {
	NegatedIdent string            `( @NegatedIdent` // First try to match as NegatedIdent
	Value        *QueryFilterValue `":" @@ )`        // Followed by colon and value
	OrNegated    bool              `| ( @"-"`        // Or try the traditional approach
	Name         string            `@Attribute`      // with separate tokens
	OrValue      *QueryFilterValue `":" @@ )`
}

// String returns the string representation of a negated attribute filter
func (naf *NegatedAttributeFilter) String() string {
	if naf == nil {
		return ""
	}

	if naf.NegatedIdent != "" {
		return naf.NegatedIdent + ":" + naf.Value.String()
	}

	return "-" + naf.Name + ":" + naf.OrValue.String()
}

// NegatedTagFilter represents a negated tag filter
type NegatedTagFilter struct {
	NegatedIdent string            `( @NegatedIdent` // First try to match as NegatedIdent
	Value        *QueryFilterValue `":" @@ )`        // Followed by colon and value
	OrNegated    bool              `| ( @"-"`        // Or try the traditional approach
	Key          string            `@Ident`          // with separate tokens
	OrValue      *QueryFilterValue `":" @@ )`
}

// String returns the string representation of a negated tag filter
func (ntf *NegatedTagFilter) String() string {
	if ntf == nil {
		return ""
	}

	if ntf.NegatedIdent != "" {
		return ntf.NegatedIdent + ":" + ntf.Value.String()
	}

	return "-" + ntf.Key + ":" + ntf.OrValue.String()
}

// Add a dedicated type for wildcard content search
type WildcardContent struct {
	Prefix  string `@"*"`
	Content string `( @Ident | @"/" )*`
	Suffix  string `( @"*" | "" )`
}

// String returns the string representation of a wildcard content search
func (wc *WildcardContent) String() string {
	if wc == nil {
		return ""
	}
	// Return "*CONTENT*" or "*CONTENT" depending on whether a suffix was matched
	if wc.Suffix == "*" {
		return "*" + wc.Content + "*"
	}
	return "*" + wc.Content
}

// Add support for CIDR function
type CIDRFunction struct {
	Attribute string   `"CIDR" "(" @Attribute`
	Blocks    []string `"," @(CIDR|IPAddress) ("," @(CIDR|IPAddress))* ")"`
}

// String returns the string representation of a CIDR function
func (cf *CIDRFunction) String() string {
	if cf == nil {
		return ""
	}

	blocks := strings.Join(cf.Blocks, ",")
	return "CIDR(" + cf.Attribute + "," + blocks + ")"
}

// NegatedSubexpression represents a negated parenthesized expression like -(expr)
type NegatedSubexpression struct {
	Pos lexer.Position

	// The inner expression that is being negated
	Subexpression *QueryFilter `"-" "(" @@ ")"`
}

// String returns the string representation of a negated subexpression
func (ns *NegatedSubexpression) String() string {
	if ns == nil {
		return ""
	}
	return "-(" + ns.Subexpression.String() + ")"
}

// Add support for NOT operator with a subexpression
type NotFunction struct {
	InnerQueryFilter *QueryFilter  `"NOT" "(" @@ ")"`
	InnerCIDR        *CIDRFunction `| "NOT" "(" @@ ")"`
}

// String returns the string representation of a NOT function
func (nf *NotFunction) String() string {
	if nf == nil {
		return ""
	}

	if nf.InnerQueryFilter != nil {
		return "NOT(" + nf.InnerQueryFilter.String() + ")"
	}

	if nf.InnerCIDR != nil {
		return "NOT(" + nf.InnerCIDR.String() + ")"
	}

	return ""
}

// NegatedQuotedText represents a negated quoted string like -"error message"
type NegatedQuotedText struct {
	Text string `"-" @String`
}

// String returns the string representation of a negated quoted text
func (nqt *NegatedQuotedText) String() string {
	if nqt == nil {
		return ""
	}
	return fmt.Sprintf("-\"%s\"", nqt.Text)
}

// FilterParam represents a parameter in a filter expression
type FilterParam struct {
	Subexpression        *QueryFilter            `( "(" @@ ")" )`
	Separator            *QueryFilterSeparator   `| @@`
	NotFunction          *NotFunction            `| @@`
	CIDRFunction         *CIDRFunction           `| @@`
	NegatedSubexpression *NegatedSubexpression   `| @@`
	NegatedQuotedText    *NegatedQuotedText      `| @@` // Add this line
	NegatedFilter        *NegatedFilter          `| "-" @@`
	NegatedTag           *NegatedTagFilter       `| @@`
	NegatedTerm          *NegatedTerm            `| @@`
	NegatedAttribute     *NegatedAttributeFilter `| @@`
	Attribute            *AttributeFilter        `| @@`
	Tag                  *TagFilter              `| @@`
	FullText             *FullTextFilter         `| @@`
	WildcardContent      *WildcardContent        `| @@`
	IPAddress            *string                 `| @IPAddress`
	NegativeFloat        *string                 `| @NegativeFloat`
	NegativeInt          *string                 `| @NegativeInt`
	Number               *float64                `| @Float`
	BoolLikeIdent        *string                 `| @BoolLikeIdent`
	Text                 *string                 `| @Ident`
	QuotedText           *string                 `| @String`
	CIDR                 *string                 `| @CIDR`
	EscapedBoolean       *string                 `| @EscapedBoolean`
}

// AttributeFilter represents an attribute filter
type AttributeFilter struct {
	Name  string            `@Attribute`
	Value *QueryFilterValue `":" @@`
}

// String returns the string representation of an attribute filter
func (af *AttributeFilter) String() string {
	if af == nil {
		return ""
	}

	return fmt.Sprintf("%s:%s", af.Name, af.Value.String())
}

// TagFilter represents a tag filter
type TagFilter struct {
	Key   string            `@Ident`
	Value *QueryFilterValue `":" @@`
}

// String returns the string representation of a tag filter
func (tf *TagFilter) String() string {
	if tf == nil {
		return ""
	}
	return fmt.Sprintf("%s:%s", tf.Key, tf.Value.String())
}

// FullTextFilter represents a full-text search filter
type FullTextFilter struct {
	Wildcard string            `@"*"`
	Value    *QueryFilterValue `":" @@`
}

// String returns the string representation of a full-text filter
func (ftf *FullTextFilter) String() string {
	if ftf == nil {
		return ""
	}

	return fmt.Sprintf("%s:%s", ftf.Wildcard, ftf.Value.String())
}

// QueryFilterSeparator represents a separator between filter terms
type QueryFilterSeparator struct {
	And bool `@("AND" | "and")`
	Or  bool `| @("OR" | "or")`
	Not bool `| @("NOT" | "not")`
}

// String returns the string representation of a query filter separator
func (fs *QueryFilterSeparator) String() string {
	if fs == nil {
		return ""
	}
	if fs.And {
		return "AND"
	}
	if fs.Or {
		return "OR"
	}
	if fs.Not {
		return "NOT"
	}
	return ""
}

// QueryFilterValue represents a value in a filter
type QueryFilterValue struct {
	Number              *float64     `@Float`
	NegativeFloat       *string      `| @NegativeFloat` // Direct negative float capture
	NegativeInt         *string      `| @NegativeInt`   // Direct negative integer capture
	Text                *string      `| @Ident`
	EscapedIdent        *string      `| @EscapedIdent`
	QuotedText          *string      `| @String`
	IPAddress           *string      `| @IPAddress`
	Wildcard            *string      `| @"*"`
	CIDR                *string      `| @CIDR`
	ComplexWildcardPath *string      `| @ComplexWildcardPath`
	Range               *RangeValue  `| @@`
	OpValue             *OpValue     `| @@`
	OperatorVal         *string      `| @OperatorValue`
	RangeVal            *string      `| @RangeValue`
	Subexpression       *QueryFilter `| "(" @@ ")"`
}

// String returns the string representation of a query filter value
func (fv *QueryFilterValue) String() string {
	if fv == nil {
		return ""
	}

	if fv.NegativeFloat != nil {
		return *fv.NegativeFloat
	}

	if fv.NegativeInt != nil {
		return *fv.NegativeInt
	}

	if fv.ComplexWildcardPath != nil {
		return *fv.ComplexWildcardPath
	}

	if fv.EscapedIdent != nil {
		// Special handling for escaped backslashes
		escaped := *fv.EscapedIdent
		result := strings.Builder{}
		i := 0

		for i < len(escaped) {
			if i+1 < len(escaped) && escaped[i] == '\\' {
				if escaped[i+1] == '\\' {
					// For escaped backslash (\\), keep one backslash
					result.WriteByte('\\')
					i += 2 // Skip both backslashes
				} else {
					// For other escaped characters, skip the backslash
					result.WriteByte(escaped[i+1])
					i += 2 // Skip the backslash and the character
				}
			} else {
				// Regular character, just copy it
				result.WriteByte(escaped[i])
				i++
			}
		}

		return result.String()
	}

	if fv.CIDR != nil {
		return *fv.CIDR
	}

	if fv.IPAddress != nil {
		return *fv.IPAddress
	}

	if fv.Number != nil {
		return fmt.Sprintf("%g", *fv.Number)
	}
	if fv.Text != nil {
		return *fv.Text
	}
	if fv.QuotedText != nil {
		return fmt.Sprintf("\"%s\"", *fv.QuotedText)
	}
	if fv.Wildcard != nil {
		return *fv.Wildcard
	}
	if fv.Range != nil {
		return fv.Range.String()
	}
	if fv.OpValue != nil {
		return fv.OpValue.String()
	}
	if fv.OperatorVal != nil {
		return *fv.OperatorVal
	}
	if fv.RangeVal != nil {
		return *fv.RangeVal
	}
	if fv.Subexpression != nil {
		return "(" + fv.Subexpression.String() + ")"
	}

	return ""
}

// OpValue represents an operator followed by a value
type OpValue struct {
	Value string `@OperatorValue`
}

// String returns the string representation of an OpValue
func (ov *OpValue) String() string {
	if ov == nil {
		return ""
	}
	return ov.Value
}

// RangeValue represents a range filter like [400 TO 499]
type RangeValue struct {
	OpenBracket  string  `@"["`
	Start        float64 `@Float`
	To           string  `@"TO"`
	End          float64 `@Float`
	CloseBracket string  `@"]"`
}

// String returns the string representation of a RangeValue
func (rv *RangeValue) String() string {
	if rv == nil {
		return ""
	}
	return fmt.Sprintf("[%.0f TO %.0f]", rv.Start, rv.End)
}

// NewQueryFilterParser returns a Parser which is capable of interpreting
// a Datadog query filter.
func NewQueryFilterParser() *QueryFilterParser {
	qfp := &QueryFilterParser{
		parser: participle.MustBuild[QueryFilter](
			participle.Lexer(lex),
			participle.Unquote("String"),
		),
	}

	return qfp
}

func DebugParse(query string) (*QueryFilter, error) {
	qfp := NewQueryFilterParser()
	result, err := qfp.Parse(query)

	if err != nil {
		return nil, err
	}

	// Debug output
	fmt.Printf("Query: %s\n", query)
	fmt.Printf("Left: %+v\n", result.Left)
	fmt.Printf("Right count: %d\n", len(result.Right))
	for i, r := range result.Right {
		fmt.Printf("Right[%d]: %+v\n", i, r)
	}

	return result, nil
}

// QueryFilterParser is parser returned when calling NewQueryFilterParser.
type QueryFilterParser struct {
	parser *participle.Parser[QueryFilter]
}

// Parse sanitizes the query string and returns the AST and any error.
func (qfp *QueryFilterParser) Parse(query string) (*QueryFilter, error) {
	// the parser doesn't handle queries that are split up across multiple lines
	sanitized := strings.ReplaceAll(query, "\n", "")

	// Pre-process to handle escaped characters within quoted strings
	sanitized = preprocessQuotedParts(sanitized)

	// return the raw parsed output
	return qfp.parser.ParseString("", sanitized)
}

// preprocessQuotedParts handles escaped characters inside quoted strings.
// The participle.Unquote("String") function expects valid Go escape sequences,
// but our custom format allows escaping any character including hyphens.
// This function converts escaped characters inside quoted strings to regular characters,
// while preserving escaped quotes.
func preprocessQuotedParts(query string) string {
	var result strings.Builder
	inQuote := false
	i := 0

	for i < len(query) {
		char := query[i]

		if char == '\\' && i+1 < len(query) {
			// Handle escape sequences
			nextChar := query[i+1]

			if inQuote {
				if nextChar == '\\' {
					// Inside quotes, preserve both backslashes for double-escaped sequences
					result.WriteByte('\\')
					result.WriteByte('\\')
				} else if nextChar == '"' {
					// Keep backslash for escaped quotes inside quotes
					result.WriteByte('\\')
					result.WriteByte('"')
				} else {
					// Inside quotes, remove backslash for other escape sequences
					result.WriteByte(nextChar)
				}
			} else {
				// Keep backslash for escaped sequences outside quotes
				result.WriteByte('\\')
				result.WriteByte(nextChar)
			}

			i += 2 // Skip both the backslash and the escaped character
			continue
		}

		if char == '"' {
			// Toggle quote state for unescaped quotes
			inQuote = !inQuote
		}

		result.WriteByte(char)
		i++
	}

	return result.String()
}

// String returns the string representation of the filter
func (qf *QueryFilter) String() string {
	if qf == nil || qf.Left == nil {
		return ""
	}

	// Special case for wildcard at the beginning of tag value
	if qf.Left.Tag != nil && qf.Left.Tag.Value != nil && qf.Left.Tag.Value.Wildcard != nil &&
		len(qf.Right) == 1 && qf.Right[0].Text != nil {
		return fmt.Sprintf("%s:*%s", qf.Left.Tag.Key, *qf.Right[0].Text)
	}

	out := []string{qf.Left.String()}

	// Handle the case where we have "NOT" followed by a parenthesized expression
	i := 0
	for i < len(qf.Right) {
		// Special case for tag:* followed by value*
		if i+1 < len(qf.Right) &&
			qf.Right[i].Tag != nil && qf.Right[i].Tag.Value != nil &&
			qf.Right[i].Tag.Value.Wildcard != nil && // tag:*
			qf.Right[i+1].Text != nil { // value*

			// Combine them into a single token without space
			out = append(out, fmt.Sprintf("%s:*%s", qf.Right[i].Tag.Key, *qf.Right[i+1].Text))
			i += 2 // Skip both elements
			continue
		} else if i+1 < len(qf.Right) &&
			qf.Right[i].Separator != nil && qf.Right[i].Separator.Not &&
			qf.Right[i+1].Subexpression != nil {
			// Found "NOT" followed by a parenthesized expression
			// Combine them without a space and add to output
			out = append(out, "NOT"+qf.Right[i+1].String())
			i += 2 // Skip both elements
		} else {
			// Regular element, add as is
			out = append(out, qf.Right[i].String())
			i++
		}
	}

	return strings.Join(out, " ")
}

func (fp *FilterParam) String() string {
	if fp == nil {
		return ""
	}

	if fp.NegatedQuotedText != nil {
		return fp.NegatedQuotedText.String()
	}

	if fp.NegatedSubexpression != nil {
		return fp.NegatedSubexpression.String()
	}

	if fp.NegatedFilter != nil {
		return fp.NegatedFilter.String()
	}

	if fp.Subexpression != nil {
		return "(" + fp.Subexpression.String() + ")"
	}
	if fp.Separator != nil {
		return fp.Separator.String()
	}
	if fp.NotFunction != nil {
		return fp.NotFunction.String()
	}
	if fp.CIDRFunction != nil {
		return fp.CIDRFunction.String()
	}
	if fp.NegatedTerm != nil {
		return fp.NegatedTerm.String()
	}
	if fp.NegatedAttribute != nil {
		return fp.NegatedAttribute.String()
	}
	if fp.NegatedTag != nil {
		return fp.NegatedTag.String()
	}
	if fp.Attribute != nil {
		return fp.Attribute.String()
	}
	if fp.Tag != nil {
		return fp.Tag.String()
	}
	if fp.FullText != nil {
		return fp.FullText.String()
	}
	if fp.WildcardContent != nil {
		return fp.WildcardContent.String()
	}
	if fp.IPAddress != nil {
		return *fp.IPAddress
	}
	if fp.NegativeFloat != nil {
		return *fp.NegativeFloat
	}

	if fp.NegativeInt != nil {
		return *fp.NegativeInt
	}
	if fp.BoolLikeIdent != nil {
		return *fp.BoolLikeIdent
	}
	if fp.Text != nil {
		return *fp.Text
	}
	if fp.QuotedText != nil {
		return fmt.Sprintf("\"%s\"", *fp.QuotedText)
	}
	if fp.Number != nil {
		return fmt.Sprintf("%g", *fp.Number)
	}
	if fp.CIDR != nil {
		return *fp.CIDR
	}
	if fp.EscapedBoolean != nil {
		// Remove the backslash from the escaped boolean
		return (*fp.EscapedBoolean)[1:]
	}
	return ""
}
