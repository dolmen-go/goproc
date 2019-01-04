# goproc

Apply [Go templates](https://golang.org/pkg/text/template/#hdr-Text_and_spaces) to JSON or YAML data.

## Build

    go build

## Usage

    goproc <Go-template-file> <JSON-or-YAML-file>

## Template syntax

See the [Go templates](https://golang.org/pkg/text/template/#hdr-Text_and_spaces) documentation.

The [Hugo documentation on Go templates](https://gohugo.io/templates/introduction/) may also be useful for a friendlier approach, but note that it contains references to features unique to Hugo (ex: `partial`).

## Functions extensions

The following functions are available in addition to the [standard functions](https://golang.org/pkg/text/template/#hdr-Functions).

### `jsonptr`

To ease the extraction of data, `jsonptr` allows to express data location using
JSON Pointer ([RFC 6901](https://tools.ietf.org/html/rfc6901)).

Usages:

    {{ jsonptr "pointer" . }}
    {{ . | jsonptr "pointer" }}

Examples:

1. `goproc` [`testdata/02.gotmpl`](testdata/02.gotmpl) [`testdata/02.json`](testdata/02.json)
2. `goproc` [`testdata/03.gotmpl`](testdata/03.gotmpl) [`testdata/03.json`](testdata/03.json)

## See also

https://github.com/naotookuda/go-template-cli
