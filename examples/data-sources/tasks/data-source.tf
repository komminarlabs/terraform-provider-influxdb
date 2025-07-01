terraform {
  required_providers {
    influxdb = {
      source = "komminarlabs/influxdb"
    }
  }
}

provider "influxdb" {}

data "influxdb_tasks" "all" {}

output "all_tasks" {
  value = data.influxdb_tasks.all
}
