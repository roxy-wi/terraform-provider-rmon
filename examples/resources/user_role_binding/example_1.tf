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

resource "rmon_user_role_binding" "example" {
  user_id  = rmon_user.example.id
  role_id  = 1
  group_id = 1
}
