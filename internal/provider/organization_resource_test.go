package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOrganizationResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + testAccOrganizationResourceConfig("test", "This is a test organization"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("influxdb_organization.test", "name", "test"),
					resource.TestCheckResourceAttr("influxdb_organization.test", "description", "This is a test organization"),
				),
			},
			// ImportState testing
			{
				ResourceName: "influxdb_organization.test",
				ImportState:  true,
			},
			// Update and Read testing
			{
				Config: providerConfig + testAccOrganizationResourceConfig("test2", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("influxdb_organization.test", "name", "test2"),
					resource.TestCheckResourceAttr("influxdb_organization.test", "description", ""),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccOrganizationResourceConfig(name string, description string) string {
	return fmt.Sprintf(`
resource "influxdb_organization" "test" {
	name = %[1]q
	description = %[2]q
}
`, name, description)
}
