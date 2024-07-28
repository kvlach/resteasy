# REST easy

Easy-to-use REST API wrapper.

## Why

This library is **not** meant for use in production.
It is designed to facilitate the creation of quick "recipes,"
resembling scripting languages more than traditional Go.
Accordingly, none of the functions return errors, they instead panic.

If JSON responses are received, the standard library's JSON parser
has been modified to enable *strict* data definitions.
The library will panic if any of the following are true:
- The struct contains a field without a `json` tag
- The response contains a key not present in the struct
- The struct contains a field not present in the response

Although it might seem counter-intuitive, Go was chosen for the following reasons:
- Strong/static typing
- Trivial parallelism
- Precise data definitions with structs:
  - Ensures that responses contain all expected data
  - Autocomplete integration, reducing the need to frequently refer to documentation.
- `go run` makes it easy and interactive to use

## Quickstart

```go
package main

import (
	rz "github.com/kvlach/resteasy"
)

type Example struct {
	Field1       string `json:"field1"`
	NestedStruct []struct {
		Field2 int `json:"field2"`
	} `json:"nested_struct"`
}

func main() {
	var resp Example

	rz.GET("https://example.com").
		Query("param1", "value1", "param2", "value2").
		Retry(5).
		JSON(true).
		Do(&resp)
}
```
