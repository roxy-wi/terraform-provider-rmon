provider "rmon" {
  base_url = "https://..."
  login    = "test"
  password = "testpass"
}

resource "rmon_country" "example" {
  name        = "My Country"
  description = "This country consists some regions"
  enabled     = true
  shared      = true
  regions     = [1, 2, 3]
}
