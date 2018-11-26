// Copyright 2018 Liu Dong <ddliuhb@gmail.com>.
// Licensed under the MIT license.

package fractal

import (
	"testing"
)

func TestTemplate(t *testing.T) {
	tpl := `Author: ${author.name}; License: ${license}; Length: ${length()}; Language: ${languages.0.name}`
	expected := `Author: Dong; License: MIT; Length: 3; Language: Golang`

	c := New(nil)
	c.SetValue("author", map[string]string{
		"name":  "Dong",
		"email": "test@example.com",
	})
	c.SetValue("license", "MIT")
	c.SetValue("languages", []map[string]string{
		map[string]string{
			"name": "Golang",
		},
		map[string]string{
			"name": "Python",
		},
	})

	parsed := c.Tpl(tpl)
	if parsed != expected {
		t.Error("Tpl error: " + parsed)
	}
}
