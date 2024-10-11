provider "rmon" {
  base_url = "https://..."
  login    = "test"
  password = "testpass"
}

resource "rmon_group" "example" {
  name        = "example_group2"
  description = "test terraform group2"
}
