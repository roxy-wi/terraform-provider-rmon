provider "rmon" {
  base_url = "https://your_rmon"
  login    = "testlog"
  password = "testpass"
}

resource "rmon_channel" "example" {
  receiver = "pd"
  channel  = "test_my_channel"
  group_id = 1
  token    = "some_token"
}
