package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

const (
	// providerConfig is a shared configuration to combine with the actual
	// test configuration so the HashiCups client is properly configured.
	providerConfig = `
  provider "influxdb" {}
  `
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"influxdb": providerserver.NewProtocol6WithError(New("test")()),
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("INFLUXDB_URL"); v == "" {
		t.Fatal("INFLUXDB_URL must be set for acceptance tests")
	}

	// Check for authentication credentials - either token OR username+password
	token := os.Getenv("INFLUXDB_TOKEN")
	username := os.Getenv("INFLUXDB_USERNAME")
	password := os.Getenv("INFLUXDB_PASSWORD")

	hasToken := token != ""
	hasUsernamePassword := username != "" && password != ""

	if !hasToken && !hasUsernamePassword {
		t.Fatal("Authentication credentials must be set for acceptance tests. " +
			"Provide either INFLUXDB_TOKEN or both INFLUXDB_USERNAME and INFLUXDB_PASSWORD environment variables")
	}

	if v := os.Getenv("INFLUXDB_ORG_ID"); v == "" {
		t.Fatal("INFLUXDB_ORG_ID must be set for acceptance tests")
	}
}
