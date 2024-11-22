provider "rmon" {
  base_url = "https://..."
  login    = "test"
  password = "testpass"
}

resource "rmon_region" "example" {
  name        = "DC01"
  description = "This region consists all Agents in DC01"
  enabled     = true
  shared      = true
  country_id  = 1
}
