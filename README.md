config
======
[![GoDoc](https://godoc.org/bitbucket.org/utils/config?status.svg)](https://godoc.org/bitbucket.org/utils/config)
[![wercker status](https://app.wercker.com/status/87bdaa59a8345f8395a9ba8e1df6a852/s "wercker status")](https://app.wercker.com/project/bykey/87bdaa59a8345f8395a9ba8e1df6a852)

Load environment specific config files

For more information read the documentation at [godoc.org](https://godoc.org/bitbucket.org/utils/config?status.svg).

Usage:

~~~go
package main

import (
	"fmt"

	"bitbucket.org/utils/config"
)

var C struct {
	Host    string
	ENV     string
	DB      string
	Cache   string
	Websrvr struct {
		ApiURL string `yaml:"api_url"`
		Creds  struct {
			Username string
			Password string
		}
	}
}

func main() {
	//load config
	config.LoadFromFile("./config.yml", &C, nil)
	fmt.Printf("%#v\n", C)
	fmt.Println(C.Websrvr.ApiURL)
	//use it as a normal struct
	//C.Websrvr.ApiURL etc,..
}
~~~

sample `config.yml` file

~~~yaml
---
default: &default
  host: localhost
  env: {{.GOENV}}
  db:  "host=/var/run/postgresql dbname=awesome_{{.GOENV}} sslmode=disable"
  cache: {{.HOME}}/cache
  websrvr:
    api_url: https://sandbox.websrvr.in/1/
    creds:
      username: foobar
      password: awesome

development:
  <<: *default

test:
  <<: *default

production:
  <<: *default
  host: minhajuddin.com/
  websrvr:
    api_url: https://api.websrvr.in/1/
    creds:
      username: {{.DROPBOX_USER}}
      password: {{.DROPBOX_PASSWORD}}
~~~
