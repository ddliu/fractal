// Copyright 2018-2019 Liu Dong <ddliuhb@gmail.com>.
// Licensed under the MIT license.

package fractal

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"

	"github.com/spf13/cast"
)

func New(data interface{}) *Context {
	return &Context{
		data: data,
	}
}

func FromJson(jsonData []byte) *Context {
	var data interface{}
	json.Unmarshal(jsonData, &data)

	return New(data)
}

// Context container
type Context struct {
	data interface{}
}

func valueOfContext(v interface{}) (interface{}, error) {
	if v == nil {
		return nil, nil
	}

	if vv, ok := v.(Context); ok {
		return vv.GetValueE()
	}

	if vv, ok := v.(*Context); ok {
		return vv.GetValueE()
	}

	return v, nil
}

func (c *Context) GetValueE(paths ...string) (interface{}, error) {
	var path string
	if len(paths) == 0 {
		path = ""
	} else {
		path = strings.Join(paths, ".")
	}

	if path == "" || path == "." {
		return valueOfContext(c.data)
	}

	parts := strings.Split(path, ".")
	v := c.data
	for _, part := range parts {
		if part == "" {
			continue
		}

		t, m, l, _, err := parseValue(v)
		if err != nil {
			return nil, err
		}

		if t == TYPE_MAP {
			if part == "length()" {
				v = len(m)
			} else {
				vv, ok := m[part]
				if !ok {
					return nil, errors.New("Path does not exist")
				}
				v = vv
			}
		} else if t == TYPE_LIST {
			if part == "length()" {
				v = len(l)
			} else {
				idx, err := cast.ToIntE(part)
				if err != nil {
					return nil, err
				}

				if idx < 0 || idx >= len(l) {
					return nil, errors.New("Index out of range")
				}

				v = l[idx]
			}
		} else {
			return nil, errors.New("Path does not exist")
		}
	}

	return valueOfContext(v)
}

func (c *Context) GetValue(paths ...string) interface{} {
	v, err := c.GetValueE(paths...)
	if err != nil {
		return nil
	}

	return v
}

func (c *Context) GetMapContextE(paths ...string) (map[string]*Context, error) {
	value, err := c.GetValueE(paths...)
	if err != nil {
		return nil, err
	}

	t, m, _, _, err := parseValue(value)
	if err != nil {
		return nil, err
	}
	if t != TYPE_MAP {
		return nil, errors.New("Not a map")
	}

	result := make(map[string]*Context)
	for k, v := range m {
		result[k] = New(v)
	}

	return result, nil
}

func (c *Context) GetMapContext(paths ...string) map[string]*Context {
	v, err := c.GetMapContextE(paths...)
	if err != nil {
		return make(map[string]*Context)
	}

	return v
}

func (c *Context) GetListContextE(paths ...string) ([]*Context, error) {
	value, err := c.GetValueE(paths...)
	if err != nil {
		return nil, err
	}
	t, _, l, _, err := parseValue(value)
	if err != nil {
		return nil, err
	}
	if t != TYPE_LIST {
		return nil, errors.New("Not a list")
	}

	var result []*Context
	for _, v := range l {
		result = append(result, New(v))
	}

	return result, nil
}

func (c *Context) GetListContext(paths ...string) []*Context {
	v, _ := c.GetListContextE(paths...)

	return v
}

func (c *Context) Keys() []string {
	_, m, _, _, _ := parseValue(c.data)
	if m == nil {
		return nil
	}

	var result []string

	for k, _ := range m {
		result = append(result, k)
	}

	return result
}

func (c *Context) Length() int {
	t, m, l, _, _ := parseValue(c.data)
	if t == TYPE_MAP {
		return len(m)
	} else if t == TYPE_LIST {
		return len(l)
	} else {
		return 0
	}
}

func (c *Context) GetContextWithTypeE(paths ...string) (SimpleType, *Context, error) {
	value, err := c.GetValueE(paths...)
	if err != nil {
		return TYPE_UNKNOWN, nil, err
	}
	t, _, _, _, err := parseValue(value)
	if err != nil {
		return TYPE_UNKNOWN, nil, err
	}

	return t, New(value), nil
}

func (c *Context) GetContextWithType(paths ...string) (SimpleType, *Context) {
	t, v, _ := c.GetContextWithTypeE(paths...)

	return t, v
}

func (c *Context) SetValue(path string, value interface{}) {
	if path == "" || path == "." {
		c.data = value
		return
	}

	parts := strings.Split(path, ".")
	c.data = setValueRecursive(c.data, parts, value)
}

