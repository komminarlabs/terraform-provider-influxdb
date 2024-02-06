data "influxdb_organizations" "all" {}

output "all_organizations" {
  value = data.influxdb_organizations.all
}
