package localhelm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/kimkeehwan/kubeapi/k8sclient"

	"github.com/Masterminds/sprig/v3"
	"github.com/pelletier/go-toml"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

type idpp2Values map[interface{}]interface{}

type Values struct {
	Values interface{}
}

// type files map[string][]byte
type Files struct {
	basedir string
}

type Idpp2ProjectHelm struct {
	template *template.Template
}

func (s *Idpp2ProjectHelm) ReadValues(valuePath string) (idpp2Values, error) {
	fBytes, err := os.ReadFile(valuePath)
	if err != nil {
		log.Printf("%s %s\n", valuePath, err.Error())
		return nil, err
	}
	values := make(idpp2Values)

	err1 := yaml.Unmarshal(fBytes, &values)
	if err1 != nil {
		log.Printf("%s %s\n", valuePath, err.Error())
		return nil, err
	}

	return values, nil
}

func NewIdpp2ProjectHelm(templateFile string, tplFile string) (*Idpp2ProjectHelm, error) {

	bTemp, _ := os.ReadFile(templateFile)
	stemp := string(bTemp)
	bTpl, _ := os.ReadFile(tplFile)
	stpl := string(bTpl)

	rootTemplate := template.New("root")
	funcMap(rootTemplate)
	// basetemplate, err := template.Parse(stemp)
	template, err := rootTemplate.Parse(stemp)
	if err != nil {
		return nil, err
	}

	template.New("gotpl").Parse(stpl)
	// tpltemplate, err := template.New("gotpl").Parse(stpl)
	// if err != nil {
	// 	return nil, err
	// }
	// log.Printf("%v", tpltemplate)

	helm := &Idpp2ProjectHelm{template: rootTemplate}

	return helm, nil
}

func (s *Idpp2ProjectHelm) delimitedYamls(data []byte) ([]k8sclient.ResourceSpecs, error) {
	// results := []map[string]interface{}{}
	results := []k8sclient.ResourceSpecs{}

	defaultYamlDelimiter := []byte("---")
	delimited := bytes.Split(data, defaultYamlDelimiter)
	for _, segment := range delimited {
		m := make(k8sclient.ResourceSpecs)
		err := yaml.Unmarshal([]byte(segment), &m)
		if err != nil {
			return nil, err
		}
		results = append(results, m)
	}

	return results, nil
}

func (s *Idpp2ProjectHelm) Execute(namespace string, valueFile string) ([]k8sclient.ResourceSpecs, error) {
	values, err := s.ReadValues(valueFile)

	if err != nil {
		return nil, err
	}

	var b bytes.Buffer

	values["namespace"] = namespace
	templValues := &Values{Values: values}

	if err := s.template.Execute(io.Writer(&b), templValues); err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	return s.delimitedYamls(b.Bytes())

	// return b.Bytes(), nil
}

func funcMap(t *template.Template) template.FuncMap {
	f := sprig.TxtFuncMap()
	delete(f, "env")
	delete(f, "expandenv")
	// Add some extra functionality
	extra := template.FuncMap{
		"toToml":        toTOML,
		"toYaml":        toYAML,
		"fromYaml":      fromYAML,
		"fromYamlArray": fromYAMLArray,
		"toJson":        toJSON,
		"fromJson":      fromJSON,
		"fromJsonArray": fromJSONArray,
		"Files":         newFiles,

		"include": func(string, interface{}) string { return "not implemented" },
	}

	for k, v := range extra {
		f[k] = v
	}

	includedNames := make(map[string]int)

	f["include"] = func(name string, data interface{}) (string, error) {
		var buf strings.Builder
		if v, ok := includedNames[name]; ok {
			if v > 1000 {
				return "", nil
			}
			includedNames[name]++
		} else {
			includedNames[name] = 1
		}
		err := t.ExecuteTemplate(&buf, name, data)
		includedNames[name]--
		return buf.String(), err
	}

	t.Funcs(f)
	return f
}

// toYAML takes an interface, marshals it to yaml, and returns a string. It will
// always return a string, even on marshal error (empty string).
//
// This is designed to be called from a template.
func toYAML(v interface{}) string {
	data, err := yaml.Marshal(v)
	if err != nil {
		// Swallow errors inside of a template.
		return ""
	}
	return strings.TrimSuffix(string(data), "\n")
}

// fromYAML converts a YAML document into a map[string]interface{}.
//
// This is not a general-purpose YAML parser, and will not parse all valid
// YAML documents. Additionally, because its intended use is within templates
// it tolerates errors. It will insert the returned error message string into
// m["Error"] in the returned map.
func fromYAML(str string) map[string]interface{} {
	m := map[string]interface{}{}

	if err := yaml.Unmarshal([]byte(str), &m); err != nil {
		m["Error"] = err.Error()
	}
	return m
}

