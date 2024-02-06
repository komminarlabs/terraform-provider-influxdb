data "influxdb_bucket" "specific_bucket" {
    name = "_tasks"
}

output "specific_bucket" {
  value = data.influxdb_bucket.specific_bucket
}
