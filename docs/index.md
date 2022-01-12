---
page_title: "Provider: Camunda"
subcategory: ""
description: |-
  Terraform provider for interacting with Camunda API.
---

# Camunda Provider

The Camunda provider is used to interact with Camunda (self hosted).

Use the navigation to the left to read about the available resources.

## Example Usage

Do not keep your authentication password in HCL for production environments, use Terraform environment variables.

```terraform
provider "camunda" {
  endpoint = "http://localhost/engine-rest"
  username = "username"
  password = "password"
  tls      = {
    insecure_skip_verify = true
  }
}
```

## Schema

### Required

- **endpoint** (String, Required) Camunda engine API address

### Optional

- **username** (String, Optional) Username to authenticate to Camunda API
- **password** (String, Optional) Password to authenticate to Camunda API
- **tls** (Tls, Optional) Tls configuration to communicate with the Camunda API
