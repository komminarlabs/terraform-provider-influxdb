package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

// Ensure the implementation satisfies the expected interfaces.
var _ provider.Provider = &InfluxDBProvider{}

// InfluxDBProvider defines the provider implementation.
type InfluxDBProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// InfluxDBProviderModel maps provider schema data to a Go type.
type InfluxDBProviderModel struct {
	URL   types.String `tfsdk:"url"`
	Token types.String `tfsdk:"token"`
}

// Metadata returns the provider type name.
func (p *InfluxDBProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "influxdb"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *InfluxDBProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"url": schema.StringAttribute{
				Description: "The InfluxDB Cloud Dedicated server URL",
				Optional:    true,
			},
			"token": schema.StringAttribute{
				Description: "An InfluxDB token string",
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
}

// Configure prepares a InfluxDB API client for data sources and resources.
func (p *InfluxDBProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration
	var config InfluxDBProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.URL.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("url"),
			"Unknown InfluxDB URL",
			"The provider cannot create the InfluxDB client as there is an unknown configuration value for the InfluxDB URL. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the INFLUXDB_URL environment variable.",
		)
	}

	if config.Token.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Unknown InfluxDB Token",
			"The provider cannot create the InfluxDB client as there is an unknown configuration value for the InfluxDB Token. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the INFLUXDB_TOKEN environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	url := os.Getenv("INFLUXDB_URL")
	token := os.Getenv("INFLUXDB_TOKEN")

	if !config.URL.IsNull() {
		url = config.URL.ValueString()
	}

	if !config.Token.IsNull() {
		token = config.Token.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if url == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("url"),
			"Missing InfluxDB URL",
			"The provider cannot create the InfluxDB client as there is a missing or empty value for the InfluxDB URL. "+
				"Set the host value in the configuration or use the INFLUXDB_URL environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if token == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Missing InfluxDB Token",
			"The provider cannot create the InfluxDB client as there is a missing or empty value for the InfluxDB Token. "+
				"Set the host value in the configuration or use the INFLUXDB_TOKEN environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "INFLUXDB_URL", url)
	ctx = tflog.SetField(ctx, "INFLUXDB_TOKEN", token)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "INFLUXDB_TOKEN")

	tflog.Debug(ctx, "Creating InfluxDB client")

	// Create a new InfluxDB client using the configuration values
	client := influxdb2.NewClient(url, token)

	_, err := client.Ping(context.Background())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create InfluxDB Client",
			"An unexpected error occurred when creating the InfluxDB client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"InfluxDB Client Error: "+err.Error(),
		)
		return
	}

	// Make the InfluxDB client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured InfluxDB client", map[string]any{"success": true})
}

// Resources defines the resources implemented in the provider.
func (p *InfluxDBProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewAuthorizationResource,
		NewBucketResource,
		NewOrganizationResource,
	}
}

// DataSources defines the data sources implemented in the provider.
func (p *InfluxDBProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewAuthorizationDataSource,
		NewAuthorizationsDataSource,
		NewBucketDataSource,
		NewBucketsDataSource,
		NewOrganizationDataSource,
		NewOrganizationsDataSource,
	}
}

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &InfluxDBProvider{
			version: version,
		}
	}
}
