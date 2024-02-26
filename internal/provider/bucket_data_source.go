package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &BucketDataSource{}
	_ datasource.DataSourceWithConfigure = &BucketDataSource{}
)

// NewBucketDataSource is a helper function to simplify the provider implementation.
func NewBucketDataSource() datasource.DataSource {
	return &BucketDataSource{}
}

// BucketsDataSource is the data source implementation.
type BucketDataSource struct {
	client influxdb2.Client
}

// Metadata returns the data source type name.
func (d *BucketDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_bucket"
}

// Schema defines the schema for the data source.
func (d *BucketDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Retrieves a bucket. Use this data source to retrieve information for a specific bucket.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "A Bucket ID.",
			},
			"org_id": schema.StringAttribute{
				Computed:    true,
				Description: "An organization ID.",
			},
			"type": schema.StringAttribute{
				Computed:    true,
				Description: "The Bucket type.",
			},
			"description": schema.StringAttribute{
				Computed:    true,
				Description: "A description of the bucket.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "A Bucket name.",
			},
			"created_at": schema.StringAttribute{
				Computed:    true,
				Description: "Bucket creation date.",
			},
			"updated_at": schema.StringAttribute{
				Computed:    true,
				Description: "Last bucket update date.",
			},
			"retention_period": schema.Int64Attribute{
				Computed:    true,
				Description: "The duration in seconds for how long data will be kept in the database. The default duration is 2592000 (30 days). 0 represents infinite retention.",
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *BucketDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(influxdb2.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected influxdb2.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

// Read refreshes the Terraform state with the latest data.
func (d *BucketDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state BucketModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	bucketName := state.Name
	if bucketName.IsNull() {
		resp.Diagnostics.AddError(
			"Name is empty",
			"Must set name",
		)

		return
	}

	bucket, err := d.client.BucketsAPI().FindBucketByName(ctx, bucketName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Bucket not found",
			err.Error(),
		)

		return
	}

	// Map response body to model
	state = BucketModel{
		Id:              types.StringPointerValue(bucket.Id),
		OrgID:           types.StringPointerValue(bucket.OrgID),
		Type:            types.StringValue(string(*bucket.Type)),
		Description:     types.StringPointerValue(bucket.Description),
		Name:            types.StringValue(bucket.Name),
		CreatedAt:       types.StringValue(bucket.CreatedAt.String()),
		UpdatedAt:       types.StringValue(bucket.UpdatedAt.String()),
		RetentionPeriod: types.Int64Value(bucket.RetentionRules[0].EverySeconds),
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
