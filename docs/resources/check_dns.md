---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "rmon_check_dns Resource - rmon"
subcategory: ""
description: |-
  This resource manages DNS check in RMON.
---

# rmon_check_dns (Resource)

This resource manages DNS check in RMON.

## Example Usage

```terraform
provider "rmon" {
  base_url = "https://..."
  login    = "test"
  password = "testpass"
}

resource "rmon_check_dns" "example" {
  name          = "DNS check"
  description   = "From TF"
  enabled       = true
  interval      = 60
  check_timeout = 10
  place         = "agent"
  entities      = [1]
  ip            = "google.com"
  port          = 53
  check_group   = "DNS"
  record_type   = "a"
  resolver      = "8.8.8.8"
}
```


<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `entities` (List of Number) List of entities where check must be created.
- `name` (String) Name of the CheckDns.
- `place` (String) Port number for binding Check DNS.

### Optional

- `check_group` (String) Name of the check group for group DNS checks.
- `check_timeout` (Number) Answer timeout in seconds.
- `description` (String) Description of the CheckDns.
- `enabled` (Boolean) Enabled state of the Check DNS.
- `interval` (Number) Interval in seconds between checks.
- `ip` (String) IP address or domain name for check.
- `mm_channel_id` (Number) Mattermost channel ID for alerts.
- `pd_channel_id` (Number) PagerDuty channel ID for alerts.
- `port` (Number) Packet size in bytes.
- `record_type` (String) DNS record type.
- `resolver` (String) DNS server where resolve DNS query.
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
  to = rmon_check_dns.example
  id = "1"
}
```

Using terraform import, import Country can be imported using the `id`, e.g. For example:

```shell
% terraform import rmon_check_dns.example 1
```