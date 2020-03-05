# sts: struct to struct: generator of transformation functions

[![codecov](https://codecov.io/gh/ekhabarov/sts/branch/master/graph/badge.svg)](https://codecov.io/gh/ekhabarov/sts)
[![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/ekhabarov/sts)](https://github.com/ekhabarov/sts/releases)
[![Travis (.org)](https://img.shields.io/travis/ekhabarov/sts)](https://travis-ci.org/ekhabarov/sts)
[![GoDoc](https://godoc.org/https://godoc.org/github.com/ekhabarov/sts?status.svg)](https://godoc.org/github.com/ekhabarov/sts)
[![Go Report Card](https://goreportcard.com/badge/github.com/ekhabarov/sts)](https://goreportcard.com/report/github.com/ekhabarov/sts)

<!-- vim-markdown-toc GFM -->

* [Install](#install)
* [Motivation](#motivation)
* [Idea](#idea)
  * [Other implementations.](#other-implementations)
* [How](#how)
  * [Step 1](#step-1)
  * [Step 2](#step-2)
    * [Example matcher](#example-matcher)
      * [Int2Bool, NullsTime2TimeTimePtr wait, what?](#int2bool-nullstime2timetimeptr-wait-what)
  * [Step 3](#step-3)
* [go generate](#go-generate)
* [License](#license)

<!-- vim-markdown-toc -->

## Install

```shell
go get -u github.com/ekhabarov/sts/cmd/sts
```

## Motivation
Working on integration between one app and different APIs (most of them,
fortunately, have Go clients) includes pretty much code which transforms one
structure into another, because for Go two structures with identical field set
and identical types are different types. Identical types could be converted one
into another with simple conversion: `targetType(destType)`, but having
[identical](https://golang.org/ref/spec#Type_identity) type is too rare case.

That means it's necessary to write such transformations manually, which is, from
one hand is tediously from another one is straightforward.

## Idea
The idea is as simple as possible: produce set of functions which allow convert
one type into another.

It can be done within three steps:

1. Source code analyze.
1. Field type matching.
1. Generations pair of functions: forward `SourceType2DestType` and reverse `DestType2SourceType`.

### Other implementations.
There is a [plugin](http://github.com/bold-commerce/protoc-gen-struct-transformer) for Protobuf with the same idea.

## How

### Step 1
On first step `sts` have to obtain information about structures which will be
involved into transformation process by analyzing source code files contained
these structures. To achieve this, packages [go/ast](https://golang.org/pkg/go/ast), [go/types](https://golang.org/pkg/go/types), etc., from
standard library can be used.

Using these packages `sts` builds a map with data types information. For details
see [parser.go](./parser.go).

### Step 2
Information from previous step is passes to matcher. Matcher lookups two
structures by name (structures names are passed via CLI params, see examples
below), source (left) and destination (right). Then it builds field pairs using
next rules:

* field on the left structure with `sts` tag will be matched with field on right side by right-side field name equals to `sts` tag value.
* if right-side field not found by name, then `sts` tag value will be compared with value of provided tag list.
* any fields without `sts` or other source tags will be skipped.

#### Example matcher
Let's say we have two structures

```go
type Source struct {
	I  int
	S  string
	I1 int        `sts:"I64"`
	I2 int        `sts:"B"`
	PT *time.Time `sts:"Nt"`
	JJ string     `sts:"json_field"`
	D  int32      `sts:"db_field"`
}
```

and

```go
type Dest struct {
	I         int
	S         string
	I64       int64
	B         bool
	Nt        nulls.Time
	JsonField string `json:"json_field"`
	DB        int64  `db:"db_field"`
}
```

after run a command

```shell
sts -src /path/to/src.go:Source -dst /path/to/dst.go:Dest -o ./output -dt json,db
```

matcher consider next combinations


Source | Destination | Conversion              | Note
-------|-------------|-------------------------|------
 `I`   | `--`        | `--`                    | source field has not tag
 `S`   | `--`        | `--`                    | source field has not tag
 `I1`  | `I64`       | direct                  | matched `sts` tag value and field name
 `I2`  | `B`         | `Int2Bool`              | matched `sts` tag value and field name
 `PT`  | `Nt`        | `NullsTime2TimeTimePtr` | matched `sts` tag value and field name
`JJ`   | `JsonField` | none                    | matched `sts` tag value and `json` tag value. `json` tag passed via `-dt` CLI parameter.
`DB`   | `D`         | direct                  | matched `sts` tag value and `db` tag value. `db` tag passed via `-dt` CLI parameter.


##### Int2Bool, NullsTime2TimeTimePtr wait, what?
Matcher uses type info provided by `go/types` package. When it compares field it
also checks paired field for [assignability](https://golang.org/pkg/go/types/#AssignableTo) and [convertibility](https://golang.org/pkg/go/types/#ConvertibleTo).
* Assignability shows can one field be assigned to another without any conversion.
* Convertibility shows can one field be directly converted to another one.

But in cases when fields in pair are not `assignable` and are not `convertable`,
the tool just generate conversion function with name of format

```go
<SourceType>2<DestType>
// and
<DestType>2<SourceType>
```

that means it's necessary to write these helper functions manually. Fortunately,
quantity of such function should be low. Number of examples can be found in
[examples](./examples/output/helpers.go) package.

### Step 3
On the last step `sts` creates a file with name `<source>_to_<dest>.go` with
pair of ready-to-use functions for each pair of structures passed as a
parameters to `sts`.

```go
// source_to_dest.go

// Auto-generated code. DO NOT EDIT!!!
// Generated by sts v0.0.1-alpha.

package output

import "github.com/ekhabarov/sts/examples"

func Source2Dest(src examples.Source) examples.Dest {
	return examples.Dest{
		I64:       int64(src.I1),
		B:         Int2Bool(src.I2),
		Nt:        TimeTimePtr2NullsTime(src.PT),
		JsonField: src.JJ,
		DB:        int64(src.D),
	}
}
func Dest2Source(src examples.Dest) examples.Source {
	return examples.Source{
		I1: int(src.I64),
		I2: Bool2Int(src.B),
		PT: NullsTime2TimeTimePtr(src.Nt),
		JJ: src.JsonField,
		D:  int32(src.DB),
	}
}
```

## go generate
Go has a command `go generate` ([blog](https://blog.golang.org/generate)|[proposal](https://docs.google.com/document/d/1V03LUfjSADDooDMhe-_K59EgpTEm3V8uvQRuNMAEnjg/edit)).
This command allows to run tools mentioned in special comments in Go code, like
this:

```go
//go:generate sts -src $GOFILE:Source -dst $GOFILE:Dest -o ./output -dt json,db
type Source struct {
	I  int
...
```

after `go generate ./...` will be run, it in turn, will run `sts` tool with
given parameters. `$GOFILE` variable will be replaced with a path to current
`.go` file by `go generate` tool.


## License

MIT License

Copyright (c) 2020 Evgeny Khabarov

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

