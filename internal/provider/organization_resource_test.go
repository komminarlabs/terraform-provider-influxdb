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
				Config: providerConfig + testAccOrganizationResourceConfig("test1", "This is a test organization"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("influxdb_organization.test", "name", "test1"),
					resource.TestCheckResourceAttr("influxdb_organization.test", "description", "This is a test organization"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "influxdb_organization.test",
				ImportState:       true,
				ImportStateVerify: true,
				// This is not normally necessary, but is here because this
				// example code does not have an actual upstream service.
				// Once the Read method is able to refresh information from
				// the upstream service, this can be removed.
				ImportStateVerifyIgnore: []string{"id", "65a453ed2e94b06f"},
			},
			// Update and Read testing
			{
				Config: providerConfig + testAccOrganizationResourceConfig("test", "test organization"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("influxdb_organization.test", "test", "test organization"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccOrganizationResourceConfig(c1 string, c2 string) string {
	return fmt.Sprintf(`
resource "influxdb_organization" "test" {
	name = %[1]q
	description = %[2]q
}
`, c1, c2)
}
