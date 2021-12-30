[![Tests](https://github.com/adwitiyaio/arka/actions/workflows/ci.yml/badge.svg)](https://github.com/adwitiyaio/arka/actions/workflows/ci.yml)

## Preface

This is the next generation backend platform similar to parse built with the following technologies and tools amongst others:
- [Golang](https://golang.org/) - Programming Language
- [Postgres](https://www.postgresql.org/) - Database
- [Redis](https://redis.io/) - Used as a cache
- [Gorm](https://gorm.io/) - ORM framework
- [AWS](https://aws.amazon.com/) - Cloud provider for storing assets amongst others

## Installing

- Run `go get github.com/adwitiyaio/arka` to get the latest version.
- Run `go get github.com/adwitiyaio/arka@v0.0.1` to get a specific version.
- Run `go get -u github.com/adwitiyaio/arka/...` to update to the latest version.

## Setup

Create a copy of `sample.env` as `test.env` and fill in the appropriate values to run tests

## Testing

### Command Line (Coverage)

```shell
go test -v ./...
```

### Command Line (Coverage)

```shell
go test -v `go list ./...` -coverprofile=coverage.out -covermode=atomic
```

### Visual Studio Code

- You can use the built in `Test Explorer` panel to run tests. 

### JetBrains IDEs

- You can use the [run configuration file](.run/Tests.run.xml) to run tests

## Developing

See [Code Structure](./code-structure.md) for more details