// TODO: Merge
// func (c *Context) Merge(value interface{}, deep bool, ommitEmpty bool) {
// }

func (c *Context) GetContextE(path string) (*Context, error) {
	v, err := c.GetValueE(path)
	if err != nil {
		return nil, err
	}

	return &Context{
		data: v,
	}, nil
}

func (c *Context) GetContext(path string) *Context {
	return &Context{
		data: c.GetValue(path),
	}
}

// Test for empty values
func (c *Context) IsEmpty(paths ...string) bool {
	v, err := c.GetValueE(paths...)
	if err != nil {
		return true
	}

	t, m, l, s, e := parseValue(v)
	if e != nil {
		return true
	}

	if t == TYPE_MAP {
		return len(m) == 0
	}

	if t == TYPE_LIST {
		return len(l) == 0
	}

	if t == TYPE_SCALAR {
		// if t.
		switch tt := s.(type) {
		case bool:
			return !tt
		case string:
			return tt == ""
		case int:
			return tt == 0
		case int8:
			return tt == 0
		case int16:
			return tt == 0
		case int32:
			return tt == 0
		case int64:
			return tt == 0
		case uint:
			return tt == 0
		case uint8:
			return tt == 0
		case uint16:
			return tt == 0
		case uint32:
			return tt == 0
		case uint64:
			return tt == 0
		case float32:
			return tt == 0
		case float64:
			return tt == 0
		}

		return true
	}

	return true
}

func (c *Context) Unmarshal(i interface{}) error {
	b, err := json.Marshal(c.data)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, i)
}

func (c *Context) UnmarshalJSON(jsonData []byte) error {
	var data interface{}
	json.Unmarshal(jsonData, &data)
	c.data = data

	return nil
}

func (c Context) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.data)
}

func (c *Context) Exist(path string) bool {
	_, err := c.GetValueE(path)
	return err == nil
}

func setValueRecursive(data interface{}, path []string, value interface{}) map[string]interface{} {
	current := path[0]
	_, m, _, _, _ := parseValue(data)

	dataMap := m

	if dataMap == nil {
		dataMap = make(map[string]interface{})
	}
	if len(path) == 1 {
		if valueAsContext, ok := value.(*Context); ok {
			if valueAsContext == nil {
				dataMap[current] = nil
			} else {
				dataMap[current] = valueAsContext.GetValue(".")
			}
		} else {
			dataMap[current] = value
		}
	} else {
		nextData, ok := dataMap[current]
		var nextDataMap map[string]interface{}
		if !ok {
			nextData = make(map[string]interface{})
		} else {
			nextDataMap, ok = nextData.(map[string]interface{})
			if !ok {
				nextData = make(map[string]interface{})
			}
		}

		dataMap[current] = setValueRecursive(nextDataMap, path[1:], value)
	}

	return dataMap
}

type SimpleType uint8

const (
	TYPE_UNKNOWN SimpleType = iota
	TYPE_MAP
	TYPE_LIST
	TYPE_SCALAR
)

func parseValue(data interface{}) (SimpleType, map[string]interface{}, []interface{}, interface{}, error) {
	if asContextRef, ok := data.(*Context); ok {
		data = asContextRef.GetValue("")
	}

	if asContext, ok := data.(Context); ok {
		data = asContext.GetValue("")
	}

	ref := reflect.ValueOf(data)
	kind := ref.Kind()

	if kind == reflect.Struct {
		// struct => map
		numField := ref.NumField()
		result := make(map[string]interface{})
		for i := 0; i < numField; i++ {
			fieldName := ref.Type().Field(i).Name
			fieldValue := ref.Field(i).Interface()
			result[fieldName] = fieldValue
		}

		return TYPE_MAP, result, nil, nil, nil
	} else if kind == reflect.Map {
		// map => map
		result := make(map[string]interface{})

		for _, k := range ref.MapKeys() {
			result[k.String()] = ref.MapIndex(k).Interface()
		}

		return TYPE_MAP, result, nil, nil, nil
	} else if kind == reflect.Slice {
		// slice => slice
		length := ref.Len()
		result := make([]interface{}, ref.Len())

		for i := 0; i < length; i++ {
			result[i] = ref.Index(i).Interface()
		}

		return TYPE_LIST, nil, result, nil, nil
	} else {
		// scalar
		return TYPE_SCALAR, nil, nil, data, nil
	}
}
