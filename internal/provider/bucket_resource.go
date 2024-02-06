package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/domain"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &BucketResource{}
var _ resource.ResourceWithImportState = &BucketResource{}

// NewBucketResource is a helper function to simplify the provider implementation.
func NewBucketResource() resource.Resource {
	return &BucketResource{}
}

// BucketResource defines the resource implementation.
type BucketResource struct {
	client influxdb2.Client
}

// Metadata returns the resource type name.
func (r *BucketResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_bucket"
}

// Schema defines the schema for the resource.
func (r *BucketResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Manages an InfluxDB bucket",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Bucket ID.",
			},
			"org_id": schema.StringAttribute{
				Required:    true,
				Description: "Organization ID.",
			},
			"type": schema.StringAttribute{
				Computed:    true,
				Description: "Bucket type.",
			},
			"schema_type": schema.StringAttribute{
				Computed:    true,
				Description: "Bucket schema type.",
				//Validators: ,
			},
			"description": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "Bucket description.",
			},
			"name": schema.StringAttribute{
				Required:    true,
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
				Optional:    true,
				Description: "Duration bucket retains data.",
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *BucketResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan BucketDataSourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	bucket := domain.Bucket{
		OrgID:          plan.OrgID.ValueStringPointer(),
		Name:           plan.Name.ValueString(),
		Description:    plan.Description.ValueStringPointer(),
		RetentionRules: []domain.RetentionRule{},
	}

	retention_days := plan.RetentionDays
	if !retention_days.IsNull() {
		bucket.RetentionRules = append(bucket.RetentionRules, domain.RetentionRule{
			EverySeconds: int64(retention_days.ValueInt64() * 24 * 60 * 60),
		})
	}

	apiResponse, err := r.client.BucketsAPI().CreateBucket(ctx, &bucket)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating bucket",
			"Could not create bucket, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.StringValue(*apiResponse.Id)
	plan.OrgID = types.StringValue(*apiResponse.OrgID)
	plan.Type = types.StringValue(string(*apiResponse.Type))
	plan.SchemaType = types.StringValue(string(*apiResponse.SchemaType))
	plan.Description = types.StringPointerValue(apiResponse.Description)
	plan.CreatedAt = types.StringValue(apiResponse.CreatedAt.String())
	plan.UpdatedAt = types.StringValue(apiResponse.UpdatedAt.String())
	plan.RetentionDays = types.Int64Value(retention_days.ValueInt64())

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *BucketResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state BucketDataSourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed bucket value from InfluxDB
	bucket, err := r.client.BucketsAPI().FindBucketByName(ctx, state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Bucket not found",
			err.Error(),
		)
		return
	}

	// Overwrite items with refreshed state
	state.ID = types.StringValue(*bucket.Id)
	state.OrgID = types.StringValue(*bucket.OrgID)
	state.Type = types.StringValue(string(*bucket.Type))
	state.SchemaType = types.StringValue(string(*bucket.SchemaType))
	state.Description = types.StringPointerValue(bucket.Description)
	state.CreatedAt = types.StringValue(bucket.CreatedAt.String())
	state.UpdatedAt = types.StringValue(bucket.UpdatedAt.String())

	retentionDays := int64(0)
	if len(bucket.RetentionRules) > 0 {
		retentionDays = int64(bucket.RetentionRules[0].EverySeconds) / 24 / 60 / 60
	}
	state.RetentionDays = types.Int64Value(retentionDays)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *BucketResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan BucketDataSourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	bucket := domain.Bucket{
		OrgID:          plan.OrgID.ValueStringPointer(),
		Id:             plan.ID.ValueStringPointer(),
		Name:           plan.Name.ValueString(),
		Description:    plan.Description.ValueStringPointer(),
		RetentionRules: []domain.RetentionRule{},
	}

	retention_days := plan.RetentionDays
	if !retention_days.IsNull() {
		bucket.RetentionRules = append(bucket.RetentionRules, domain.RetentionRule{
			EverySeconds: int64(retention_days.ValueInt64() * 24 * 60 * 60),
		})
	}

	// Update existing bucket
	apiResponse, err := r.client.BucketsAPI().UpdateBucket(ctx, &bucket)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating bucket",
			"Could not update bucket, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.StringValue(*apiResponse.Id)
	plan.OrgID = types.StringValue(*apiResponse.OrgID)
	plan.Type = types.StringValue(string(*apiResponse.Type))
	plan.SchemaType = types.StringValue(string(*apiResponse.SchemaType))
	plan.Description = types.StringPointerValue(apiResponse.Description)
	plan.CreatedAt = types.StringValue(apiResponse.CreatedAt.String())
	plan.UpdatedAt = types.StringValue(apiResponse.UpdatedAt.String())
	plan.RetentionDays = types.Int64Value(retention_days.ValueInt64())

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *BucketResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data BucketDataSourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete example, got error: %s", err))
	//     return
	// }
}

// Configure adds the provider configured client to the resource.
func (r *BucketResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = client
}

func (r *BucketResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
