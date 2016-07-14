package mustache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRenderText(t *testing.T) {
	tmpl, err := Compile("This is some text.")
	if !assert.NoError(t, err) || !assert.NotNil(t, tmpl) {
		return
	}
	contents, err := tmpl.Render(nil)
	if assert.NoError(t, err) {
		assert.Equal(t, "This is some text.", contents)
	}
}

func TestRenderVariable(t *testing.T) {
	tmpl, err := Compile("{{a}}")
	if !assert.NoError(t, err) || !assert.NotNil(t, tmpl) {
		return
	}
	contents, err := tmpl.Render(map[string]interface{}{"a": "b"})
	if assert.NoError(t, err) {
		assert.Equal(t, "b", contents)
	}
}
