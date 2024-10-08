## [1.1.2] - 2024-09-12

## Updated:
* Updated docs to include supported influxdb flavours.

## [1.1.1] - 2024-08-23

## Fixed:
* fixed overwriting the token value during read in `influxdb_authorization` resource.
  
## [1.1.0] - 2024-04-16

## Updated:
* Updated `influxdb_authorization` resource and made the `id` & `org_id` as optional in `permissions.resource` inline with Influx api.
* Updated `influxdb_authorization` resource and made the `name` as read-only in `permissions.resource`. This is due to how Influx api returns the response. This will be modified in the future versions.

## [1.0.1] - 2024-03-05

## Updated:
* Updated provider docs
* Added Issue Templates
* Bump github.com/hashicorp/terraform-plugin-go from `0.21.0` to `0.22.0`
* Bump github.com/hashicorp/terraform-plugin-framework from `1.5.0` to `1.6.0`
* Upgrade golang
* Renamed module in `go.mod`

## [1.0.0] - 2024-02-27

## Added:

* Added an optional attribute `name` to `influxdb_authorization` resource.
* Acceptance tests for all data sources and resources.

## Updated:

* `retention_days` is renamed to `retention_period` in `influxdb_bucket` resource.
* Made some document changes.

## [0.1.0] - 2024-02-14

### Added:

* **New Data Source:** `influxdb_authorization`
* **New Data Source:** `influxdb_authorizations`
* **New Data Source:** `influxdb_bucket`
* **New Data Source:** `influxdb_buckets`
* **New Data Source:** `influxdb_organization`
* **New Data Source:** `influxdb_organizations`
* **New Resource:** `influxdb_authorization`
* **New Resource:** `influxdb_bucket`
* **New Resource:** `influxdb_organization`
