package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccUsersDataSource(t *testing.T) {
	userName1 := acctest.RandomWithPrefix("tf-user-test-1")
	userName2 := acctest.RandomWithPrefix("tf-user-test-2")
	password1 := acctest.RandomWithPrefix("password1")
	password2 := acctest.RandomWithPrefix("password2")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing - multiple users
			{
				Config: providerConfig + testAccUsersDataSourceConfig(userName1, password1, userName2, password2),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Check that we have at least the 2 users we created
					resource.TestCheckResourceAttrWith("data.influxdb_users.all", "users.#", func(value string) error {
						if value == "0" {
							return fmt.Errorf("expected at least 1 user, got %s", value)
						}
						return nil
					}),
					// Check that users have required attributes
					resource.TestCheckTypeSetElemNestedAttrs("data.influxdb_users.all", "users.*", map[string]string{
						"name":   userName1,
						"status": "active",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("data.influxdb_users.all", "users.*", map[string]string{
						"name":   userName2,
						"status": "inactive",
					}),
				),
			},
		},
	})
}

func TestAccUsersDataSourceWithOrganizations(t *testing.T) {
	userName1 := acctest.RandomWithPrefix("tf-user-test-1")
	userName2 := acctest.RandomWithPrefix("tf-user-test-2")
	password1 := acctest.RandomWithPrefix("password1")
	password2 := acctest.RandomWithPrefix("password2")
	orgName1 := acctest.RandomWithPrefix("tf-org-test-1")
	orgName2 := acctest.RandomWithPrefix("tf-org-test-2")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing - users with different organization roles
			{
				Config: providerConfig + testAccUsersDataSourceWithOrgsConfig(userName1, password1, userName2, password2, orgName1, orgName2),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Check that we have at least the users we created
					resource.TestCheckResourceAttrWith("data.influxdb_users.all", "users.#", func(value string) error {
						if value == "0" {
							return fmt.Errorf("expected at least 1 user, got %s", value)
						}
						return nil
					}),
					// Check that users have organization information
					resource.TestCheckTypeSetElemNestedAttrs("data.influxdb_users.all", "users.*", map[string]string{
						"name":     userName1,
						"org_role": "member",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("data.influxdb_users.all", "users.*", map[string]string{
						"name":     userName2,
						"org_role": "owner",
					}),
				),
			},
		},
	})
}

func TestAccUsersDataSourceEmpty(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing - should not fail even if no users exist
			{
				Config: providerConfig + testAccUsersDataSourceEmptyConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Should have users attribute even if empty
					resource.TestCheckResourceAttrSet("data.influxdb_users.all", "users.#"),
				),
			},
		},
	})
}

func testAccUsersDataSourceConfig(userName1, password1, userName2, password2 string) string {
	return fmt.Sprintf(`
resource "influxdb_user" "test1" {
	name     = "%s"
	password = "%s"
	status   = "active"
}

resource "influxdb_user" "test2" {
	name     = "%s"
	password = "%s"
	status   = "inactive"
}

data "influxdb_users" "all" {
	depends_on = [influxdb_user.test1, influxdb_user.test2]
}
`, userName1, password1, userName2, password2)
}

func testAccUsersDataSourceWithOrgsConfig(userName1, password1, userName2, password2, orgName1, orgName2 string) string {
	return fmt.Sprintf(`
resource "influxdb_organization" "test1" {
	name        = "%s"
	description = "Test organization 1 for users data source"
}

resource "influxdb_organization" "test2" {
	name        = "%s"
	description = "Test organization 2 for users data source"
}

resource "influxdb_user" "test1" {
	name     = "%s"
	password = "%s"
	status   = "active"
	org_id   = influxdb_organization.test1.id
	org_role = "member"
}

resource "influxdb_user" "test2" {
	name     = "%s"
	password = "%s"
	status   = "active"
	org_id   = influxdb_organization.test2.id
	org_role = "owner"
}

data "influxdb_users" "all" {
	depends_on = [influxdb_user.test1, influxdb_user.test2]
}
`, orgName1, orgName2, userName1, password1, userName2, password2)
}

const testAccUsersDataSourceEmptyConfig = `
data "influxdb_users" "all" {}
`
