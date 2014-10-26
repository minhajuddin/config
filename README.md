config
======
[![GoDoc](https://godoc.org/bitbucket.org/utils/config?status.svg)](https://godoc.org/bitbucket.org/utils/config)
[![wercker status](https://app.wercker.com/status/87bdaa59a8345f8395a9ba8e1df6a852/s "wercker status")](https://app.wercker.com/project/bykey/87bdaa59a8345f8395a9ba8e1df6a852)

Load environment specific config files

sample `config.yml` file

For more information read the documentation at [godoc.org](https://godoc.org/bitbucket.org/utils/config?status.svg).

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
