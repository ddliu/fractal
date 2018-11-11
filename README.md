# Fractal

Fractal is a Go package that makes it easy to work with dynamic and nested data types, with encoding/decoding support.

## Features

- Nested data type
- Dot path support
- JSON encoding/decoding
- Simple template replacement
- Inplace update
- Common data type support: Struct, Map...

## Install

```
go get -u github.com/ddliu/fractal
```

## Usage

Work with struct

```go
data := myStruct {
    Key1: "Value1",
    Key2: anotherStruct {
        Key3: "Value3"
    }
}

// Create context
ctx := fractal.New(data)
println(ctx.String("Key2.Key3"))
// output: Value3
```

Work with json

```go
ctx := fractal.FromJson([]byte(`{"key1": "value1", "key2": {"key3": "value3"}}`))
println(ctx.String("key2.key3"))
```

Work with map

```go
ctx := fractal.New(map[string]interface{
    "key1": "value1",
    "key2": "value2",
})
```

Update:

```go
ctx.SetValue("key2.new_key", 3)
```

Simple template:

```go
tpl := `Author: ${author.name}; License: ${license}; `

c := New(nil)
c.SetValue("author", map[string]string{
    "name":  "Dong",
    "email": "test@example.com",
})

c.SetValue("license", "MIT")

println(c.Tpl(tpl))

// output: Author: Dong; License: MIT; 
```
