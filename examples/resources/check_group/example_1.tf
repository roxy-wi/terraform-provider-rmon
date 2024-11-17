provider "rmon" {
  base_url = "https://..."
  login    = "test"
  password = "testpass"
}

resource "rmon_check_group" "example" {
  name = "example_group"
}
