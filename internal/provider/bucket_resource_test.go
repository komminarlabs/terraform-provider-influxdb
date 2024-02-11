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
				Config: providerConfig + testAccBucketResourceConfig("test1", "test bucket", "12c1df6c262377a5"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("influxdb_bucket.test", "name", "test1"),
					resource.TestCheckResourceAttr("influxdb_bucket.test", "description", "test bucket"),
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
				ImportStateVerifyIgnore: []string{"id", "defaulted"},
			},
			// Update and Read testing
			{
				Config: providerConfig + testAccBucketResourceConfig("test1", "", "12c1df6c262377a5"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("influxdb_bucket.test", "name", "test1"),
					resource.TestCheckResourceAttr("influxdb_bucket.test", "description", ""),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccBucketResourceConfig(c1 string, c2 string, c3 string) string {
	return fmt.Sprintf(`
resource "influxdb_bucket" "test" {
  name = %[1]q
  description = %[2]q 
  org_id = %[3]q
}
`, c1, c2, c3)
} */
