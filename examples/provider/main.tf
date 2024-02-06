terraform {
  required_providers {
    influxdb = {
      source = "registry.terraform.io/hashicorp/influxdb"
    }
  }
}

provider "influxdb" {}
