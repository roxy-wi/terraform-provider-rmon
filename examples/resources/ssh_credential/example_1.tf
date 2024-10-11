provider "rmon" {
  base_url = "https://..."
  login    = "test"
  password = "testpass"
}

resource "rmon_ssh_credential" "example" {
  group_id    = 1
  name        = "test_cred"
  username    = "root test"
  password    = "test23"
  key_enabled = false
  shared      = true
}