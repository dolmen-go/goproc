package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/dolmen-go/flagx"
	"github.com/dolmen-go/jsonptr"
	yaml "gopkg.in/yaml.v3"
)

func main() {
	log.SetPrefix("")
	log.SetFlags(0)

	if err := _main(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func _main() error {

	var templates []string
	var inlineTemplates []string
	var withEnv bool
	flag.Var(flagx.Slice(&templates, "", nil), "i", "input template `file`")
	flag.Var(flagx.Slice(&inlineTemplates, "", nil), "e", "inline template")
	flag.BoolVar(&withEnv, "env", false, "enable env function")

	loadData := loadJSON
	flag.Var(flagx.BoolFunc(func(stdinAsYAML bool) error {
		if stdinAsYAML {
			loadData = loadYAML
		} else {
			loadData = loadJSON
		}
		return nil
	}), "yaml", "load data from stdin as YAML (default is JSON)")

	flag.Parse()

	// TODO handle -version

	args := flag.Args()
	if len(inlineTemplates)+len(templates) == 0 {
		if len(args) < 1 {
			return errors.New("missing input template arguments")
		}
		templates = args[:1]
		args = args[1:]
	}

	tmpl := template.New("")
	funcs := template.FuncMap{
		"error": func(msg string) string {
			panic(errors.New(msg))
		},
		"jsonptr": func(ptr string, doc interface{}) (interface{}, error) {
			return jsonptr.Get(doc, ptr)
		},
		"json": func(doc interface{}) (string, error) {
			b, err := json.Marshal(doc)
			return string(b), err
		},
		"yaml": func(doc interface{}) (string, error) {
			b, err := yaml.Marshal(doc)
			return string(b), err
		},
	}

	if withEnv {
		funcs["env"] = func(names ...string) interface{} {
			switch len(names) {
			case 0:
				envNative := os.Environ()
				env := make(map[string]string, len(envNative))
				for _, v := range envNative {
					p := strings.IndexByte(v, '=')
					if p == -1 {
						continue
					}
					env[v[:p]] = v[p+1:]
				}
				return env
			case 1:
				return os.Getenv(names[0])
			default:
				envNative := os.Environ()
				env := make(map[string]string, len(envNative))
				for _, v := range envNative {
					p := strings.IndexByte(v, '=')
					if p == -1 {
						continue
					}
					n := v[:p]
					for _, name := range names {
						if name == n {
							env[name] = v[p+1:]
						}
					}
				}
				return env
			}
		}
	}

	tmpl.Funcs(funcs)

	var err error

	for _, s := range inlineTemplates {
		_, err := tmpl.Parse(s)
		if err != nil {
			return err
		}
	}

	var templateName string
	if len(templates) > 0 {
		_, err = tmpl.ParseFiles(templates...)
		if err != nil {
			return err
		}
		templateName = filepath.Base(templates[0])
	} else {
		templateName = tmpl.Name()
	}

	var data interface{}
	if len(args) > 0 {
		data, err = loadFile(args[0])
	} else {
		data, err = loadData(os.Stdin)
	}
	if err != nil {
		return err
	}

	return tmpl.ExecuteTemplate(os.Stdout, templateName, data)
}

func loadFile(pth string) (interface{}, error) {
	f, err := os.Open(pth)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	switch filepath.Ext(pth) {
	case ".json":
		return loadJSON(f)
	case ".yaml", ".yml":
		return loadYAML(f)
	default:
		return nil, errors.New("unsupported file extension")
	}
}

func loadJSON(r io.Reader) (data interface{}, err error) {
	err = json.NewDecoder(r).Decode(&data)
	if err != nil {
		return nil, err
	}
	return data, err
}

func loadYAML(r io.Reader) (data interface{}, err error) {
	dec := yaml.NewDecoder(r)
	err = dec.Decode(&data)
	if err != nil {
		return nil, err
	}
	return fixMaps(data), err
}

func fixMaps(v interface{}) interface{} {
	switch v := v.(type) {
	case nil, bool, string, int, int64, float64:
	case []interface{}:
		for i, item := range v {
			v[i] = fixMaps(item)
		}
	case map[interface{}]interface{}:
		m := make(map[string]interface{}, len(v))
		for key, val := range v {
			m[fmt.Sprint(key)] = fixMaps(val)
		}
		return m
	case map[string]interface{}:
		for key, value := range v {
			v[key] = fixMaps(value)
		}
	}
	return v
}
