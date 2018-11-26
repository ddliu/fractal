// Copyright 2018 Liu Dong <ddliuhb@gmail.com>.
// Licensed under the MIT license.

package fractal

import (
	"testing"
)

func TestTemplate(t *testing.T) {
	tpl := `Author: ${author.name}; License: ${license}; Length: ${length()}`
	expected := `Author: Dong; License: MIT; Length: 2`

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
