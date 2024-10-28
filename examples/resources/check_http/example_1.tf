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
