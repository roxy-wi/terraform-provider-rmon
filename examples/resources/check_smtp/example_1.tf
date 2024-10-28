provider "rmon" {
  base_url = "https://..."
  login    = "test"
  password = "testpass"
}

resource "rmon_check_smtp" "example" {
  name          = "SMTP check"
  description   = "From TF"
  enabled       = true
  interval      = 60
  check_timeout = 10
  place         = "region"
  entities      = [1]
  ip            = "smtp.example.com"
  port          = 25
  username      = "some@example.com"
  password      = "password"
  check_group   = "SMTP checks"
}
