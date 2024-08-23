terraform {
  required_providers {
    influxdb = {
      source = "komminarlabs/influxdb"
    }
  }
}

provider "influxdb" {}

data "influxdb_bucket" "signals" {
  name = "signals"
}

output "signals_bucket" {
  value = data.influxdb_bucket.signals
}
