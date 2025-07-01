package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccLabelDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create a label resource first, then read it with data source
			{
				Config: providerConfig + testAccLabelDataSourceConfig("test-label-data"),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Check the created resource
					resource.TestCheckResourceAttr("influxdb_label.test", "name", "test-label-data"),
					resource.TestCheckResourceAttr("influxdb_label.test", "org_id", os.Getenv("INFLUXDB_ORG_ID")),
					// Check the data source
					resource.TestCheckResourceAttr("data.influxdb_label.test", "name", "test-label-data"),
					resource.TestCheckResourceAttr("data.influxdb_label.test", "org_id", os.Getenv("INFLUXDB_ORG_ID")),
					resource.TestCheckResourceAttrPair("data.influxdb_label.test", "id", "influxdb_label.test", "id"),
				),
			},
		},
	})
}

func TestAccLabelDataSourceWithProperties(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create a label resource with properties, then read it with data source
			{
				Config: providerConfig + testAccLabelDataSourceWithPropertiesConfig("test-label-data-props"),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Check the created resource
					resource.TestCheckResourceAttr("influxdb_label.test", "name", "test-label-data-props"),
					resource.TestCheckResourceAttr("influxdb_label.test", "properties.color", "green"),
					resource.TestCheckResourceAttr("influxdb_label.test", "properties.environment", "testing"),
					// Check the data source
					resource.TestCheckResourceAttr("data.influxdb_label.test", "name", "test-label-data-props"),
					resource.TestCheckResourceAttr("data.influxdb_label.test", "properties.color", "green"),
					resource.TestCheckResourceAttr("data.influxdb_label.test", "properties.environment", "testing"),
					resource.TestCheckResourceAttrPair("data.influxdb_label.test", "id", "influxdb_label.test", "id"),
				),
			},
		},
	})
}

func testAccLabelDataSourceConfig(name string) string {
	return fmt.Sprintf(`
resource "influxdb_label" "test" {
  name   = %[1]q
  org_id = "`+os.Getenv("INFLUXDB_ORG_ID")+`"
}

data "influxdb_label" "test" {
  id = influxdb_label.test.id
}
`, name)
}

func testAccLabelDataSourceWithPropertiesConfig(name string) string {
	return fmt.Sprintf(`
resource "influxdb_label" "test" {
  name   = %[1]q
  org_id = "`+os.Getenv("INFLUXDB_ORG_ID")+`"
  properties = {
    "color"       = "green"
    "environment" = "testing"
  }
}

data "influxdb_label" "test" {
  id = influxdb_label.test.id
}
`, name)
}
