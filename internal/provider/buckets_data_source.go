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
	_ datasource.DataSource              = &BucketsDataSource{}
	_ datasource.DataSourceWithConfigure = &BucketsDataSource{}
)

// NewBucketsDataSource is a helper function to simplify the provider implementation.
func NewBucketsDataSource() datasource.DataSource {
	return &BucketsDataSource{}
}

// BucketsDataSource is the data source implementation.
type BucketsDataSource struct {
	client influxdb2.Client
}

// BucketsDataSourceModel describes the data source data model.
type BucketsDataSourceModel struct {
	Buckets []BucketDataSourceModel `tfsdk:"buckets"`
}

// Metadata returns the data source type name.
func (d *BucketsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_buckets"
}

// Schema defines the schema for the data source.
func (d *BucketsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Attributes: map[string]schema.Attribute{
			"buckets": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "Bucket ID.",
						},
						"org_id": schema.StringAttribute{
							Computed:    true,
							Description: "Organization ID.",
						},
						"type": schema.StringAttribute{
							Computed:    true,
							Description: "Bucket type.",
						},
						"schema_type": schema.StringAttribute{
							Computed:    true,
							Description: "Bucket schema type.",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "Bucket description.",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "Bucket name.",
						},
						"created_at": schema.StringAttribute{
							Computed:    true,
							Description: "Bucket creation date.",
						},
						"updated_at": schema.StringAttribute{
							Computed:    true,
							Description: "Last bucket update date.",
						},
						"retention_days": schema.Int64Attribute{
							Computed:    true,
							Description: "Duration bucket retains data.",
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *BucketsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *BucketsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state BucketsDataSourceModel

	buckets, err := d.client.BucketsAPI().GetBuckets(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to list buckets",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, bucket := range *buckets {
		retentionDays := int64(0)
		if len(bucket.RetentionRules) > 0 {
			retentionDays = int64(bucket.RetentionRules[0].EverySeconds) / 24 / 60 / 60
		}

		bucketState := BucketDataSourceModel{
			ID:            types.StringValue(*bucket.Id),
			OrgID:         types.StringValue(*bucket.OrgID),
			Type:          types.StringValue(string(*bucket.Type)),
			SchemaType:    types.StringValue(string(*bucket.SchemaType)),
			Description:   types.StringPointerValue(bucket.Description),
			Name:          types.StringValue(bucket.Name),
			CreatedAt:     types.StringValue(bucket.CreatedAt.String()),
			UpdatedAt:     types.StringValue(bucket.UpdatedAt.String()),
			RetentionDays: types.Int64Value(retentionDays),
		}

		state.Buckets = append(state.Buckets, bucketState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
