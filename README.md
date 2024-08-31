# RIPE DB Provider for Terraform

The [RIPE DB provider for Terraform](https://registry.terraform.io/providers/frederic-arr/ripe/latest) is a plugin that exposes the [RIPE Database](https://apps.db.ripe.net/) to Terraform.


## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.8
- [Go](https://golang.org/doc/install) >= 1.21

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```shell
make testacc
```

If you want to use the development version of the provider, you can configure the `~/.terraformrc` file to use the local provider:

```hcl
provider_installation {
    dev_overrides {
        "frederic-arr/ripe" = "/home/user/go/bin"
    }

    direct {}
}
```
