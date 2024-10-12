provider "rmon" {
  base_url = "https://your_rmon"
  login    = "testlog"
  password = "testpass"
}

resource "rmon_user" "example" {
  email    = "test23@gmail.com"
  enabled  = true
  password = "testpassword"
  username = "testuser2"
}
