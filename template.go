// Copyright 2018 Liu Dong <ddliuhb@gmail.com>.
// Licensed under the MIT license.

package fractal

import (
	"regexp"
)

var re = regexp.MustCompile(`\$\{[a-zA-Z0-9\._]+\}`)

func (c *Context) Tpl(tpl string) string {
	return re.ReplaceAllStringFunc(tpl, c.rep)
}

func (c *Context) rep(str string) string {
	str = str[2 : len(str)-1]
	return c.String(str)
}
