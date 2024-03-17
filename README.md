# Terraform provider for couchbase
![GitHub release (with filter)](https://img.shields.io/github/v/release/lukasbudisky/terraform-provider-couchbase?style=flat-square&logo=terraform&logoColor=blue&label=latest%20version&labelColor=grey&link=https%3A%2F%2Fregistry.terraform.io%2Fproviders%2Flukasbudisky%2Fcouchbase%2Flatest&link=https%3A%2F%2Fregistry.terraform.io%2Fproviders%2Flukasbudisky%2Fcouchbase%2Flatest)
![GitHub Workflow Status (with event)](https://img.shields.io/github/actions/workflow/status/lukasbudisky/terraform-provider-couchbase/.github%2Fworkflows%2Fmain_branch.yml?style=flat-square&logo=github&logoColor=white&label=tests&labelColor=grey&link=https%3A%2F%2Fgithub.com%2Flukasbudisky%2Fterraform-provider-couchbase%2Freleases&link=https%3A%2F%2Fgithub.com%2Flukasbudisky%2Fterraform-provider-couchbase%2Freleases)
![GitHub all releases](https://img.shields.io/github/downloads/lukasbudisky/terraform-provider-couchbase/total?style=flat-square&logo=terraform&logoColor=blue&labelColor=grey&color=yellow&link=https%3A%2F%2Fregistry.terraform.io%2Fproviders%2Flukasbudisky%2Fcouchbase%2Flatest&link=https%3A%2F%2Fregistry.terraform.io%2Fproviders%2Flukasbudisky%2Fcouchbase%2Flatest)
![GitHub contributors](https://img.shields.io/github/contributors-anon/lukasbudisky/terraform-provider-couchbase?style=flat-square&logo=github&logoColor=white&labelColor=grey&color=yellow&link=https%3A%2F%2Fgithub.com%2Flukasbudisky%2Fterraform-provider-couchbase%2Fgraphs%2Fcontributors&link=https%3A%2F%2Fgithub.com%2Flukasbudisky%2Fterraform-provider-couchbase%2Fgraphs%2Fcontributors)
![GitHub Repo stars](https://img.shields.io/github/stars/lukasbudisky/terraform-provider-couchbase?style=flat-square&logo=github&logoColor=white&labelColor=grey&color=yellow&link=https%3A%2F%2Fgithub.com%2Flukasbudisky%2Fterraform-provider-couchbase%2Fstargazers&link=https%3A%2F%2Fgithub.com%2Flukasbudisky%2Fterraform-provider-couchbase%2Fstargazers)


Terraform provider for Couchbase allow manage resources in couchbase cluster

## Requirements
- terraform 1.7.4
- go 1.22.1 (for plugin build)
- docker-compose v2.24.6-desktop.1
- docker desktop 4.28.0

## Run couchbase on localhost
In terraform_example folder is docker-compose.yml with couchbase server.

How to run couchbase on localhost. (Works on Ubuntu)
```bash
# Add couchbase to your /etc/hosts file
echo "127.0.0.1 couchbase" >> /etc/hosts

# Create couchbase network
make cbnetup

# Create couchbase
make cbup

# Couchbase initialization
make cbinit
```
How to destroy local infrastructure
```bash
# Destroy couchbase
make cbdown

# Destroy couchbase network
make cbnetdown
```

## Provider

> **WARNING**
>
> If you create multiple query indexes at once you can
> get internal server failure error because you can't
> create next index until previous is created.
>
> Suggested solution is to reduce parallelism.
> Add -parallelism=1 parameter during terraform apply
>
> Example:
>```bash
>terraform apply -parallelism=1
>```

### Base provider configuration
```terraform
terraform {
  required_version = ">= 1.7.4"
  required_providers {
    couchbase = {
      version = "~> 1.1.1"
      source  = "lukasbudisky/couchbase"
    }
  }
}

provider "couchbase" {
  address                   = "couchbase.couchbase"
  client_port               = 8091
  node_port                 = 11210
  username                  = "Administrator"
  password                  = "123456"
  management_timeout        = 10
  tls_root_cert_skip_verify = false
}
```

### TLS provider configuration
```terraform
terraform {
  required_version = ">= 1.7.4"
  required_providers {
    couchbase = {
      version = "~> 1.1.1"
      source  = "lukasbudisky/couchbase"
    }
  }
}

provider "couchbase" {
  address                   = "couchbase.couchbase"
  client_port               = 18091
  node_port                 = 11207
  username                  = "Administrator"
  password                  = "123456"
  management_timeout        = 10
  tls_root_cert_skip_verify = true
  tls_root_cert             = "certificate.pem"
  allow_sasl_mechanism      = "SCRAM-SHA1,SCRAM-SHA256,SCRAM-SHA512"
}
```

## Resources
We currently manage these operations via terraform resources
- buckets: ```couchbase_bucket_manager```
- groups: ```couchbase_security_group```
- users: ```couchbase_security_user```
- primary: query indexes ```couchbase_primary_query_index```
- query indexes: ```couchbase_query_index```

## Developing provider
Provider tests
```bash
make test
```
Acceptance tests
```bash
make testacc
```
Build provider
```bash
make build
```
Install provider
```bash
make install
```

