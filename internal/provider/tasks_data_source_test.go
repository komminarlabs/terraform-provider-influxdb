package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTasksDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Test data source that lists all tasks
			{
				Config: testAccTasksDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.influxdb_tasks.all", "tasks.#"),
					// Check that tasks are properly populated
					resource.TestCheckResourceAttrSet("data.influxdb_tasks.all", "tasks.0.id"),
					resource.TestCheckResourceAttrSet("data.influxdb_tasks.all", "tasks.0.name"),
					resource.TestCheckResourceAttrSet("data.influxdb_tasks.all", "tasks.0.org_id"),
					resource.TestCheckResourceAttrSet("data.influxdb_tasks.all", "tasks.0.flux"),
					resource.TestCheckResourceAttrSet("data.influxdb_tasks.all", "tasks.0.status"),
					resource.TestCheckResourceAttrSet("data.influxdb_tasks.all", "tasks.0.created_at"),
				),
			},
		},
	})
}

func TestAccTasksDataSourceWithMultipleTasks(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Test data source with multiple tasks created
			{
				Config: testAccTasksDataSourceConfigMultiple,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Should have at least 2 tasks (the ones we created)
					resource.TestCheckResourceAttr("data.influxdb_tasks.all", "tasks.#", "2"),
					// Check first task
					resource.TestCheckResourceAttrSet("data.influxdb_tasks.all", "tasks.0.id"),
					resource.TestCheckResourceAttrSet("data.influxdb_tasks.all", "tasks.0.name"),
					resource.TestCheckResourceAttrSet("data.influxdb_tasks.all", "tasks.0.flux"),
					// Check second task
					resource.TestCheckResourceAttrSet("data.influxdb_tasks.all", "tasks.1.id"),
					resource.TestCheckResourceAttrSet("data.influxdb_tasks.all", "tasks.1.name"),
					resource.TestCheckResourceAttrSet("data.influxdb_tasks.all", "tasks.1.flux"),
				),
			},
		},
	})
}

const testAccTasksDataSourceConfig = `
# First get the organization to use its ID
data "influxdb_organizations" "all" {}

# Create a test task first
resource "influxdb_task" "test_for_datasource" {
  org_id = length(data.influxdb_organizations.all.organizations) > 0 ? data.influxdb_organizations.all.organizations[0].id : "dummy-org-id"
  flux   = <<-EOT
    option task = {
      name: "Test Task for DataSource",
      every: 1h,
    }
    
    from(bucket: "test-bucket")
      |> range(start: -1h)
      |> filter(fn: (r) => r._measurement == "cpu")
      |> mean()
      |> to(bucket: "output-bucket", org: "test-org")
  EOT
}

# Query all tasks
data "influxdb_tasks" "all" {
  depends_on = [influxdb_task.test_for_datasource]
}
`

const testAccTasksDataSourceConfigMultiple = `
# First get the organization to use its ID
data "influxdb_organizations" "all" {}

# Create first test task
resource "influxdb_task" "test_cron_for_datasource" {
  org_id = length(data.influxdb_organizations.all.organizations) > 0 ? data.influxdb_organizations.all.organizations[0].id : "dummy-org-id"
  flux   = <<-EOT
    option task = {
      name: "Test Cron Task for DataSource",
      cron: "0 0 * * *",
    }
    
    from(bucket: "test-bucket")
      |> range(start: -24h)
      |> filter(fn: (r) => r._measurement == "temperature")
      |> mean()
      |> to(bucket: "daily-temps", org: "test-org")
  EOT
}

# Create second test task
resource "influxdb_task" "test_every_for_datasource" {
  org_id = length(data.influxdb_organizations.all.organizations) > 0 ? data.influxdb_organizations.all.organizations[0].id : "dummy-org-id"
  flux   = <<-EOT
    option task = {
      name: "Test Every Task for DataSource",
      every: 30m,
    }
    
    from(bucket: "test-bucket")
      |> range(start: -30m)
      |> filter(fn: (r) => r._measurement == "pressure")
      |> mean()
      |> to(bucket: "pressure-data", org: "test-org")
  EOT
}

# Query all tasks
data "influxdb_tasks" "all" {
  depends_on = [
    influxdb_task.test_cron_for_datasource,
    influxdb_task.test_every_for_datasource
  ]
}
`
