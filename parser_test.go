package mustache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestText(t *testing.T) {
	tmpl, err := Compile("This is an example string")
	if assert.NoError(t, err) {
		if assert.Len(t, tmpl.result.tokens, 1) && assert.IsType(t, &text{}, tmpl.result.tokens[0]) {
			assert.Equal(t, "This is an example string", tmpl.result.tokens[0].(*text).value)
		}
	}
}

func TestVariable(t *testing.T) {
	templates := []string{
		"Welcome to {{ place }}!",
		"Welcome to {{ place}}!",
		"Welcome to {{{ place }}}!",
		"Welcome to {{{place }}}!",
	}
	for _, template := range templates {
		tmpl, err := Compile(template)
		if assert.NoError(t, err, template) && assert.Len(t, tmpl.result.tokens, 3, template) {
			if assert.IsType(t, &text{}, tmpl.result.tokens[0]) {
				assert.Equal(t, "Welcome to ", tmpl.result.tokens[0].(*text).value)
			}
			if assert.IsType(t, &variable{}, tmpl.result.tokens[1]) {
				assert.Equal(t, "place", tmpl.result.tokens[1].(*variable).name)
			}
			if assert.IsType(t, &text{}, tmpl.result.tokens[2]) {
				assert.Equal(t, "!", tmpl.result.tokens[2].(*text).value)
			}
		}
	}
}

func TestComment(t *testing.T) {
	templates := []string{
		"12{{! comment }}34",
		"12{{!comment }}34",
		"12{{! comment}}34",
		"12{{!comment}}34",
		"12{{! comment !}}34",
		"12{{!comment !}}34",
		"12{{! comment!}}34",
		"12{{!comment!}}34",
	}
	for _, template := range templates {
		tmpl, err := Compile(template)
		if assert.NoError(t, err, template) && assert.Len(t, tmpl.result.tokens, 2, template) {
			if assert.IsType(t, &text{}, tmpl.result.tokens[0]) {
				assert.Equal(t, "12", tmpl.result.tokens[0].(*text).value)
			}
			if assert.IsType(t, &text{}, tmpl.result.tokens[1]) {
				assert.Equal(t, "34", tmpl.result.tokens[1].(*text).value)
			}
		}
	}
}

func TestStandaloneComment(t *testing.T) {
	templates := []string{
		"12\n{{! comment }}\n34",
	}
	for _, template := range templates {
		tmpl, err := Compile(template)
		if assert.NoError(t, err, template) && assert.Len(t, tmpl.result.tokens, 2, template) {
			if assert.IsType(t, &text{}, tmpl.result.tokens[0]) {
				assert.Equal(t, "12\n", tmpl.result.tokens[0].(*text).value)
			}
			if assert.IsType(t, &text{}, tmpl.result.tokens[1]) {
				assert.Equal(t, "34", tmpl.result.tokens[1].(*text).value)
			}
		}
	}
}

func TestMultilineComment(t *testing.T) {
	templates := []string{
		"Begin.\n{{!\nSomething's going on here...\n}}\nEnd.",
		"Begin.\n{{! Something's going on here... }}\nEnd.",
		"Begin.\n  {{!\n    Something's going on here...\n  }}\nEnd.",
	}
	for _, template := range templates {
		tmpl, err := Compile(template)
		if assert.NoError(t, err, template) && assert.Len(t, tmpl.result.tokens, 2, template) {
			if assert.IsType(t, &text{}, tmpl.result.tokens[0]) {
				assert.Equal(t, "Begin.\n", tmpl.result.tokens[0].(*text).value)
			}
			if assert.IsType(t, &text{}, tmpl.result.tokens[1]) {
				assert.Equal(t, "End.", tmpl.result.tokens[1].(*text).value)
			}
		}
	}
}
