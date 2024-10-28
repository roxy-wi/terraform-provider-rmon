provider "rmon" {
  base_url = "https://..."
  login    = "test"
  password = "testpass"
}

resource "rmon_check_rabbitmq" "example" {
  name          = "RabbitMQ check"
  description   = "From TF"
  enabled       = true
  interval      = 60
  check_timeout = 10
  place         = "all"
  entities      = []
  ip            = "10.0.10.1"
  port          = 5672
  username      = "guest"
  password      = "guest"
  vhost         = "/"
}
