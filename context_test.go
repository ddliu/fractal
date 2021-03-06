// Copyright 2018-2019 Liu Dong <ddliuhb@gmail.com>.
// Licensed under the MIT license.

package fractal

import (
	"encoding/json"
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

func TestUnmarshalJSON(t *testing.T) {
	var c Context
	json.Unmarshal([]byte(`{"key": "value"}`), &c)
	if c.String("key") != "value" {
		t.Error()
	}
}

func TestKeys(t *testing.T) {
	var c Context
	json.Unmarshal([]byte(`{"key1": "value1", "key2": "value2"}`), &c)
	keys := c.Keys()

	if !listEquals(keys, []string{"key1", "key2"}) {
		t.Error()
	}
}

func TestList(t *testing.T) {
	var c Context
	json.Unmarshal([]byte(`{"key": ["a", {"b": "b"}]}`), &c)
	if c.String("key.1.b") != "b" {
		t.Error()
	}
	if c.String("key.0") != "a" {
		t.Error()
	}
}

func TestNested(t *testing.T) {
	c := New(nil)
	child := New(nil)
	child.SetValue("a", "b")
	c.SetValue("c", child)

	if c.String("c.a") != "b" {
		t.Error()
	}

	c.SetValue("d", *child)

	if c.String("d.a") != "b" {
		t.Error()
	}

	// x, _ := c.MarshalJSON()
	// print(string(x))
}

func TestLength(t *testing.T) {
	var c Context
	json.Unmarshal([]byte(`{"key": ["a", "b"]}`), &c)
	if c.Int("key.length()") != 2 {
		t.Error()
	}
}

func TestEmpty(t *testing.T) {
	emptyTests := []interface{}{
		nil,
		0,
		0.0,
		false,
		"",
		[]int{},
		map[string]string{},
	}

	noneEmptyTests := []interface{}{
		1,
		1.1,
		true,
		"test",
		[]int{1},
		map[string]string{"test": "test"},
	}

	for _, v := range emptyTests {
		c := New(v)
		if !c.IsEmpty() {
			t.Errorf("%v is not empty", v)
		}
	}

	for _, v := range noneEmptyTests {
		c := New(v)
		if c.IsEmpty() {
			t.Errorf("%v is empty", v)
		}
	}
}

func listEquals(l1, l2 []string) bool {
	if len(l1) != len(l2) {
		return false
	}

	for _, v1 := range l1 {
		exist := false
		for _, v2 := range l2 {
			if v1 == v2 {
				exist = true
			}
		}

		if !exist {
			return false
		}
	}

	return true
}

func TestNil(t *testing.T) {
	c := New(nil)
	if c.GetValue("a") != nil {
		t.Error()
	}
}

func TestGetMap(t *testing.T) {
	c := FromJson([]byte(`{"a":{"k1": "v1", "k2": "v2"}}`))
	if c.GetMapContext("a")["k2"].String() != "v2" {
		t.Error()
	}
}

func TestGetList(t *testing.T) {
	c := FromJson([]byte(`{"a":[{"id": "id1"}, {"id": "id2"}, "simple string"]}`))
	if c.GetListContext("a")[1].String("id") != "id2" {
		t.Error()
	}

	if c.GetListContext("a")[2].String() != "simple string" {
		t.Error()
	}
}

func TestStringNoChild(t *testing.T) {
	c := New("hello")

	if c.String("none.exist") != "" {
		t.Error(c.String("none.exist") + "!=" + "\"\"")
	}
}

func TestNestedInStruct(t *testing.T) {
	type MyStruct struct {
		Name string  `json:"name"`
		Ctx  Context `json:"ctx"`
	}

	data := MyStruct{
		Name: "Hello",
	}

	data.Ctx.SetValue("name1", "value1")

	data1 := MyStruct{}

	b, _ := json.Marshal(data)
	err := json.Unmarshal(b, &data1)
	if err != nil || data1.Ctx.String("name1") != "value1" {
		t.Error("Test tested failed", err, data1)
	}
}

func Test1(t *testing.T) {
	data1 := New(map[string]interface{}{
		"a": "a",
		"b": "b",
	})

	data2 := New(map[string]interface{}{
		"c": "c",
		"d": "d",
	})

	if !data1.IsEmpty("data") {
		t.Error()
	}

	data1.SetValue("data", data2)
	if data1.IsEmpty("data") {
		t.Error()
	}

	if data1.String("data.c") != "c" {
		t.Error()
	}
}
