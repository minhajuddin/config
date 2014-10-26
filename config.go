//Package config allows you to use an environment based
//configuration file with ability to interpolate environment
//variables inside the configuration file.
//e.g. ./config.yml
//    default: &default
//      database: awesome_{{.GOENV}}
//      mongo:
//        host: localhost:3030
//        username: foob
//        password: boof
//    development:
//      <<: *default
//
//    production:
//      <<: *default
//      mongo:
//        host: {{.MONGO_URL}}
//        username: {{.MONGO_USER}}
//        password: {{.MONGO_PASSWORD}}
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
	//GOENV holds the current value if it is invoked with it
	//e.g. GOENV=test go run * will set GOENV to test
	//if the app is invoked without any GOENV it is set to the DEFAULTENV
	GOENV string
	//DEFAULTENV is set to development by default but you can change this
	//before calling config.Load to change the default GOENV value
	DEFAULTENV = "development"
)

//LoadFromFile allows you to read the configuration from a file on disk
//path is the path to the config file, e.g. "./config.yml
//config is a pointer to your config variable
//logfunc is the function used to log errors if any, you can pass nil if you don't care about the log function
func LoadFromFile(path string, config interface{}, logFunc func(args ...interface{})) {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	Load(f, config, logFunc)
}

//Load allows you to read the configuration from an io.Reader instead of a file
//Check the LoadFromFile function for more information on arguments
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
