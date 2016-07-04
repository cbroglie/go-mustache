package mustache

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScan(t *testing.T) {
	scanner := &stringScanner{input: `This is an example string`}

	assert.Equal(t, "This", first(scanner.Scan(regexp.MustCompile(`(\w+)`))))
	assert.Nil(t, scanner.Scan(regexp.MustCompile(`\w+`)))
	assert.Equal(t, " ", first(scanner.Scan(regexp.MustCompile(`(\s+)`))))
	assert.Nil(t, scanner.Scan(regexp.MustCompile(`\s+`)))
	assert.Equal(t, "is", second(scanner.Scan(regexp.MustCompile(`(\w+)`))))
	assert.False(t, scanner.Done())

	assert.Equal(t, " ", first(scanner.Scan(regexp.MustCompile(`\s+`))))
	assert.Equal(t, "an", second(scanner.Scan(regexp.MustCompile(`(\w+)`))))
	assert.Equal(t, " ", first(scanner.Scan(regexp.MustCompile(`\s+`))))
	assert.Equal(t, "example", first(scanner.Scan(regexp.MustCompile(`(\w+)`))))
	assert.Equal(t, " ", second(scanner.Scan(regexp.MustCompile(`(\s+)`))))
	assert.Equal(t, "string", first(scanner.Scan(regexp.MustCompile(`(\w+)`))))
	assert.True(t, scanner.Done())

	assert.Nil(t, scanner.Scan(regexp.MustCompile(`\s+`)))
	assert.Nil(t, scanner.Scan(regexp.MustCompile(`\w+`)))
}

func TestScanUntil(t *testing.T) {
	scanner := &stringScanner{input: `Fri Dec 12 1975 14:39`}
	assert.Equal(t, "Fri Dec 1", first(scanner.ScanUntil(regexp.MustCompile(`1`))))
	assert.Nil(t, scanner.ScanUntil(regexp.MustCompile(`XYZ`)))
	assert.False(t, scanner.Done())
	assert.Equal(t, "2 1975 14:39", first(scanner.ScanUntil(regexp.MustCompile(`(.*)`))))
	assert.True(t, scanner.Done())
}

func TestScanUntilMultiline(t *testing.T) {
	scanner := &stringScanner{input: `<html>
  <body>
  {{!
    this is a comment
  !}}
  </body>
</html>`}

	matches := scanner.CheckUntil(regexp.MustCompile(`(?m)(^[ \t]*)?(\{\{)`))
	if assert.NotNil(t, matches) && assert.Len(t, matches, 4) {
		assert.Equal(t, "<html>\n  <body>\n  {{", matches[0])
		assert.Equal(t, "  {{", matches[1])
		assert.Equal(t, "  ", matches[2])
		assert.Equal(t, "{{", matches[3])
	}
	assert.Equal(t, 0, scanner.Pos())

	matches = scanner.ScanUntil(regexp.MustCompile(`(?m)(^[ \t]*)?(\{\{)`))
	if assert.NotNil(t, matches) && assert.Len(t, matches, 4) {
		assert.Equal(t, "<html>\n  <body>\n  {{", matches[0])
		assert.Equal(t, "  {{", matches[1])
		assert.Equal(t, "  ", matches[2])
		assert.Equal(t, "{{", matches[3])
	}
}

func TestSubstring(t *testing.T) {
	scanner := &stringScanner{input: `<html>
  <body>
  {{!
    this is a comment
  !}}
  </body>
</html>`}
	s, err := scanner.Substring(1, 16)
	if assert.NoError(t, err) {
		assert.Equal(t, "html>\n  <body>\n", s)
	}
	s, err = scanner.Substring(1, 100)
	if assert.Error(t, err) {
		assert.Equal(t, "", s)
	}
}

func TestSetPos(t *testing.T) {
	const s = "test string"
	scanner := &stringScanner{input: s}
	assert.Equal(t, 0, scanner.Pos())
	assert.NoError(t, scanner.SetPos(3))
	assert.Equal(t, 3, scanner.Pos())
	assert.NoError(t, scanner.SetPos(len(s)))
	assert.Equal(t, len(s), scanner.Pos())
	assert.Error(t, scanner.SetPos(len(s)+1))
	assert.Error(t, scanner.SetPos(-1))
	assert.Equal(t, len(s), scanner.Pos())
}

func TestLen(t *testing.T) {
	const s = "test string"
	scanner := &stringScanner{input: s}
	assert.Equal(t, len(s), scanner.Len())
}

func first(v []string) string {
	if v == nil || len(v) == 0 {
		return ""
	}
	return v[0]
}

func second(v []string) string {
	if v == nil || len(v) <= 1 {
		return ""
	}
	return v[1]
}

func TestStartOfLine(t *testing.T) {
	scanner := &stringScanner{input: "test\ntest\n"}
	assert.True(t, scanner.StartOfLine())
	scanner.Scan(regexp.MustCompile(`te`))
	assert.False(t, scanner.StartOfLine())
	scanner.Scan(regexp.MustCompile(`st\n`))
	assert.True(t, scanner.StartOfLine())
	scanner.Scan(regexp.MustCompile(`te`))
	assert.False(t, scanner.StartOfLine())
}
