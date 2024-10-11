provider "rmon" {
  base_url = "https://demo.roxy-wi.org"
  login    = "testlog"
  password = "testpass"
}

resource "rmon_user" "example" {
  email    = "test23@yandex.ru"
  enabled  = true
  password = "testpassword"
  username = "testuser2"
}
