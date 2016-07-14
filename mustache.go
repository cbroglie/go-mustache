package mustache

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"reflect"
)

// TokenType represents the different tokens the parser extracts from a string
// Template.
type TokenType uint

// Defines for the possible TokenType values.
const (
	Text TokenType = iota
	Section
	InvertedSection
	Variable
	UnescapedVariable
	Partial
	comment      // not exported as comment tags are not part of the final parse tree
	closeSection // not exported as section close tags are not part of the final parse tree
)

// Token represents a mustache token.
//
// Not all methods apply to all kinds of tokens. Restrictions, if any, are noted
// in the documentation for each method. Use the Type method to find out the
// type of token before calling type-specific methods. Calling a method
// inappropriate to the type of token causes a run time panic.
type Token interface {
	// Type returns the type of the token.
	Type() TokenType
	// Name returns the name of the token. It panics for token types which are
	// not named (i.e. text tokens).
	Name() string
	// Tokens returns any child tokens. It panics for token types which cannot
	// contain child tokens (i.e. text and variable tokens).
	Tokens() []Token
}

// Template represents a compiled mustache Template. Its methods are safe for
// concurrent use.
type Template struct {
	result   *section
	sections []*section
	scanner  *stringScanner
	error    error
}

// Compile takes a string mustache template and compiles it so that it can be
// rendered.
func Compile(contents string) (*Template, error) {
	t := &Template{}
	if err := t.parse(contents); err != nil {
		return nil, err
	}
	return t, nil
}

// Render uses the given context to render a compiled mustache template.
func (t *Template) Render(context interface{}) (string, error) {
	var buf bytes.Buffer
	err := t.renderSection(&buf, t.result, reflect.ValueOf(context))
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

// Tokens returns the list of tokens from a parsed mustache template.
func (t *Template) Tokens() []Token {
	return t.result.tokens
}

func (t *Template) renderSection(out io.Writer, s *section, context reflect.Value) error {
	for _, token := range s.tokens {
		switch v := token.(type) {
		case *text:
			out.Write([]byte(v.value))
		case *variable:
			rv, err := lookup(context, v.name)
			if err != nil {
				return err
			}
			if rv.IsValid() {
				if v.escape {
					template.HTMLEscape(out, []byte(fmt.Sprint(rv.Interface())))
				} else {
					fmt.Fprint(out, rv.Interface())
				}
			}
		}
	}
	return nil
}

func lookup(context reflect.Value, name string) (reflect.Value, error) {
	k := context.Kind()
	if k != reflect.Map {
		return reflect.Value{}, fmt.Errorf("Only map is supported, found context with type %s", k)
	}
	v := context.MapIndex(reflect.ValueOf(name))
	if v.IsValid() {
		return v, nil
	}
	return reflect.Value{}, nil
}
