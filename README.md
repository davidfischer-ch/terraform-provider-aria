# Terraform Provider Aria (Terraform Plugin Framework)

This is the [Terraform](https://www.terraform.io) provider for VMWare's Aria Automation Platform.

We'll [publish it on the Terraform Registry](https://developer.hashicorp.com/terraform/registry/providers/publishing) so that others can use it.

It has been developped by the CSC Team from the IT department of the State of Geneva (Switzerland).

Please be aware that Broadcom is not responsible neither involved on this project.

_This provider is built on the [Terraform Plugin Framework](https://github.com/hashicorp/terraform-plugin-framework). See [Which SDK Should I Use?](https://developer.hashicorp.com/terraform/plugin/framework-benefits) in the Terraform documentation for additional information._

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.21

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Using the provider

Fill this in for each provider

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate ./...`.
To format the code `make fmt` or `find . -name "*.go" -exec gofmt -s -w {} \;`

You have to set the following environment variables `ARIA_HOST` and `ARIA_REFRESH_TOKEN` before running tests. For the unit tests you can set those to dummy values.

To run the full suite of Unit tests, run `go test ./...`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```shell
make testacc
```
