---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "RMON Provider"
subcategory: ""
description: |-
  
---

# RMON Provider

Use the RMON provider to interact with the many resources supported by RMON. You must configure the provider with the proper credentials before you can use it.

Use the navigation to the left to read about the available resources.

To learn the basics of Terraform using this provider, follow the hands-on get started tutorials.

## Example Usage

Terraform 1.7.0 and later:

```terraform
terraform {
  required_providers {
    rmon = {
      source = "Roxy-wi/rmon"
    }
  }
}

provider "rmon" {
  base_url = "https://..."
  login    = "test"
  password = "testpass"
}
```

<!-- schema generated by tfplugindocs -->

### Provider Configuration

!> **Warning:** Hard-coded credentials are not recommended in any Terraform configuration and risks secret leakage should this file ever be committed to a public version control system.

### Environment Variables

Credentials can be provided by using the `RMON_USERNAME`, `RMON_PASSWORD` and url auth `RMON_BASE_URL` environment variables.

For example:

```terraform
provider "rmon" {}
```

```shell
% export RMON_USERNAME="user"
% export RMON_PASSWORD="password"
% export RMON_BASE_URL="https://your_rmon"
% terraform plan
```

## Schema

### Required

- `login` (String) Username for RMON.
- `password` (String) Password for RMON.
- `base_url` (String) URL to connect for RMON.
