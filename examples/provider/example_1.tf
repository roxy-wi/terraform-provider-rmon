terraform {
  required_providers {
    rmon = {
      source = "Roxy-wi/rmon"
    }
  }
}

provider "rmon" {
  base_url = "https://..."
  login    = "test"
  password = "testpass"
}