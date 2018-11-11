// Copyright 2018 Liu Dong <ddliuhb@gmail.com>.
// Licensed under the MIT license.

package fractal

import (
	"testing"
)

func TestContext(t *testing.T) {
	c := New(nil)
	c.SetValue("a.b.c.d", 99)
	if c.Int("a.b.c.d") != 99 {
		t.Error("Get value error")
	}

	c.SetValue("a1.b.c.d", map[string]string{
		"k1": "v1",
		"k2": "v2",
	})

	if c.String("a1.b.c.d.k1") != "v1" {
		t.Error("Get value error")
	}
}

func TestJson(t *testing.T) {
	c := FromJson([]byte(`{"key": "value"}`))
	if c.String("key") != "value" {
		t.Error()
	}
}

type testStruct1 struct {
	Key1 string
	Key2 testStruct2
}

type testStruct2 struct {
	Key3 string
}

func TestStruct(t *testing.T) {
	v := testStruct1{
		Key1: "Value1",
		Key2: testStruct2{
			Key3: "Value3",
		},
	}
	c := New(v)
	if c.String("Key2.Key3") != "Value3" {
		t.Error()
	}
}
