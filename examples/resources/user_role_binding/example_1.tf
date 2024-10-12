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

resource "rmon_user_role_binding" "example" {
  user_id  = rmon_user.example.id
  role_id  = 1
  group_id = 1
}
