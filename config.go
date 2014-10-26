//Allows you to use yaml config files for configuration
//You can even use ENVIRONMENT VARS in your config
//e.g. ./config.yml
//    development:
//      database: websrvr_{{.GOENV}}
//
//    test:
//      database: websrvr_{{.GOENV}}
package config

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strings"
	"text/template"

	"gopkg.in/v1/yaml"
)

var (
	GOENV string
)

func getEnv() map[string]string {

	//prep the env
	env := make(map[string]string)

	//parse env vars
	for _, e := range os.Environ() {
		idx := strings.Index(e, "=")
		//if env var is malformed
		if idx < 0 {
			continue
		}
		env[e[:idx]] = e[idx+1:]
	}

	//set default env var
	var ok bool
	GOENV, ok = env["GOENV"]
	if !ok {
		env["GOENV"] = DEFAULTENV
		GOENV = DEFAULTENV
	}

	return env
}

//This env will be loaded if GOENV is not set
var DEFAULTENV = "development"

//Pass the address of the config variable
//e.g. config.Load(r, &config, nil)
func Load(r io.Reader, config interface{}, logFunc func(args ...interface{})) {
	if logFunc == nil {
		logFunc = log.Println
	}

	configBytes, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}
	logFunc("CONFIG: ", string(configBytes))

	tpl := template.Must(template.New("config").Parse(string(configBytes)))

	//pass env to config
	var b bytes.Buffer
	tpl.Execute(&b, getEnv())

	kt := reflect.TypeOf("")
	vt := reflect.TypeOf(config)
	m := reflect.MakeMap(reflect.MapOf(kt, vt))

	configData := m.Interface()

	yaml.Unmarshal(b.Bytes(), configData)

	c := m.MapIndex(reflect.ValueOf(GOENV))

	cptr := reflect.ValueOf(config)

	el := cptr.Elem()
	if el.CanSet() {
		el.Set(c.Elem())
	} else {
		logFunc("ERROR: the config variable should pass the address the config struct")
	}
}

//Loads the config file from the `path` e.g. "./config.yml
//config is a pointer to your config type
//logfunc is the function used to log errors if any,
//you can pass nil if you don't care about the log function
func LoadFromFile(path string, config interface{}, logFunc func(args ...interface{})) {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	Load(f, config, logFunc)
}
