package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccBucketDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + testAccBucketDataSourceConfig("_monitoring"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.influxdb_bucket.test", "name", "_monitoring"),
					resource.TestCheckResourceAttr("data.influxdb_bucket.test", "type", "system"),
				),
			},
		},
	})
}

func testAccBucketDataSourceConfig(name string) string {
	return fmt.Sprintf(`
data "influxdb_bucket" "test" {
	name = %[1]q
}
`, name)
}
