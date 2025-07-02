terraform {
  required_providers {
    influxdb = {
      source = "komminarlabs/influxdb"
    }
  }
}

provider "influxdb" {}

data "influxdb_label" "test" {
  id = "0f1d6e0c19f39000"
}

output "test_labels" {
  value = data.influxdb_label.test
}
