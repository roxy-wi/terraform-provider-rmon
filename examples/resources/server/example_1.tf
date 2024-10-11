provider "rmon" {
  base_url = "https://..."
  login    = "test"
  password = "testpass"
}

resource "rmon_server" "example" {
  cred_id     = 2
  description = "test server"
  enabled     = true
  group_id    = 1
  hostname    = "d-infra-redis01"
  ip          = "192.168.1.101"
  port        = 5673
}
