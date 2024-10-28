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
