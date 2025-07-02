terraform {
  required_providers {
    influxdb = {
      source = "komminarlabs/influxdb"
    }
  }
}

provider "influxdb" {}

data "influxdb_users" "all" {}

output "all_users" {
  value     = data.influxdb_users.all
  sensitive = true
}
