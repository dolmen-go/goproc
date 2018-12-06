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
	"text/template"

	"github.com/dolmen-go/flagx"
	"github.com/dolmen-go/jsonptr"
	yaml "gopkg.in/yaml.v2"
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
	flag.Var(flagx.Slice(&templates, "", nil), "i", "input template `file`")

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
	if len(templates) == 0 {
		if len(args) < 1 {
			return errors.New("missing input template arguments")
		}
		templates = args[:1]
		args = args[1:]
	}

	tmpl := template.New("")
	tmpl.Funcs(template.FuncMap{
		"jsonptr": func(ptr string, doc interface{}) (interface{}, error) {
			return jsonptr.Get(doc, ptr)
		},
	})

	var err error
	_, err = tmpl.ParseFiles(templates...)
	if err != nil {
		return err
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

	return tmpl.ExecuteTemplate(os.Stdout, filepath.Base(templates[0]), data)
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
	dec.SetStrict(true)
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
