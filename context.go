// Copyright 2018 Liu Dong <ddliuhb@gmail.com>.
// Licensed under the MIT license.

package fractal

import (
	"encoding/json"
	"errors"
	"github.com/spf13/cast"
	"reflect"
	"strings"
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

func (c *Context) GetValueE(paths ...string) (interface{}, error) {
	var path string
	if len(paths) == 0 {
		path = ""
	} else {
		path = strings.Join(paths, ".")
	}

	if path == "" || path == "." {
		return c.data, nil
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
		}
	}

	return v, nil
}

func (c *Context) GetValue(paths ...string) interface{} {
	v, err := c.GetValueE(paths...)
	if err != nil {
		return nil
	}

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

func (c *Context) SetValue(path string, value interface{}) {
	if path == "" || path == "." {
		c.data = value
		return
	}

	parts := strings.Split(path, ".")
	c.data = setValueRecursive(c.data, parts, value)
}

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

func (c *Context) MarshalJSON() ([]byte, error) {
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
