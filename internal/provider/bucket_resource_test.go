package provider

/* import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccBucketResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccBucketResourceConfig("one"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("influxdb_bucket.test", "configurable_attribute", "one"),
					resource.TestCheckResourceAttr("influxdb_bucket.test", "defaulted", "example value when not configured"),
					resource.TestCheckResourceAttr("influxdb_bucket.test", "id", "example-id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "influxdb_bucket.test",
				ImportState:       true,
				ImportStateVerify: true,
				// This is not normally necessary, but is here because this
				// example code does not have an actual upstream service.
				// Once the Read method is able to refresh information from
				// the upstream service, this can be removed.
				ImportStateVerifyIgnore: []string{"configurable_attribute", "defaulted"},
			},
			// Update and Read testing
			{
				Config: testAccBucketResourceConfig("two"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("influxdb_bucket.test", "configurable_attribute", "two"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccBucketResourceConfig(configurableAttribute string) string {
	return fmt.Sprintf(`
resource "influxdb_bucket" "test" {
  configurable_attribute = %[1]q
}
`, configurableAttribute)
}*/
