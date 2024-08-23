terraform {
  required_providers {
    influxdb = {
      source = "komminarlabs/influxdb"
    }
  }
}

provider "influxdb" {}

data "influxdb_authorization" "signals_authorization" {
  id = "0c91163b3f53e000"
}

output "signals_authorization" {
  value = data.influxdb_authorization.signals_authorization.id
}
