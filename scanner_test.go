package mustache

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScan(t *testing.T) {
	scanner := &stringScanner{input: `This is an example string`}

	assert.Equal(t, "This", first(scanner.scan(regexp.MustCompile(`(\w+)`))))
	assert.Nil(t, scanner.scan(regexp.MustCompile(`\w+`)))
	assert.Equal(t, " ", first(scanner.scan(regexp.MustCompile(`(\s+)`))))
	assert.Nil(t, scanner.scan(regexp.MustCompile(`\s+`)))
	assert.Equal(t, "is", second(scanner.scan(regexp.MustCompile(`(\w+)`))))
	assert.False(t, scanner.done())

	assert.Equal(t, " ", first(scanner.scan(regexp.MustCompile(`\s+`))))
	assert.Equal(t, "an", second(scanner.scan(regexp.MustCompile(`(\w+)`))))
	assert.Equal(t, " ", first(scanner.scan(regexp.MustCompile(`\s+`))))
	assert.Equal(t, "example", first(scanner.scan(regexp.MustCompile(`(\w+)`))))
	assert.Equal(t, " ", second(scanner.scan(regexp.MustCompile(`(\s+)`))))
	assert.Equal(t, "string", first(scanner.scan(regexp.MustCompile(`(\w+)`))))
	assert.True(t, scanner.done())

	assert.Nil(t, scanner.scan(regexp.MustCompile(`\s+`)))
	assert.Nil(t, scanner.scan(regexp.MustCompile(`\w+`)))
}

func TestScanUntil(t *testing.T) {
	scanner := &stringScanner{input: `Fri Dec 12 1975 14:39`}
	assert.Equal(t, "Fri Dec 1", first(scanner.scanUntil(regexp.MustCompile(`1`))))
	assert.Nil(t, scanner.scanUntil(regexp.MustCompile(`XYZ`)))
	assert.False(t, scanner.done())
	assert.Equal(t, "2 1975 14:39", first(scanner.scanUntil(regexp.MustCompile(`(.*)`))))
	assert.True(t, scanner.done())
}

func TestScanUntilMultiline(t *testing.T) {
	scanner := &stringScanner{input: `<html>
  <body>
  {{!
    this is a comment
  !}}
  </body>
</html>`}

	matches := scanner.scanUntil(regexp.MustCompile(`(?m)(^[ \t]*)?(\{\{)`))
	if assert.NotNil(t, matches) && assert.Len(t, matches, 4) {
		assert.Equal(t, "<html>\n  <body>\n  {{", matches[0])
		assert.Equal(t, "  {{", matches[1])
		assert.Equal(t, "  ", matches[2])
		assert.Equal(t, "{{", matches[3])
	}
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
