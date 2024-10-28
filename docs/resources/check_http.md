---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "rmon_check_http Resource - rmon"
subcategory: ""
description: |-
  This resource manages HTTP(s) check in RMON.
---

# rmon_check_http (Resource)

This resource manages HTTP(s) check in RMON.

## Example Usage

```terraform
provider "rmon" {
  base_url = "https://..."
  login    = "test"
  password = "testpass"
}

resource "rmon_check_http" "example" {
  name                  = "HTTP check"
  description           = "From TF"
  enabled               = true
  interval              = 60
  check_timeout         = 10
  place                 = "region"
  entities              = [2]
  url                   = "https://google.com"
  body_req              = <<EOF
  {"test": "test"}
EOF
  http_method           = "get"
  accepted_status_codes = 200
}
```


<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `entities` (List of Number) List of entities where check must be created.
- `http_method` (String) HTTP method for HTTP(s) check.
- `name` (String) Name of the Check Http.
- `place` (String) Port number for binding Check HTTP(s).
- `url` (String) URL what must be checked.

### Optional

- `accepted_status_codes` (Number) Expected status code (default to 200, optional).
- `body` (String) Check body answer.
- `body_req` (String) Send body to server. In JSON.
- `check_group` (String) Name of the check group for group HTTP(s) checks.
- `check_timeout` (Number) Answer timeout in seconds.
- `description` (String) Description of the Check HTTP(s).
- `enabled` (Boolean) Enabled state of the Check HTTP(s).
- `header_req` (String) Send headers to server. In JSON.
- `ignore_ssl_error` (Boolean) Ignore TLS/SSL error.
- `interval` (Number) Interval in seconds between checks.
- `mm_channel_id` (Number) Mattermost channel ID for alerts.
- `pd_channel_id` (Number) PagerDuty channel ID for alerts.
- `slack_channel_id` (Number) Slack channel ID for alerts.
- `telegram_channel_id` (Number) Telegram channel ID for alerts.
- `timeouts` (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String)
- `delete` (String)
- `update` (String)

## Import

In Terraform v1.7.0 and later, use an import block to import Country. For example:

```terraform
import {
  to = rmon_check_http.example
  id = "1"
}
```

Using terraform import, import Country can be imported using the `id`, e.g. For example:

```shell
% terraform import rmon_check_http.example 1
```