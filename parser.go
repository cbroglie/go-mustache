package mustache

import (
	"fmt"
	"regexp"
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
	comment // not exported as comment tags are not part of the final parse tree
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

type text struct {
	value string
}

type section struct {
	name     string
	inverted bool
	tokens   []Token
}

type variable struct {
	name   string
	escape bool
}

type partial struct {
	name  string
	value *section
}

// Template represents a compiled mustache Template. Its methods are safe for
// concurrent use.
type Template struct {
	result   *section
	sections []*section
	scanner  *stringScanner
	error    error
}

var (
	openTag        = regexp.MustCompile(`([ \t]*)?(\{\{)`)
	notOpenTag     = regexp.MustCompile(`(?m)(^[ \t]*)?(\{\{)`)
	tagType        = regexp.MustCompile(`(\!|\{)`)
	allowedContent = regexp.MustCompile(`(\w|[?!\/.-])*`)
	closeTag       = map[TokenType]*regexp.Regexp{
		Variable:          regexp.MustCompile(`([ \t]*)?(\}\})`),
		UnescapedVariable: regexp.MustCompile(`([ \t]*)?(\}\}\})`),
		comment:           regexp.MustCompile(`([ \t]*)?(\!?\}\})`),
	}
)

// Compile takes a string mustache Template and compiles it so that it can be
// rendered.
func Compile(contents string) (*Template, error) {
	t := &Template{}
	if err := t.parse(contents); err != nil {
		return nil, err
	}
	return t, nil
}

func (t *Template) Render(context interface{}) (string, error) {
	return "", nil
}

func (t *Template) Tokens() []Token {
	return t.result.tokens
}

func newSection() *section {
	return &section{
		tokens: make([]Token, 0),
	}
}

func (t *Template) parse(contents string) error {
	t.scanner = &stringScanner{input: contents}
	t.result = newSection()
	t.sections = make([]*section, 0)

	for t.error == nil && !t.scanner.Done() {
		if t.parseTags() {
			continue
		}
		t.parseText()
	}
	if t.error != nil {
		return t.error
	}

	// We have parsed the whole Template, but there are still open sections.
	if len(t.sections) != 0 {
		return fmt.Errorf("Unclosed section %q", t.sections[0].name)
	}

	return nil
}

func (t *Template) parseTokenType() (TokenType, error) {
	// Parse the next character after the opening tag "{{"
	// If it's not one of the special control characters assume it's a regular
	// variable tag.
	matches := t.scanner.Scan(tagType)
	if len(matches) == 0 {
		return Variable, nil
	}
	switch matches[0] {
	case "!":
		return comment, nil
	case "{":
		return UnescapedVariable, nil
	default:
		return 0, fmt.Errorf("Unexpected tag type %s", matches[0])
	}
}

func (t *Template) parseContent(tokenType TokenType) (string, error) {
	if tokenType == comment {
		matches := t.scanner.ScanUntil(closeTag[tokenType])
		// Backup the scan pointer to just before the match.
		if len(matches) > 0 {
			t.scanner.SetPos(t.scanner.Pos() - len(matches[1]))
			return matches[0], nil
		}
		return "", nil
	}

	matches := t.scanner.Scan(allowedContent)
	if len(matches) == 0 {
		return "", fmt.Errorf("Illegal content in tag")
	}
	return matches[0], nil
}

func (t *Template) addTokens(tokenType TokenType, content string) error {
	switch tokenType {
	case Variable, UnescapedVariable:
		t.result.tokens = append(t.result.tokens, &variable{name: content, escape: tokenType == Variable})
	}
	return nil
}

func (t *Template) parseTags() bool {
	startOfLine := t.scanner.StartOfLine()

	// Look for an opening tag.
	matches := t.scanner.Scan(openTag)
	if len(matches) == 0 {
		return false
	}
	if len(matches) != 3 {
		t.error = fmt.Errorf("Unexpected regex match %v", matches)
		return true
	}

	// If we're matching the start of a new line we hold off on adding the
	// whitespace; it may be skipped based on the type of tag we've matched.
	padding := matches[1]
	if !startOfLine && len(padding) > 0 {
		t.result.tokens = append(t.result.tokens, &text{value: padding})
	}

	// Scan ahead to figure out which kind of token this is.
	tokenType, err := t.parseTokenType()
	if err != nil {
		t.error = err
		return true
	}

	// Skip over any whitespace between the opening tag and the content.
	t.scanner.Scan(regexp.MustCompile(`\s*`))

	// Parse the content in the tag. The rules vary by type.
	content, err := t.parseContent(tokenType)
	if err != nil {
		t.error = err
		return true
	}

	// Add the token to the parse tree.
	err = t.addTokens(tokenType, content)
	if err != nil {
		t.error = err
		return true
	}

	// Skip over any whitespace between the content and the closing tag.
	t.scanner.Scan(regexp.MustCompile(`\s*`))

	// Find the closing tag.
	matches = t.scanner.Scan(closeTag[tokenType])
	if len(matches) == 0 {
		t.error = fmt.Errorf("Unclosed tag")
		return true
	}

	// If this tag was the only non-whitespace content on this line, strip the
	// remaining whitespace. If not, but we've been hanging on to padding from
	// the beginning of the line, re-insert the padding as static text.
	if startOfLine && !t.scanner.Done() {
		if skipWhitespace(tokenType) && t.scanner.Check(regexp.MustCompile(`[\t ]*\r?\n`)) != nil {
			t.scanner.Scan(regexp.MustCompile(`[\t ]*\r?\n`))
		} else if len(padding) > 0 {
			t.result.tokens = append(t.result.tokens, &text{value: padding})
		}
	}

	return false
}

func skipWhitespace(tokenType TokenType) bool {
	// After these types of tags, all whitespace until the end of the line will
	// be skipped if they are the first (and only) non-whitespace content on the
	// line.
	switch tokenType {
	case Section, InvertedSection, comment:
		return true
	}
	return true
}

func (t *Template) parseText() bool {
	if t.scanner.Done() {
		return false
	}

	// Scan up to the next open tag.
	matches := t.scanner.ScanUntil(notOpenTag)

	// No more open tags, add the remaining string as text.
	if len(matches) == 0 {
		rest, err := t.scanner.Substring(t.scanner.Pos(), t.scanner.Len())
		if err != nil {
			t.error = err
		} else {
			t.result.tokens = append(t.result.tokens, &text{value: rest})
			t.scanner.SetPos(t.scanner.Len())
		}
		return true
	}

	// Sanity check the regex.
	if len(matches) != 4 {
		t.error = fmt.Errorf("Unexpected regex match %v", matches)
		return true
	}

	// Backup the scan pointer to just before the match.
	t.scanner.SetPos(t.scanner.Pos() - len(matches[1]))

	// Add the text up to the match.
	t.result.tokens = append(t.result.tokens, &text{value: matches[0][0 : len(matches[0])-len(matches[1])]})
	return true
}

func (v *variable) Type() TokenType {
	if v.escape {
		return Variable
	}
	return UnescapedVariable
}

func (v *variable) Name() string {
	return v.name
}

func (v *variable) Tokens() []Token {
	if v.escape {
		panic("mustache: Tokens on Variable type")
	}
	panic("mustache: Tokens on UnescapedVariable type")
}

func (t *text) Type() TokenType {
	return Text
}

func (t *text) Name() string {
	panic("mustache: Name on Text type")
}

func (t *text) Tokens() []Token {
	panic("mustache: Tokens on Text type")
}
