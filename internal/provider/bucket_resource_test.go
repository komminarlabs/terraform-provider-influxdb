package provider

import (
	"fmt"
	"os"
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
				Config: providerConfig + testAccBucketResourceWithRetentionConfig("test", "test bucket", "0"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("influxdb_bucket.test", "name", "test"),
					resource.TestCheckResourceAttr("influxdb_bucket.test", "description", "test bucket"),
					resource.TestCheckResourceAttr("influxdb_bucket.test", "retention_period", "0"),
				),
			},
			// ImportState testing
			{
				ResourceName: "influxdb_bucket.test",
				ImportState:  true,
			},
			// Update and Read testing
			{
				Config: providerConfig + testAccBucketResourceConfig("test-bucket", "test-bucket"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("influxdb_bucket.test", "name", "test-bucket"),
					resource.TestCheckResourceAttr("influxdb_bucket.test", "description", "test-bucket"),
					resource.TestCheckResourceAttr("influxdb_bucket.test", "retention_period", "2592000"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccBucketResourceWithRetentionConfig(name string, description string, retention_period string) string {
	return fmt.Sprintf(`
resource "influxdb_bucket" "test" {
  name = %[1]q
  description = %[2]q
  retention_period = %[3]q
  org_id = "`+os.Getenv("INFLUXDB_ORG_ID")+`"
}
`, name, description, retention_period)
}

func testAccBucketResourceConfig(name string, description string) string {
	return fmt.Sprintf(`
resource "influxdb_bucket" "test" {
  name = %[1]q
  description = %[2]q
  org_id = "`+os.Getenv("INFLUXDB_ORG_ID")+`"
}
`, name, description)
}
