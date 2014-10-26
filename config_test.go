package config

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var config = struct {
	Host    string
	ENV     string
	DB      string
	Cache   string
	Websrvr struct {
		APIURL string `yaml:"api_url"`
		Creds  struct {
			Username string
			Password string
		}
	}
}{}

func loadConfig(t *testing.T) {
	r := strings.NewReader(`---
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
`)

	Load(r, &config, t.Log)
}

func TestDev(t *testing.T) {
	loadConfig(t)
	assert := assert.New(t)

	//Test the dev environment
	assert.Equal(config.Host, "localhost")
	assert.Equal(config.ENV, "development")
	assert.Equal(config.DB, "host=/var/run/postgresql dbname=awesome_development sslmode=disable")
	assert.Equal(config.Websrvr.APIURL, "https://sandbox.websrvr.in/1/")
	assert.Equal(config.Websrvr.Creds.Username, "foobar")
	assert.Equal(config.Websrvr.Creds.Password, "awesome")

	home := os.Getenv("HOME")
	assert.Equal(config.Cache, home+"/cache")
}

func TestProduction(t *testing.T) {
	assert := assert.New(t)

	//prep the env
	os.Setenv("GOENV", "production")
	os.Setenv("DROPBOX_USER", "supersecretuser")
	os.Setenv("DROPBOX_PASSWORD", "supersecretpassword")

	loadConfig(t)

	assert.Equal(config.Host, "minhajuddin.com/")
	assert.Equal(config.ENV, "production")
	assert.Equal(config.DB, "host=/var/run/postgresql dbname=awesome_production sslmode=disable")
	assert.Equal(config.Websrvr.APIURL, "https://api.websrvr.in/1/")
	assert.Equal(config.Websrvr.Creds.Username, "supersecretuser")
	assert.Equal(config.Websrvr.Creds.Password, "supersecretpassword")

}
