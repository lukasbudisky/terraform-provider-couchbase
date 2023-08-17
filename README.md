# Terraform provider for couchbase
[![Build Status](https://drone.budisky.com/api/badges/lukasbudisky/terraform-provider-couchbase/status.svg)](https://drone.budisky.com/lukasbudisky/terraform-provider-couchbase)

Terraform provider for Couchbase allow manage resources in couchbase cluster

## Requirements
- terraform 1.4.0
- go 1.21 (for plugin build)
- docker-compose v2.15.1
- docker desktop 4.17.0

## Run couchbase on localhost
In terraform_example folder is docker-compose.yml with couchbase server.

How to run couchbase on localhost. (Works on Ubuntu)
```
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
```
# Destroy couchbase
make cbdown

# Destroy couchbase network
make cbnetdown
```

## Provider

> **WARNING**
>
> If you create multiple query indexes at once you can get internal server failure error because you > > can't create next index until previous is created.
> Add -parallelism=1 parameter during terraform apply
>
> Example:
>```
>terraform apply -parallelism=1
>```

#### Base provider configuration
```
terraform {
  required_version = ">= 1.4.0"
  required_providers {
    couchbase = {
      version = "~> 0.0.6"
      source  = "budisky.com/couchbase/couchbase"
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

#### TLS provider configuration
```
terraform {
  required_version = ">= 1.4.0"
  required_providers {
    couchbase = {
      version = "~> 0.0.6"
      source  = "budisky.com/couchbase/couchbase"
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
```
make test
```
Acceptance tests
```
make testacc
```
Build provider
```
make build
```
Install provider
```
make install
```

