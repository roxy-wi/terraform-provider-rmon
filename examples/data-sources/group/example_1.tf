provider "rmon" {
  base_url = "https://..."
  login    = "test"
  password = "testpass"
}

data "rmon_group" "example_id" {
  id = "4"
}

output "view" {
  value = data.rmon_group.example_id
}

// ------------------------------------

data "rmon_group" "example_name" {
  name = "test"
}

output "data" {
  value = data.rmon_group.example_name
}
