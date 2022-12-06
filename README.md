# pets-api

[![CI](https://github.com/tj-actions/coverage-badge-go/workflows/CI/badge.svg)](https://github.com/tj-actions/coverage-badge-go/actions?query=workflow%3ACI)
![Coverage](https://img.shields.io/badge/Coverage-80.6%25-brightgreen)
[![Update release version.](https://github.com/tj-actions/coverage-badge-go/workflows/Update%20release%20version./badge.svg)](https://github.com/tj-actions/coverage-badge-go/actions?query=workflow%3A%22Update+release+version.%22)

## Table of contents

* [Layers](#layers)
* [Command](#command)
* [Project setup](#project-setup)
    * [Required](#required)
    * [Optional](#optional)
    * [Script](#script)
* [Tasks](#tasks)
    * [Test](#test)
    * [Build](#build)
* [Configuration](#configuration)
* [RestClient](#restclient)
* [Environment](#environment)
    * [Local services](#local-services)
    * [Docker compose](#docker-compose)

## Layers

![layers.png](src/resources/images/layers.png)

## Command

[![asciicast](https://asciinema.org/a/jI7h9PfaAnZBO31Pj9ysD0Ovk.svg)](https://asciinema.org/a/jI7h9PfaAnZBO31Pj9ysD0Ovk)

## Project setup

Install the following dependencies

### Required

- [Golang Lint](https://golangci-lint.run/)
- [Golang Task](https://taskfile.dev/)
- [Golang Dependencies Update](https://github.com/oligot/go-mod-upgrade)
- [ent - An Entity Framework For Go](https://github.com/ent/ent)
- [MySQL](https://www.mysql.com/)

### Optional

If you want to browse api with SSL self-signed certificate

- [NGINX](https://www.nginx.com/)
- [mkcert](https://github.com/FiloSottile/mkcert)

### Script

Some tools are required, go-task is a task runner similar to Gradle, Gulp, NPM, etc.

```shell
brew install go-task/tap/go-task golangci-lint mysql
go install github.com/oligot/go-mod-upgrade@latest
go install entgo.io/ent/cmd/ent@latest
```

Optional if you want to have a local self-signed certificate.

```shell
brew install mkcert nginx
```

## Tasks

### Test

```shell
task test
```

### Build

```shell
task build
```

## Configuration

Environment configuration is based on **Archaius Config**, you should use a similar folder structure.
*SCOPE* env variable in remote environment is required

```
└── config
    ├── config.yml (shared config)
    └── dev
        └── config.yml (for local development)
    └── prod (for remote environment)
        └── config.yml (base config)
        └── {environment}.config.yml (base config)
```

The SDK provides a simple configuration hierarchy

* resources/config/config.properties (shared config)
* resources/config/{environment}/config.properties (override shared config by environment)
* resources/config/{environment}/{scope}.config.properties (override env and shared config by scope)

example *test.pets-api.internal.com*

```
└── config
    ├── config.yml                              3th (third)
    └── dev
        └── config.yml                          <ignored>
    └── prod
        └── config.yml (base config)            2nd (second)
        └── test.config.yml (base config)       1st (first)
```

* 1st (first)   prod/test.config.yml
* 2nd (second)  prod/config.yml
* 3th (third)   config.yml

```
2022/11/20 13:24:26 INFO: Two files have same priority. keeping
    /resources/config/prod/test.config.yml value
2022/11/20 13:24:26 INFO: Configuration files:
    /resources/config/prod/test.config.yml,
    /resources/config/prod/config.yml,
    /resources/config/config.yml
2022/11/20 13:24:26 INFO: invoke dynamic handler:FileSource
2022/11/20 13:24:26 INFO: enable env source
2022/11/20 13:24:26 INFO: invoke dynamic handler:EnvironmentSource
2022/11/20 13:24:26 INFO: archaius init success
2022/11/20 13:24:26 INFO: ENV: prod, SCOPE: test
2022/11/20 13:24:26 INFO: create new watcher
2022/11/20 13:24:26 Listening on port 8080
2022/11/20 13:24:26 Open http://127.0.0.1:8080/ping in the browser
```

## RestClient

RestPool and RestClient are dynamically created by config files.

```yaml
# gitlab
rest:
  pool:
    default:
      pool:
        size: 20
        timeout: 2000
        connection-timeout: 5000
  client:
    gitlab:
      pool: default
```

```go
package main

import (
	"github.com/src/main/app/clients/gitlab"
	"github.com/src/main/app/config"
	"log"
)

func main() {
    factory := config.ProvideRestClients()
	gitLabClient := gitlab.NewGitLabClient(factory.Get("gitlab"), nil)
	result, err := gitLabClient.GetGroups()
	if err != nil {
	    log.Fatal(err)
    }
	log.Println(result)
}

```

## Environment

Set up your local environment. For now, you must set some env values as environment variable (export command).

    * PROD_CONNECTION_STRING: {mysql_connection_string}
    * GITLAB_TOKEN: {your_access_token}

Also, you can set the same values inside resources/config/dev/config.yml

### Local services

If you want to run program.go you must install mysql in your environment.

### Docker compose

If you want to run isolated local environment you must install docker and run docker-compose.
