package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccUserResource(t *testing.T) {
	userName := acctest.RandomWithPrefix("tf-user-test")
	password := acctest.RandomWithPrefix("password")
	updatedPassword := acctest.RandomWithPrefix("updated-password")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + testAccUserResourceConfig(userName, password),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("influxdb_user.test", "name", userName),
					resource.TestCheckResourceAttr("influxdb_user.test", "password", password),
					resource.TestCheckResourceAttr("influxdb_user.test", "status", "active"),
					resource.TestCheckResourceAttrSet("influxdb_user.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "influxdb_user.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "org_id", "org_role"}, // Password is not returned by API, org attributes need to be set manually after import
			},
			// Update and Read testing
			{
				Config: providerConfig + testAccUserResourceConfigUpdated(userName, updatedPassword),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("influxdb_user.test", "name", userName),
					resource.TestCheckResourceAttr("influxdb_user.test", "password", updatedPassword),
					resource.TestCheckResourceAttr("influxdb_user.test", "status", "inactive"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccUserResourceWithOrganization(t *testing.T) {
	userName := acctest.RandomWithPrefix("tf-user-org-test")
	password := acctest.RandomWithPrefix("password")
	orgName := acctest.RandomWithPrefix("tf-org-test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create user with organization membership as member
			{
				Config: providerConfig + testAccUserResourceWithOrgConfig(userName, password, orgName, "member"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("influxdb_user.test", "name", userName),
					resource.TestCheckResourceAttr("influxdb_user.test", "password", password),
					resource.TestCheckResourceAttr("influxdb_user.test", "org_role", "member"),
					resource.TestCheckResourceAttrSet("influxdb_user.test", "id"),
					resource.TestCheckResourceAttrSet("influxdb_user.test", "org_id"),
				),
			},
			// Update role to owner
			{
				Config: providerConfig + testAccUserResourceWithOrgConfig(userName, password, orgName, "owner"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("influxdb_user.test", "name", userName),
					resource.TestCheckResourceAttr("influxdb_user.test", "org_role", "owner"),
				),
			},
			// Remove from organization
			{
				Config: providerConfig + testAccUserResourceConfig(userName, password),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("influxdb_user.test", "name", userName),
					resource.TestCheckNoResourceAttr("influxdb_user.test", "org_id"),
					resource.TestCheckNoResourceAttr("influxdb_user.test", "org_role"),
				),
			},
		},
	})
}

func TestAccUserResourceOrgSwitching(t *testing.T) {
	userName := acctest.RandomWithPrefix("tf-user-switch-test")
	password := acctest.RandomWithPrefix("password")
	org1Name := acctest.RandomWithPrefix("tf-org1-test")
	org2Name := acctest.RandomWithPrefix("tf-org2-test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create user with first organization as member
			{
				Config: providerConfig + testAccUserResourceWithTwoOrgsConfig(userName, password, org1Name, org2Name, "org1", "member"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("influxdb_user.test", "name", userName),
					resource.TestCheckResourceAttr("influxdb_user.test", "org_role", "member"),
					resource.TestCheckResourceAttrSet("influxdb_user.test", "org_id"),
				),
			},
			// Switch to second organization as owner
			{
				Config: providerConfig + testAccUserResourceWithTwoOrgsConfig(userName, password, org1Name, org2Name, "org2", "owner"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("influxdb_user.test", "name", userName),
					resource.TestCheckResourceAttr("influxdb_user.test", "org_role", "owner"),
					resource.TestCheckResourceAttrSet("influxdb_user.test", "org_id"),
				),
			},
			// Switch back to first organization as owner
			{
				Config: providerConfig + testAccUserResourceWithTwoOrgsConfig(userName, password, org1Name, org2Name, "org1", "owner"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("influxdb_user.test", "name", userName),
					resource.TestCheckResourceAttr("influxdb_user.test", "org_role", "owner"),
				),
			},
			// Remove from organization completely
			{
				Config: providerConfig + testAccUserResourceConfig(userName, password),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("influxdb_user.test", "name", userName),
					resource.TestCheckNoResourceAttr("influxdb_user.test", "org_id"),
					resource.TestCheckNoResourceAttr("influxdb_user.test", "org_role"),
				),
			},
		},
	})
}

