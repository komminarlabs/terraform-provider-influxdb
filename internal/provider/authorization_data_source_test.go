package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAuthorizationDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + testAccAuthorizationDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.influxdb_authorization.test", "permissions.#", "2"),
					resource.TestCheckResourceAttr("data.influxdb_authorization.test", "description", "Access test bucket"),
				),
			},
		},
	})
}

func testAccAuthorizationDataSourceConfig() string {
	return `
resource "influxdb_bucket" "test" {
	name = "test"
	org_id = "` + os.Getenv("INFLUXDB_ORG_ID") + `"
  }

resource "influxdb_authorization" "test" {
	org_id      = "` + os.Getenv("INFLUXDB_ORG_ID") + `"
	description = "Access test bucket"
  
	permissions = [{
	  action = "read"
	  resource = {
		id     = influxdb_bucket.test.id
		org_id = "` + os.Getenv("INFLUXDB_ORG_ID") + `"
		type   = "buckets"
	  }
	  },
	  {
		action = "write"
		resource = {
		  id     = influxdb_bucket.test.id
		  org_id = "` + os.Getenv("INFLUXDB_ORG_ID") + `"
		  type   = "buckets"
		}
	}]
  }

  data "influxdb_authorization" "test" {
	id = influxdb_authorization.test.id
  }
`
}
