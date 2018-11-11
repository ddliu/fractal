// Copyright 2018 Liu Dong <ddliuhb@gmail.com>.
// Licensed under the MIT license.

package fractal

import (
	"github.com/spf13/cast"
)

func (c *Context) Bool(paths ...string) bool {
	return cast.ToBool(c.GetValue(paths...))
}

func (c *Context) Float32(paths ...string) float32 {
	return cast.ToFloat32(c.GetValue(paths...))
}

func (c *Context) Float64(paths ...string) float64 {
	return cast.ToFloat64(c.GetValue(paths...))
}

func (c *Context) Int(paths ...string) int {
	return cast.ToInt(c.GetValue(paths...))
}

func (c *Context) Int16(paths ...string) int16 {
	return cast.ToInt16(c.GetValue(paths...))
}

func (c *Context) Int32(paths ...string) int32 {
	return cast.ToInt32(c.GetValue(paths...))
}

func (c *Context) Int64(paths ...string) int64 {
	return cast.ToInt64(c.GetValue(paths...))
}

func (c *Context) Int8(paths ...string) int8 {
	return cast.ToInt8(c.GetValue(paths...))
}

func (c *Context) Uint(paths ...string) uint {
	return cast.ToUint(c.GetValue(paths...))
}

func (c *Context) Uint16(paths ...string) uint16 {
	return cast.ToUint16(c.GetValue(paths...))
}

func (c *Context) Uint32(paths ...string) uint32 {
	return cast.ToUint32(c.GetValue(paths...))
}

func (c *Context) Uint64(paths ...string) uint64 {
	return cast.ToUint64(c.GetValue(paths...))
}

func (c *Context) Uint8(paths ...string) uint8 {
	return cast.ToUint8(c.GetValue(paths...))
}

func (c *Context) String(paths ...string) string {
	return cast.ToString(c.GetValue(paths...))
}
