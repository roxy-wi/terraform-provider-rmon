provider "rmon" {
  base_url = "https://..."
  login    = "test"
  password = "testpass"
}

resource "rmon_check_tcp" "example" {
  name          = "TCP check"
  description   = "From TF"
  enabled       = true
  interval      = 60
  check_timeout = 10
  place         = "agent"
  entities      = [1]
  ip            = "example.com"
  port          = 22
  check_group   = "TCP"
}
