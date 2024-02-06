data "influxdb_buckets" "all" {}

output "all_buckets" {
  value = data.influxdb_buckets.all
}
