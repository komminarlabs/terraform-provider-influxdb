data "influxdb_authorizations" "all" {}

output "all_authorizations" {
  value = data.influxdb_authorizations.all.authorizations[*].id
}