// fromYAMLArray converts a YAML array into a []interface{}.
//
// This is not a general-purpose YAML parser, and will not parse all valid
// YAML documents. Additionally, because its intended use is within templates
// it tolerates errors. It will insert the returned error message string as
// the first and only item in the returned array.
func fromYAMLArray(str string) []interface{} {
	a := []interface{}{}

	if err := yaml.Unmarshal([]byte(str), &a); err != nil {
		a = []interface{}{err.Error()}
	}
	return a
}

// toTOML takes an interface, marshals it to toml, and returns a string. It will
// always return a string, even on marshal error (empty string).
//
// This is designed to be called from a template.
func toTOML(v interface{}) string {
	b := bytes.NewBuffer(nil)
	e := toml.NewEncoder(b)
	err := e.Encode(v)
	if err != nil {
		return err.Error()
	}
	return b.String()
}

// toJSON takes an interface, marshals it to json, and returns a string. It will
// always return a string, even on marshal error (empty string).
//
// This is designed to be called from a template.
func toJSON(v interface{}) string {
	data, err := json.Marshal(v)
	if err != nil {
		// Swallow errors inside of a template.
		return ""
	}
	return string(data)
}

// fromJSON converts a JSON document into a map[string]interface{}.
//
// This is not a general-purpose JSON parser, and will not parse all valid
// JSON documents. Additionally, because its intended use is within templates
// it tolerates errors. It will insert the returned error message string into
// m["Error"] in the returned map.
func fromJSON(str string) map[string]interface{} {
	m := make(map[string]interface{})

	if err := json.Unmarshal([]byte(str), &m); err != nil {
		m["Error"] = err.Error()
	}
	return m
}

// fromJSONArray converts a JSON array into a []interface{}.
//
// This is not a general-purpose JSON parser, and will not parse all valid
// JSON documents. Additionally, because its intended use is within templates
// it tolerates errors. It will insert the returned error message string as
// the first and only item in the returned array.
func fromJSONArray(str string) []interface{} {
	a := []interface{}{}

	if err := json.Unmarshal([]byte(str), &a); err != nil {
		a = []interface{}{err.Error()}
	}
	return a
}
func newFiles() Files {
	staticPath := viper.GetString("StaticPath")
	// staticPath := "/workspaces/projectManager/data/static"
	return Files{basedir: staticPath}
	// return newFilesIn(staticPath)
}

// func newFilesIn(basedir string) files {
// 	files := make(map[string][]byte)
// 	// basedir := "/workspaces/projectManager/data/static"
// 	readfiles, err := os.ReadDir(basedir)
// 	if err != nil {
// 		log.Fatalf("readdir %s", err.Error())
// 	}
// 	for _, file := range readfiles {

// 		if !file.IsDir() {
// 			data, err := os.ReadFile(fmt.Sprintf("%v/%v", basedir, file.Name()))
// 			if err != nil {
// 				log.Fatalf("read fail %s", err.Error())
// 			}
// 			files[file.Name()] = data
// 		}
// 	}

// 	return files
// }

// func (f files) GetBytes(name string) []byte {

// 	data, err := os.ReadFile(fmt.Sprintf("%v/%v", f.basedir, name))
// 	if v, ok := f[name]; ok {
// 		return v
// 	}
// 	return []byte{}
// }

func (f Files) Get(name string) string {
	path := fmt.Sprintf("%v/%v", f.basedir, name)
	data, err := os.ReadFile(fmt.Sprintf("%v/%v", f.basedir, name))
	if err != nil {
		log.Fatalf("read fail %s", path)
	}

	return string(string(data))
}

func (f Files) GetConfig(name string) string {
	data, err := os.ReadFile(fmt.Sprintf("%v/%v", f.basedir, name))
	if err != nil {
		log.Fatalf("read fail")
	}

	m := make(map[string]string)
	m[name] = string(data)
	return toYAML(m)
}

// func (f files) AsConfig() string {
// 	if f == nil {
// 		return ""
// 	}

// 	m := make(map[string]string)

// 	// Explicitly convert to strings, and file names
// 	for k, v := range f {
// 		m[path.Base(k)] = string(v)
// 	}

// 	return toYAML(m)
// }

// func (f files) AsSecrets() string {
// 	if f == nil {
// 		return ""
// 	}

// 	m := make(map[string]string)

// 	for k, v := range f {
// 		m[path.Base(k)] = base64.StdEncoding.EncodeToString(v)
// 	}

// 	return toYAML(m)
// }

// func (f files) Lines(path string) []string {
// 	if f == nil || f[path] == nil {
// 		return []string{}
// 	}

// 	return strings.Split(string(f[path]), "\n")
// }
