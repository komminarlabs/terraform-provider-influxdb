package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOrganizationsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + testAccOrganizationsDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.influxdb_organizations.all", "organizations.0.name", "default"),
					resource.TestCheckResourceAttr("data.influxdb_organizations.all", "organizations.0.description", ""),
				),
			},
		},
	})
}

const testAccOrganizationsDataSourceConfig = `
data "influxdb_organizations" "all" {}
`
