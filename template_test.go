package fractal

import (
	"testing"
)

func TestTemplate(t *testing.T) {
	tpl := `Author: ${author.name}; License: ${license}; `
	expected := `Author: Dong; License: MIT; `

	c := New(nil)
	c.SetValue("author", map[string]string{
		"name":  "Dong",
		"email": "test@example.com",
	})
	c.SetValue("license", "MIT")

	parsed := c.Tpl(tpl)
	if parsed != expected {
		t.Error("Tpl error: " + parsed)
	}
}