func TestAccUserResourceValidation(t *testing.T) {
	userName := acctest.RandomWithPrefix("tf-user-validation-test")
	password := acctest.RandomWithPrefix("password")
	orgName := acctest.RandomWithPrefix("tf-org-test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Test validation: org_id without org_role should fail
			{
				Config:      providerConfig + testAccUserResourceWithOrgIdOnlyConfig(userName, password, orgName),
				ExpectError: regexp.MustCompile(`Attribute "org_role" must be specified when "org_id" is specified`),
			},
			// Test validation: org_role without org_id should fail
			{
				Config:      providerConfig + testAccUserResourceWithOrgRoleOnlyConfig(userName, password),
				ExpectError: regexp.MustCompile(`Attribute "org_id" must be specified when "org_role" is specified`),
			},
			// Test valid configuration: both org_id and org_role specified
			{
				Config: providerConfig + testAccUserResourceWithOrgConfig(userName, password, orgName, "member"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("influxdb_user.test", "name", userName),
					resource.TestCheckResourceAttr("influxdb_user.test", "org_role", "member"),
					resource.TestCheckResourceAttrSet("influxdb_user.test", "org_id"),
				),
			},
		},
	})
}

func TestAccUserResourceRoleChange(t *testing.T) {
	userName := acctest.RandomWithPrefix("tf-user-role-test")
	password := acctest.RandomWithPrefix("password")
	orgName := acctest.RandomWithPrefix("tf-org-test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create user as member
			{
				Config: providerConfig + testAccUserResourceWithOrgConfig(userName, password, orgName, "member"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("influxdb_user.test", "name", userName),
					resource.TestCheckResourceAttr("influxdb_user.test", "org_role", "member"),
					resource.TestCheckResourceAttrSet("influxdb_user.test", "org_id"),
				),
			},
			// Change role to owner in same organization
			{
				Config: providerConfig + testAccUserResourceWithOrgConfig(userName, password, orgName, "owner"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("influxdb_user.test", "name", userName),
					resource.TestCheckResourceAttr("influxdb_user.test", "org_role", "owner"),
					resource.TestCheckResourceAttrSet("influxdb_user.test", "org_id"),
				),
			},
			// Change role back to member
			{
				Config: providerConfig + testAccUserResourceWithOrgConfig(userName, password, orgName, "member"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("influxdb_user.test", "name", userName),
					resource.TestCheckResourceAttr("influxdb_user.test", "org_role", "member"),
				),
			},
		},
	})
}

func testAccUserResourceConfig(name, password string) string {
	return fmt.Sprintf(`
resource "influxdb_user" "test" {
  name     = %[1]q
  password = %[2]q
}
`, name, password)
}

func testAccUserResourceConfigUpdated(name, password string) string {
	return fmt.Sprintf(`
resource "influxdb_user" "test" {
  name     = %[1]q
  password = %[2]q
  status   = "inactive"
}
`, name, password)
}

func testAccUserResourceWithOrgConfig(userName, password, orgName, role string) string {
	return fmt.Sprintf(`
resource "influxdb_organization" "test" {
  name = %[3]q
}

resource "influxdb_user" "test" {
  name     = %[1]q
  password = %[2]q
  org_id   = influxdb_organization.test.id
  org_role = %[4]q
}
`, userName, password, orgName, role)
}

func testAccUserResourceWithTwoOrgsConfig(userName, password, org1Name, org2Name, selectedOrg, role string) string {
	orgRef := fmt.Sprintf("influxdb_organization.%s.id", selectedOrg)
	return fmt.Sprintf(`
resource "influxdb_organization" "org1" {
  name = %[3]q
}

resource "influxdb_organization" "org2" {
  name = %[4]q
}

resource "influxdb_user" "test" {
  name     = %[1]q
  password = %[2]q
  org_id   = %[6]s
  org_role = %[7]q
}
`, userName, password, org1Name, org2Name, selectedOrg, orgRef, role)
}

func testAccUserResourceWithOrgIdOnlyConfig(userName, password, orgName string) string {
	return fmt.Sprintf(`
resource "influxdb_organization" "test" {
  name = %[3]q
}

resource "influxdb_user" "test" {
  name     = %[1]q
  password = %[2]q
  org_id   = influxdb_organization.test.id
}
`, userName, password, orgName)
}

func testAccUserResourceWithOrgRoleOnlyConfig(userName, password string) string {
	return fmt.Sprintf(`
resource "influxdb_user" "test" {
  name     = %[1]q
  password = %[2]q
  org_role = "member"
}
`, userName, password)
}
