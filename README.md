# pets-api

[![CI](https://github.com/tj-actions/coverage-badge-go/workflows/CI/badge.svg)](https://github.com/tj-actions/coverage-badge-go/actions?query=workflow%3ACI)
![Coverage](https://img.shields.io/badge/Coverage-76.8%25-brightgreen)
[![Update release version.](https://github.com/tj-actions/coverage-badge-go/workflows/Update%20release%20version./badge.svg)](https://github.com/tj-actions/coverage-badge-go/actions?query=workflow%3A%22Update+release+version.%22)

## Table of contents
* [Project setup](#project-setup)
    * [Required](#required)
    * [Optional](#optional)
    * [Script](#script)
* [Tasks](#task)
    * [Test](#test)
    * [Build](#build)
* [Configuration](#configuration)
* [Environment](#enviroment)

## Project setup

Install the following dependencies

### Required

- [Golang Lint](https://golangci-lint.run/)
- [Golang Task](https://taskfile.dev/)
- [Golang Dependencies Update](https://github.com/oligot/go-mod-upgrade)
- [ent - An Entity Framework For Go](https://github.com/ent/ent)

### Optional

If you want to browse api with SSL self-signed certificate

- [NGINX](https://www.nginx.com/)
- [mkcert](https://github.com/FiloSottile/mkcert)

### Script

```shell
brew install go-task/tap/go-task
brew install golangci-lint
go install github.com/oligot/go-mod-upgrade@latest
go install entgo.io/ent/cmd/ent@latest
brew install mkcert
brew install nginx
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

Environment configuration is based on Archaius Config, you should use a similar folder structure.
SCOPE env variable in remote environment is required

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

## Environment

SECRETS_STORE_GITLAB_TOKEN_KEY_NAME={gitlab_access_token}
SECRETS_STORE_PROD_CONNECTION_STRING_KEY_NAME={mysql_connection_string}