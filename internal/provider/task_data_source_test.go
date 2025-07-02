package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTaskDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Test data source with id
			{
				Config: testAccTaskDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.influxdb_task.test", "id"),
					resource.TestCheckResourceAttrSet("data.influxdb_task.test", "name"),
					resource.TestCheckResourceAttrSet("data.influxdb_task.test", "org_id"),
					resource.TestCheckResourceAttrSet("data.influxdb_task.test", "flux"),
					resource.TestCheckResourceAttrSet("data.influxdb_task.test", "status"),
				),
			},
		},
	})
}

const testAccTaskDataSourceConfig = `
# First get the organization to use its ID
data "influxdb_organizations" "all" {}

resource "influxdb_task" "test" {
  org_id = length(data.influxdb_organizations.all.organizations) > 0 ? data.influxdb_organizations.all.organizations[0].id : "dummy-org-id"
  flux   = <<-EOT
    option task = {
      name: "Test Task Data Source",
      every: 1h,
    }
    
    from(bucket: "test-bucket")
      |> range(start: -1h)
      |> filter(fn: (r) => r._measurement == "cpu")
      |> mean()
      |> to(bucket: "output-bucket", org: "test-org")
  EOT
}

data "influxdb_task" "test" {
  id = influxdb_task.test.id
}
`
