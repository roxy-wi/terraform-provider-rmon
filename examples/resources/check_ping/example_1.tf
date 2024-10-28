provider "rmon" {
  base_url = "https://..."
  login    = "test"
  password = "testpass"
}

resource "rmon_check_ping" "example" {
  name          = "Ping check"
  description   = "From TF"
  enabled       = true
  packet_size   = 56
  interval      = 60
  check_timeout = 10
  place         = "agent"
  entities      = [1]
  ip            = "example.com"
}
