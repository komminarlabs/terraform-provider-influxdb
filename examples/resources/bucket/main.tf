terraform {
  required_providers {
    influxdb = {
      source = "registry.terraform.io/hashicorp/influxdb"
    }
  }
}

provider "influxdb" {}

resource "influxdb_bucket" "sample" {
  org_id         = "12c1df6c262377a5"
  name           = "sample"
  description    = "This is a sample bucket 1"
  retention_days = 30
}

output "sample_bucket" {
  value = influxdb_bucket.sample
}
