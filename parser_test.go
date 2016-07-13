package mustache

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type parserTest struct {
	template string
	tokens   []Token
}

func runTests(t *testing.T, tests []parserTest) {
	for _, test := range tests {
		tmpl, err := Compile(test.template)
		if assert.NoError(t, err, test.template) && assert.Len(t, tmpl.result.tokens, len(test.tokens)) {
			for i := range test.tokens {
				assert.Equal(t, test.tokens[i], tmpl.result.tokens[i], fmt.Sprintf("tokens at index %d are not equal", i))
			}
		}
	}
}

func TestText(t *testing.T) {
	tokens := []Token{
		&text{value: "This is an example string"},
	}
	tests := []parserTest{
		{
			template: "This is an example string",
			tokens:   tokens,
		},
	}
	runTests(t, tests)
}

func TestVariable(t *testing.T) {
	escapedTokens := []Token{
		&text{value: "Welcome to "},
		&variable{name: "place", escape: true},
		&text{value: "!"},
	}
	unescapedTokens := []Token{
		&text{value: "Welcome to "},
		&variable{name: "place", escape: false},
		&text{value: "!"},
	}
	tests := []parserTest{
		{
			template: "Welcome to {{ place }}!",
			tokens:   escapedTokens,
		},
		{
			template: "Welcome to {{ place}}!",
			tokens:   escapedTokens,
		},
		{
			template: "Welcome to {{{ place }}}!",
			tokens:   unescapedTokens,
		},
		{
			template: "Welcome to {{{place }}}!",
			tokens:   unescapedTokens,
		},
	}
	runTests(t, tests)
}

func TestComment(t *testing.T) {
	tokens := []Token{
		&text{value: "12"},
		&text{value: "34"},
	}
	tests := []parserTest{
		{
			template: "12{{! comment }}34",
			tokens:   tokens,
		},
		{
			template: "12{{!comment }}34",
			tokens:   tokens,
		},
		{
			template: "12{{! comment}}34",
			tokens:   tokens,
		},
		{
			template: "12{{!comment}}34",
			tokens:   tokens,
		},
		{
			template: "12{{! comment !}}34",
			tokens:   tokens,
		},
		{
			template: "12{{!comment !}}34",
			tokens:   tokens,
		},
		{
			template: "12{{! comment!}}34",
			tokens:   tokens,
		},
		{
			template: "12{{!comment!}}34",
			tokens:   tokens,
		},
	}
	runTests(t, tests)
}

func TestStandaloneComment(t *testing.T) {
	tokens := []Token{
		&text{value: "12\n"},
		&text{value: "34"},
	}
	tests := []parserTest{
		{
			template: "12\n{{! comment }}\n34",
			tokens:   tokens,
		},
	}
	runTests(t, tests)
}

func TestMultilineComment(t *testing.T) {
	tokens := []Token{
		&text{value: "Begin.\n"},
		&text{value: "End."},
	}
	tests := []parserTest{
		{
			template: "Begin.\n{{!\nSomething's going on here...\n}}\nEnd.",
			tokens:   tokens,
		},
		{
			template: "Begin.\n{{! Something's going on here... }}\nEnd.",
			tokens:   tokens,
		},
		{
			template: "Begin.\n  {{!\n    Something's going on here...\n  }}\nEnd.",
			tokens:   tokens,
		},
	}
	runTests(t, tests)
}
