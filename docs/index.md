---
layout: "couchbase"
page_title: "terraform-provider-couchbase"
sidebar_current: "docs-couchbase-index"
description: |-
  Terraform provider for couchbase
---

# terraform-provider-couchbase

The terraform couchbase provider `terraform-provider-couchbase` for management resources in couchbase

## Argument reference

The following arguments are supported
### Required

<ul>
  <li><b>address</b> (String) Couchase server address</li>
  <ul>
    <li><b>CB_ADDRESS</b> Environment variable</li>
  </ul>
  <li><b>client_port</b> (Int) Couchase server port: client-to-node</li>
  <ul>
    <li><b>CB_CLIENT_PORT</b> Environment variable</li>
  </ul>
  <li><b>node_port</b> (Int) Couchase server port: node-to-node</li>
  <ul>
    <li><b>CB_NODE_PORT</b> Environment variable</li>
  </ul>
  <li><b>username</b> (String) Couchase username</li>
  <ul>
    <li><b>CB_USERNAME</b> Environment variable</li>
  </ul>
  <li><b>password</b> (String) Couchase password</li>
  <ul>
    <li><b>CB_PASSWORD</b> Environment variable</li>
  </ul>
</ul>

**More information about client-to-node and node-to-node ports are in couchbase documentation**
<a href=https://docs.couchbase.com/server/current/install/install-ports.html>Couchbase Ports Documentation</a>

### Optional
<ul>
  <li><b>management_timeout</b> (String) Couchase management timeout. Read more about couchbase timeouts in documentation</li>
  <ul>
    <li><b>CB_MANAGEMENT_TIMEOUT</b> Environment variable</li>
  </ul>
  <li><b>tls_root_cert_skip_verify</b> (Bool) Skip root CA verification</li>
  <ul>
    <li><b>TLS_ROOT_CERT_SKIP_VERIFY</b>Environment variable</li>
  </ul>
  <li><b>allow_sasl_mechanism</b> (String) Allow sasl mechanism - multiple values shoud be separated with commas in one string ","</li>
  <ul>
    <li><b>TLS_ROOT_CERT_ALLOW_SASL_MECHANISM</b>Environment variable</li>
  </ul>
  <li><b>tls_root_cert</b> (String) Path to certificate</li>
  <ul>
    <li><b>TLS_ROOT_CERT</b>Environment variable</li>
  </ul>
</ul>

**More information about timeouts are in client settings couchbase documentation**
<a href="https://docs.couchbase.com/ruby-sdk/current/ref/client-settings.html">Client Settings / Timeouts</a>

## Example usage
####Minimal configuration
```
terraform {
  required_version = ">= 0.13"
  required_providers {
    couchbase = {
      version = "~> 0.0.3"
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
}
```

####TLS Configuration
```
terraform {
  required_version = ">= 0.13"
  required_providers {
    couchbase = {
      version = "~> 0.0.3"
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

####Example create new bucket
```
terraform {
  required_version = ">= 0.13"
  required_providers {
    couchbase = {
      version = "~> 0.0.3"
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

resource "couchbase_bucket_manager" "bucket_1" {
  name                     = "bucket_1"
  ram_quota_mb             = 512
  flush_enabled            = false
  max_expire               = 0
  conflict_resolution_type = "seqno"
  compression_mode         = "passive"
  num_replicas             = 1
}
```
