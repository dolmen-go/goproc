# goproc

Apply [Go templates](https://golang.org/pkg/text/template/#hdr-Text_and_spaces) to JSON or YAML data.

## Build

    go build

## Usage

    goproc <Go-template-file> <JSON-or-YAML-file>

## Template syntax

See the [Go templates](https://golang.org/pkg/text/template/#hdr-Text_and_spaces) documentation.

## Functions extensions

The following functions are added in addition to the standard functions.

### `jsonptr`

To ease the extraction of data, `jsonptr` allows to express data location using
JSON Pointer ([RFC 6901](https://tools.ietf.org/html/rfc6901)).

Usages:

    {{ jsonptr "pointer" . }}
    {{ . | jsonptr "pointer" }}

Examples:

1. [->template](testdata/02.gotml) [->data](testdata/02.json)
2. [->template](testdata/03.gotml) [->data](testdata/03.json)
