terraform {
  required_providers {
    influxdb = {
      source = "komminarlabs/influxdb"
    }
  }
}

provider "influxdb" {}
