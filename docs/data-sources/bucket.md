---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "influxdb_bucket Data Source - terraform-provider-influxdb"
subcategory: ""
description: |-
  Retrieves a bucket. Use this data source to retrieve information for a specific bucket.
---

# influxdb_bucket (Data Source)

Retrieves a bucket. Use this data source to retrieve information for a specific bucket.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) A Bucket name.

### Read-Only

- `created_at` (String) Bucket creation date.
- `description` (String) A description of the bucket.
- `id` (String) A Bucket ID.
- `org_id` (String) An organization ID.
- `retention_period` (Number) The duration in seconds for how long data will be kept in the database. `0` represents infinite retention.
- `type` (String) The Bucket type.
- `updated_at` (String) Last bucket update date.
