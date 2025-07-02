package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccLabelResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + testAccLabelResourceConfig("test-label", "test label description"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("influxdb_label.test", "name", "test-label"),
					resource.TestCheckResourceAttrSet("influxdb_label.test", "id"),
					resource.TestCheckResourceAttr("influxdb_label.test", "org_id", os.Getenv("INFLUXDB_ORG_ID")),
				),
			},
			// ImportState testing
			{
				ResourceName: "influxdb_label.test",
				ImportState:  true,
			},
			// Update and Read testing
			{
				Config: providerConfig + testAccLabelResourceConfig("test-label-updated", "updated test label"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("influxdb_label.test", "name", "test-label-updated"),
					resource.TestCheckResourceAttrSet("influxdb_label.test", "id"),
					resource.TestCheckResourceAttr("influxdb_label.test", "org_id", os.Getenv("INFLUXDB_ORG_ID")),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccLabelResourceWithProperties(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with properties
			{
				Config: providerConfig + testAccLabelResourceWithPropertiesConfig("test-label-props", map[string]string{
					"color":       "blue",
					"environment": "test",
					"team":        "engineering",
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("influxdb_label.test", "name", "test-label-props"),
					resource.TestCheckResourceAttr("influxdb_label.test", "properties.color", "blue"),
					resource.TestCheckResourceAttr("influxdb_label.test", "properties.environment", "test"),
					resource.TestCheckResourceAttr("influxdb_label.test", "properties.team", "engineering"),
				),
			},
			// Update properties - modify, add, and remove
			{
				Config: providerConfig + testAccLabelResourceWithPropertiesConfig("test-label-props", map[string]string{
					"color":       "red",        // modified
					"environment": "production", // modified
					"priority":    "high",       // added
					// "team" removed
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("influxdb_label.test", "name", "test-label-props"),
					resource.TestCheckResourceAttr("influxdb_label.test", "properties.color", "red"),
					resource.TestCheckResourceAttr("influxdb_label.test", "properties.environment", "production"),
					resource.TestCheckResourceAttr("influxdb_label.test", "properties.priority", "high"),
					resource.TestCheckNoResourceAttr("influxdb_label.test", "properties.team"),
				),
			},
			// Remove all properties
			{
				Config: providerConfig + testAccLabelResourceConfig("test-label-props", "no properties"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("influxdb_label.test", "name", "test-label-props"),
					resource.TestCheckNoResourceAttr("influxdb_label.test", "properties.color"),
					resource.TestCheckNoResourceAttr("influxdb_label.test", "properties.environment"),
					resource.TestCheckNoResourceAttr("influxdb_label.test", "properties.priority"),
				),
			},
		},
	})
}

func TestAccLabelResourcePropertyRemoval(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with properties
			{
				Config: providerConfig + testAccLabelResourceWithPropertiesConfig("test-property-removal", map[string]string{
					"key1": "value1",
					"key2": "value2",
					"key3": "value3",
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("influxdb_label.test", "properties.key1", "value1"),
					resource.TestCheckResourceAttr("influxdb_label.test", "properties.key2", "value2"),
					resource.TestCheckResourceAttr("influxdb_label.test", "properties.key3", "value3"),
				),
			},
			// Remove specific property by setting to empty string
			{
				Config: providerConfig + testAccLabelResourceWithPropertiesConfig("test-property-removal", map[string]string{
					"key1": "value1",
					"key3": "value3",
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("influxdb_label.test", "properties.key1", "value1"),
					resource.TestCheckNoResourceAttr("influxdb_label.test", "properties.key2"),
					resource.TestCheckResourceAttr("influxdb_label.test", "properties.key3", "value3"),
				),
			},
		},
	})
}

func testAccLabelResourceConfig(name string, _ string) string {
	return fmt.Sprintf(`
resource "influxdb_label" "test" {
  name   = %[1]q
  org_id = "`+os.Getenv("INFLUXDB_ORG_ID")+`"
}
`, name)
}

func testAccLabelResourceWithPropertiesConfig(name string, properties map[string]string) string {
	propertiesStr := ""
	if len(properties) > 0 {
		propertiesStr = "properties = {\n"
		for key, value := range properties {
			propertiesStr += fmt.Sprintf("    %q = %q\n", key, value)
		}
		propertiesStr += "  }"
	}

	return fmt.Sprintf(`
resource "influxdb_label" "test" {
  name   = %[1]q
  org_id = "`+os.Getenv("INFLUXDB_ORG_ID")+`"
  %[2]s
}
`, name, propertiesStr)
}
