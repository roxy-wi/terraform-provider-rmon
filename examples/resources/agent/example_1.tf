provider "rmon" {
  base_url = "https://..."
  login    = "test"
  password = "testpass"
}

resource "rmon_agent" "example" {
  server_id   = 1
  name        = "Agent01"
  description = "Agent in my DC"
  enabled     = true
  shared      = true
  port        = 5101
  region_id   = 1
}
