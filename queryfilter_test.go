package ddqp

import (
	"github.com/alecthomas/repr"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueryFilterParser(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		expected string
		wantErr  bool
		printAST bool // For debugging, can opt in to print AST
	}{
		// Basic boolean operators
		{
			name:     "AND operator",
			query:    `authentication AND failure`,
			expected: `authentication AND failure`,
			wantErr:  false,
		},
		{
			name:     "OR operator",
			query:    `authentication OR password`,
			expected: `authentication OR password`,
			wantErr:  false,
		},
		{
			name:     "exclusion operator",
			query:    `authentication AND -password`,
			expected: `authentication AND -password`,
			wantErr:  false,
		},

		// Full-text search
		{
			name:     "full-text search exact match",
			query:    `*:hello`,
			expected: `*:hello`,
			wantErr:  false,
		},
		{
			name:     "full-text search with wildcard",
			query:    `*:hello*`,
			expected: `*:hello*`,
			wantErr:  false,
		},
		{
			name:     "full-text search with quoted terms",
			query:    `*:"hello world"`,
			expected: `*:"hello world"`,
			wantErr:  false,
		},
		{
			name:     "free text search",
			query:    `hello world`,
			expected: `hello world`,
			wantErr:  false,
		},

		// Escaped special characters
		{
			name:     "escape plus character",
			query:    `@attribute:value\+extra`,
			expected: `@attribute:value+extra`,
			wantErr:  false,
		},
		{
			name:     "escape minus character",
			query:    `service:quota\-manager`,
			expected: `service:quota-manager`,
			wantErr:  false,
		},
		{
			name:     "escape equals character",
			query:    `@attribute:value\=equals`,
			expected: `@attribute:value=equals`,
			wantErr:  false,
		},
		{
			name:     "escape ampersand character",
			query:    `@attribute:value\&more`,
			expected: `@attribute:value&more`,
			wantErr:  false,
		},
		{
			name:     "escape pipe character",
			query:    `@attribute:value\|other`,
			expected: `@attribute:value|other`,
			wantErr:  false,
		},
		{
			name:     "escape greater than character",
			query:    `@attribute:value\>10`,
			expected: `@attribute:value>10`,
			wantErr:  false,
		},
		{
			name:     "escape less than character",
			query:    `@attribute:value\<10`,
			expected: `@attribute:value<10`,
			wantErr:  false,
		},
		{
			name:     "escape exclamation mark",
			query:    `@attribute:value\!important`,
			expected: `@attribute:value!important`,
			wantErr:  false,
		},
		{
			name:     "escape parentheses",
			query:    `@attribute:value\(in-parentheses\)`,
			expected: `@attribute:value(in-parentheses)`,
			wantErr:  false,
		},
		{
			name:     "escape curly braces",
			query:    `@attribute:value\{in-braces\}`,
			expected: `@attribute:value{in-braces}`,
			wantErr:  false,
		},
		{
			name:     "escape square brackets",
			query:    `@attribute:value\[in-brackets\]`,
			expected: `@attribute:value[in-brackets]`,
			wantErr:  false,
		},
		{
			name:     "escape caret character",
			query:    `@attribute:value\^power`,
			expected: `@attribute:value^power`,
			wantErr:  false,
		},
		{
			name:     "escape quotes",
			query:    `@attribute:value\"quoted\"`,
			expected: `@attribute:value"quoted"`,
			wantErr:  false,
		},
		{
			name:     "escape tilde character",
			query:    `@attribute:value\~approx`,
			expected: `@attribute:value~approx`,
			wantErr:  false,
		},
		{
			name:     "escape asterisk character",
			query:    `@attribute:value\*wildcard`,
			expected: `@attribute:value*wildcard`,
			wantErr:  false,
		},
		{
			name:     "escape question mark character",
			query:    `@attribute:value\?uncertain`,
			expected: `@attribute:value?uncertain`,
			wantErr:  false,
		},
		{
			name:     "escape colon character",
			query:    `@attribute:value\:more`,
			expected: `@attribute:value:more`,
			wantErr:  false,
		},
		{
			name:     "escape backslash character",
			query:    `@attribute:value\\escaped`,
			expected: `@attribute:value\escaped`,
			wantErr:  false,
		},
		{
			name:     "escape hash character",
			query:    `@attribute:value\#number`,
			expected: `@attribute:value#number`,
			wantErr:  false,
		},
		{
			name:     "escape space character",
			query:    `@attribute:value\ with\ spaces`,
			expected: `@attribute:value with spaces`,
			wantErr:  false,
		},
		{
			name:     "multiple escaped spaces",
			query:    `@message:this\ is\ a\ test\ message`,
			expected: `@message:this is a test message`,
			wantErr:  false,
		},

		// Attribute search
		{
			name:     "simple attribute search",
			query:    `@url:www.datadoghq.com`,
			expected: `@url:www.datadoghq.com`,
			wantErr:  false,
		},
		{
			name:     "attribute with quoted value",
			query:    `@http.url_details.path:"/api/v1/test"`,
			expected: `@http.url_details.path:"/api/v1/test"`,
			wantErr:  false,
		},
		{
			name:     "attribute with wildcard",
			query:    `@http.url:/api\-v1/*`,
			expected: `@http.url:/api-v1/*`,
			wantErr:  false,
		},
		{
			name:     "complex attribute search with range",
			query:    `@http.status_code:[200 TO 299] @http.url_details.path:/api\-v1/*`,
			expected: `@http.status_code:[200 TO 299] @http.url_details.path:/api-v1/*`,
			wantErr:  false,
		},
		{
			name:     "negated attribute existence",
			query:    `-@http.status_code:*`,
			expected: `-@http.status_code:*`,
			wantErr:  false,
		},
		{
			name:     "attribute with question mark wildcard",
			query:    `@my_attribute:hello?world`,
			expected: `@my_attribute:hello?world`,
			wantErr:  false,
		},

		// CIDR notation
		{
			name:     "CIDR function with single block",
			query:    `CIDR(@network.client.ip,13.0.0.0/8)`,
			expected: `CIDR(@network.client.ip,13.0.0.0/8)`,
			wantErr:  false,
		},
		{
			name:     "CIDR function with multiple blocks",
			query:    `CIDR(@network.ip.list,13.0.0.0/8, 15.0.0.0/8)`,
			expected: `CIDR(@network.ip.list,13.0.0.0/8,15.0.0.0/8)`,
			wantErr:  false,
		},
		{
			name:     "CIDR function with additional filters",
			query:    `source:pan.firewall evt.name:reject CIDR(@network.client.ip, 13.0.0.0/8)`,
			expected: `source:pan.firewall evt.name:reject CIDR(@network.client.ip,13.0.0.0/8)`,
			wantErr:  false,
		},
		{
			name:     "CIDR function with negation",
			query:    `source:vpc NOT(CIDR(@network.client.ip, 13.0.0.0/8)) CIDR(@network.destination.ip, 15.0.0.0/8)`,
			expected: `source:vpc NOT(CIDR(@network.client.ip,13.0.0.0/8)) CIDR(@network.destination.ip,15.0.0.0/8)`,
			wantErr:  false,
		},

		// Wildcards
		{
			name:     "multi-character wildcard at end",
			query:    `service:web*`,
			expected: `service:web*`,
			wantErr:  false,
		},
		{
			name:     "multi-character wildcard at beginning",
			query:    `service:*mongo`,
			expected: `service:*mongo`,
			wantErr:  false,
		},
		{
			name:     "wildcard in content search",
			query:    `*NETWORK*`,
			expected: `*NETWORK*`,
			wantErr:  false,
		},
		{
			name:     "quoted wildcard (literal asterisk)",
			query:    `"*test*"`,
			expected: `"*test*"`,
			wantErr:  false,
		},

		// Numerical values
		{
			name:     "numerical greater than",
			query:    `@http.response_time:>100`,
			expected: `@http.response_time:>100`,
			wantErr:  false,
		},
		{
			name:     "numerical less than",
			query:    `@http.response_time:<100`,
			expected: `@http.response_time:<100`,
			wantErr:  false,
		},
		{
			name:     "numerical greater than or equal",
			query:    `@http.response_time:>=100`,
			expected: `@http.response_time:>=100`,
			wantErr:  false,
		},
		{
			name:     "numerical less than or equal",
			query:    `@http.response_time:<=100`,
			expected: `@http.response_time:<=100`,
			wantErr:  false,
		},
		{
			name:     "numerical range",
			query:    `@http.status_code:[400 TO 499]`,
			expected: `@http.status_code:[400 TO 499]`,
			wantErr:  false,
		},

		// Tags
		{
			name:     "single tag",
			query:    `env:prod`,
			expected: `env:prod`,
			wantErr:  false,
		},
		{
			name:     "OR tags",
			query:    `env:(prod OR test)`,
			expected: `env:(prod OR test)`,
			wantErr:  false,
		},
		{
			name:     "AND tags with negation",
			query:    `(env:prod AND -version:beta)`,
			expected: `(env:prod AND -version:beta)`,
			wantErr:  false,
		},
		{
			name:     "tag without key-value",
			query:    `tags:MY_TAG`,
			expected: `tags:MY_TAG`,
			wantErr:  false,
		},

		// Arrays and JSON objects
		{
			name:     "array attribute search",
			query:    `@users.names:Peter`,
			expected: `@users.names:Peter`,
			wantErr:  false,
		},
		{
			name:     "nested JSON object search",
			query:    `@Event.EventData.Data.Name:ObjectServer`,
			expected: `@Event.EventData.Data.Name:ObjectServer`,
			wantErr:  false,
		},

		// Calculated fields
		{
			name:     "calculated field search",
			query:    `#request_duration:<100`,
			expected: `#request_duration:<100`,
			wantErr:  true,
		},
		{
			name:     "calculated field with boolean operator",
			query:    `#request_duration:>100 AND service:api`,
			expected: `#request_duration:>100 AND service:api`,
			wantErr:  true,
		},

		// Complex query combinations
		{
			name:     "complex query with multiple operators",
			query:    `(service:web OR service:api) AND @http.status_code:[500 TO 599] AND -@error:timeout`,
			expected: `(service:web OR service:api) AND @http.status_code:[500 TO 599] AND -@error:timeout`,
			wantErr:  false,
		},
		{
			name:     "complex query with escaped characters",
			query:    `service:payment\-api @path:"/api/v1/process\-payment" @http.status_code:>=400`,
			expected: `service:payment-api @path:"/api/v1/process-payment" @http.status_code:>=400`,
			wantErr:  false,
		},
		{
			name:     "complex query with multiple escaped operators",
			query:    `tag:value\=equals \AND not-an-operator \OR another-term`,
			expected: `tag:value=equals AND not-an-operator OR another-term`,
			wantErr:  false,
		},
		{
			name:     "complex boolean expression with escaped components",
			query:    `(service:api AND env:prod) OR (service:web \AND not-a-boolean)`,
			expected: `(service:api AND env:prod) OR (service:web AND not-a-boolean)`,
			wantErr:  false,
		},
		{
			name:     "complex query with mixed features",
			query:    `@http.url:"/api\-v1/users/*" @http.status_code:[400 TO 499] *:"error message" -service:legacy`,
			expected: `@http.url:"/api-v1/users/*" @http.status_code:[400 TO 499] *:"error message" -service:legacy`,
			wantErr:  false,
		},

		// Edge cases from documentation
		{
			name:     "query with multiple attribute filters",
			query:    `service:quota-manager status:error host:blabla @blabla:ff David`,
			expected: `service:quota-manager status:error host:blabla @blabla:ff David`,
			wantErr:  false,
		},
		{
			name:     "query with quoted strings",
			query:    `service:quota-manager status:error host:blabla @blabla:ff "skipping quota" "as it is inactive"`,
			expected: `service:quota-manager status:error host:blabla @blabla:ff "skipping quota" "as it is inactive"`,
			wantErr:  false,
		},
		{
			name:     "complex path with escaped characters",
			query:    `service:quota\-manager status:error host:blabla @path:"/api\-v1/quota\-check"`,
			expected: `service:quota-manager status:error host:blabla @path:"/api-v1/quota-check"`,
			wantErr:  false,
		},
		{
			name:     "quoted string with escaped quotes",
			query:    `@message:"This has an \"escaped\" quote"`,
			expected: `@message:"This has an "escaped" quote"`,
			wantErr:  false,
		},

		// Multiple nested subexpressions with various operators
		{
			name:     "deeply nested subexpressions",
			query:    `((service:api AND @status_code:200) OR (service:web AND (@response_time:<300 OR @error:*)))`,
			expected: `((service:api AND @status_code:200) OR (service:web AND (@response_time:<300 OR @error:*)))`,
			wantErr:  false,
		},

		// Complex queries with multiple negations at different levels
		{
			name:     "multiple negations at different levels",
			query:    `service:api -env:dev -@status_code:[500 TO 599] NOT(source:internal AND -region:us-east)`,
			expected: `service:api -env:dev -@status_code:[500 TO 599] NOT(source:internal AND -region:us-east)`,
			wantErr:  false,
		},

		// Mixed wildcard usage across different filter types
		{
			name:     "mixed wildcards across filter types",
			query:    `service:api* @path:*/users/*/* *:"error*" -tag:temp*`,
			expected: `service:api* @path:*/users/*/* *:"error*" -tag:temp*`,
			wantErr:  false,
		},

		// Complex combination of CIDR and other filters
		{
			name:     "complex CIDR combination",
			query:    `CIDR(@network.client.ip,10.0.0.0/8) AND NOT(CIDR(@network.destination.ip,192.168.0.0/16)) service:security*`,
			expected: `CIDR(@network.client.ip,10.0.0.0/8) AND NOT(CIDR(@network.destination.ip,192.168.0.0/16)) service:security*`,
			wantErr:  false,
		},

		// Query with multiple attributes, tags, and full-text search with various operators
		{
			name:     "kitchen sink query",
			query:    `(@http.method:GET OR @http.method:POST) @duration:>500 service:api-* env:(prod OR staging) -version:beta *:"timeout exceeded" NOT(@error:connection_reset)`,
			expected: `(@http.method:GET OR @http.method:POST) @duration:>500 service:api-* env:(prod OR staging) -version:beta *:"timeout exceeded" NOT(@error:connection_reset)`,
			wantErr:  false,
		},

		// Complex escaped characters in various positions
		{
			name:     "complex escaped characters",
			query:    `service:payment\-gateway @path:"/api/v2/process\-payment\?debug\=true" tag:value\=with\=equals \AND literal \OR literal`,
			expected: `service:payment-gateway @path:"/api/v2/process-payment?debug=true" tag:value=with=equals AND literal OR literal`,
			wantErr:  false,
		},

		// Multiple number ranges and operators in one query
		{
			name:     "multiple numerical filters",
			query:    `@status_code:[400 TO 499] @latency:>200 @retry_count:<5 @success_rate:>=0.95 @duration:[100 TO 500]`,
			expected: `@status_code:[400 TO 499] @latency:>200 @retry_count:<5 @success_rate:>=0.95 @duration:[100 TO 500]`,
			wantErr:  false,
		},

		// Complex combinations with quotes and escapes
		{
			name:     "quotes and escapes combination",
			query:    `@message:"User \"john.doe\" attempted to access \"/restricted\\area\"" service:auth\-service`,
			expected: `@message:"User "john.doe" attempted to access "/restricted\area"" service:auth-service`,
			wantErr:  false,
		},

		// Complex boolean logic with parentheses
		{
			name:     "complex boolean logic",
			query:    `((env:prod AND NOT(region:eu-*)) OR (env:staging AND service:api)) AND @duration:>100`,
			expected: `((env:prod AND NOT(region:eu-*)) OR (env:staging AND service:api)) AND @duration:>100`,
			wantErr:  false,
		},

		// Mixed case operators and identifiers
		{
			name:     "mixed case operators",
			query:    `service:API and @Status_Code:200 OR environment:Prod`,
			expected: `service:API AND @Status_Code:200 OR environment:Prod`,
			wantErr:  false,
		},

		// Complex negation within subexpressions
		{
			name:     "complex negation in subexpressions",
			query:    `(service:api AND -(env:dev OR region:eu-west)) OR (@status:error AND -(@error_type:timeout OR @error_type:connection))`,
			expected: `(service:api AND -(env:dev OR region:eu-west)) OR (@status:error AND -(@error_type:timeout OR @error_type:connection))`,
			wantErr:  false,
		},

		// Extreme wildcard usage
		{
			name:     "extreme wildcard usage",
			query:    `*:* -service:* @status_code:* tag:*value* @path:*/v1/*/*`,
			expected: `*:* -service:* @status_code:* tag:*value* @path:*/v1/*/*`,
			wantErr:  false,
		},

		// Multiple escaped booleans
		{
			name:     "multiple escaped booleans",
			query:    `\AND in-text \OR more-text \NOT another-term service:api`,
			expected: `AND in-text OR more-text NOT another-term service:api`,
			wantErr:  false,
		},

		// Multiple consecutive negations
		{
			name:     "multiple consecutive negations",
			query:    `-service:api -@status:error -env:dev -@duration:>500 -"error message"`,
			expected: `-service:api -@status:error -env:dev -@duration:>500 -"error message"`,
			wantErr:  false,
		},

		// Combination of quoted strings with wildcards
		{
			name:     "quoted strings with wildcards",
			query:    `@message:"*important* alert*" @path:"*/secure/*" service:"api-*-service"`,
			expected: `@message:"*important* alert*" @path:"*/secure/*" service:"api-*-service"`,
			wantErr:  false,
		},

		// Escaped spaces in identifiers and values
		{
			name:     "escaped spaces in identifiers",
			query:    `service:my\ service\ name @path:"/my\ path\ with\ spaces" tag:value\ with\ spaces`,
			expected: `service:my service name @path:"/my path with spaces" tag:value with spaces`,
			wantErr:  false,
		},

		// Combinations of IP addresses and CIDR
		{
			name:     "IP addresses and CIDR combinations",
			query:    `@network.src_ip:192.168.1.1 @network.dst_ip:10.0.0.1 CIDR(@network.subnet,172.16.0.0/16,192.168.0.0/16)`,
			expected: `@network.src_ip:192.168.1.1 @network.dst_ip:10.0.0.1 CIDR(@network.subnet,172.16.0.0/16,192.168.0.0/16)`,
			wantErr:  false,
		},

		// Edge case with many operators chained together
		{
			name:     "many chained operators",
			query:    `@a:1 AND @b:2 OR @c:3 AND @d:4 OR @e:5 AND NOT @f:6 OR @g:7`,
			expected: `@a:1 AND @b:2 OR @c:3 AND @d:4 OR @e:5 AND NOT @f:6 OR @g:7`,
			wantErr:  false,
		},

		// Really complex quote and escape patterns
		{
			name:     "complex quote and escape patterns",
			query:    `@json.field:"{ \"key\": \"value with \\"nested\\" quotes\" }"`,
			expected: `@json.field:"{ "key": "value with \"nested\" quotes" }"`,
			wantErr:  true,
		},

		// Edge case with number formats
		{
			name:     "various number formats",
			query:    `@value:42 @decimal:3.14159 @negative:-10.5 @scientific:1.23e-4 negativeint:-14`,
			expected: `@value:42 @decimal:3.14159 @negative:-10.5 @scientific:0.000123 negativeint:-14`,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qfp := NewQueryFilterParser()
			result, err := qfp.Parse(tt.query)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			if tt.printAST {
				repr.Println(result)
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result.String())
		})
	}
}
