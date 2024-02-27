package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOrganizationDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + testAccOrganizationDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.influxdb_organization.test", "name", "test"),
					resource.TestCheckResourceAttr("data.influxdb_organization.test", "description", "This is a test organization"),
				),
			},
		},
	})
}

const testAccOrganizationDataSourceConfig = `
resource "influxdb_organization" "test" {
	name = "test"
	description = "This is a test organization"
}

data "influxdb_organization" "test" {
	name = "test"
	depends_on = [influxdb_organization.test]
}
`
