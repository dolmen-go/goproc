# goproc

Apply [Go templates](https://golang.org/pkg/text/template/#hdr-Text_and_spaces) to JSON or YAML data.

## Install

    go install github.com/dolmen-go/goproc@latest

## Usage

    goproc [ -env | -env=true | -env=VAR ] -e <Go-template-text> [ <JSON-or-YAML-file> ]
    goproc [ -env | -env=true | -env=VAR ] -i <Go-template-file> [ <JSON-or-YAML-file> ]

The default input is STDIN in JSON format. Use `-yaml` flag to handle STDIN as YAML.

When an input file is given, the file extension determines if it is parsed as JSON or YAML.

In any case there is no magic detection (that usually lead to security issues).

## Template syntax

See the [Go templates](https://golang.org/pkg/text/template/#hdr-Text_and_spaces) documentation.

The [Hugo documentation on Go templates](https://gohugo.io/templates/introduction/) may also be useful for a friendlier approach, but note that it contains references to features unique to Hugo (ex: `partial`).

## Functions extensions

The following functions are available in addition to the [standard functions](https://golang.org/pkg/text/template/#hdr-Functions).

### `env`

Allow to use environment variables (with security proctections).

This function must be explicitely enabled using the `-env` flag:

* `-env`: enables the `env` function. Any environment variable can be used, but listing variables is blocked.
* `-env=`: enables the `env` function, but whitelist of allowed variables is cleared (no variables allowed).
* `-env=name1,name2`: enables the `env` function. Only the variables `name1` and `name2` are visible in calls to `env`.

Usage:

    {{ env "HOME" }}     {{- /* Get value of HOME environment variable */}}
    {{ "HOME" | env }}   {{- /* Get value of HOME environment variable */}}
    {{ range $name, $value := env "HOME" "LANG" -}}
      {{ $name }}={{ $value }}
    {{ end }}


### `error`

Usage:

    {{ error "message" }}

Example:

    echo 0 | goproc -e '{{ error "fail!" }}'
    template: :1:3: executing "" at <error "fail!">: error calling error: fail!

### `jsonptr`

To ease the extraction of data, `jsonptr` allows to express data location using
JSON Pointer ([RFC 6901](https://tools.ietf.org/html/rfc6901)).

Usages:

    {{ jsonptr "pointer" . }}
    {{ . | jsonptr "pointer" }}


Examples:

1. `goproc` [`testdata/02.gotmpl`](testdata/02.gotmpl) [`testdata/02.json`](testdata/02.json)
2. `goproc` [`testdata/03.gotmpl`](testdata/03.gotmpl) [`testdata/03.json`](testdata/03.json)

### `json`

Convert input to JSON.

Usage:

    {{ json }}

Example:

    echo '{"data": ["x"]}' | goproc -e '{{.data | json}}{{print "\n"}}'
    ["x"]

### `yaml`

Convert input to YAML.

Usage:

    {{ yaml }}

Example:

    echo '{"data": ["x"]}' | goproc -e '{{.data | yaml}}'
    - x

## See also

https://github.com/naotookuda/go-template-cli

## License

Copyright 2018-2024 Olivier Mengué

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
