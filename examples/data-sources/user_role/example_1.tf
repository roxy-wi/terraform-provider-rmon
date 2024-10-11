provider "rmon" {
  base_url = "https://..."
  login    = "test"
  password = "testpass"
}

data "rmon_user_role" "example" {}

output "test" {
  value = data.rmon_user_role.example.roles
}
