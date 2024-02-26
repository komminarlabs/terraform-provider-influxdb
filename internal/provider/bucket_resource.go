package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/domain"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &BucketResource{}
	_ resource.ResourceWithImportState = &BucketResource{}
	_ resource.ResourceWithImportState = &BucketResource{}
)

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
		MarkdownDescription: "Creates and manages a bucket.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "A Bucket ID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"org_id": schema.StringAttribute{
				Required:    true,
				Description: "An organization ID.",
			},
			"type": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "The Bucket type.",
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"user", "system"}...),
				},
			},
			"description": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
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
				Optional:    true,
				Default:     int64default.StaticInt64(2592000),
				Description: "The duration in seconds for how long data will be kept in the database. The default duration is 2592000 (30 days). 0 represents infinite retention.",
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *BucketResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan BucketModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	createBucket := domain.Bucket{
		OrgID:       plan.OrgID.ValueStringPointer(),
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueStringPointer(),
		RetentionRules: []domain.RetentionRule{{
			EverySeconds: plan.RetentionPeriod.ValueInt64(),
		}},
	}

	apiResponse, err := r.client.BucketsAPI().CreateBucket(ctx, &createBucket)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating bucket",
			"Could not create bucket, unexpected error: "+err.Error(),
		)

		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.Id = types.StringPointerValue(apiResponse.Id)
	plan.OrgID = types.StringPointerValue(apiResponse.OrgID)
	plan.Name = types.StringValue(apiResponse.Name)
	plan.Type = types.StringValue(string(*apiResponse.Type))
	plan.Description = types.StringPointerValue(apiResponse.Description)
	plan.CreatedAt = types.StringValue(apiResponse.CreatedAt.String())
	plan.UpdatedAt = types.StringValue(apiResponse.UpdatedAt.String())
	plan.RetentionPeriod = types.Int64Value(apiResponse.RetentionRules[0].EverySeconds)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *BucketResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state BucketModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed bucket value from InfluxDB
	readBucket, err := r.client.BucketsAPI().FindBucketByID(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Bucket not found",
			err.Error(),
		)

		return
	}

	// Overwrite items with refreshed state
	state.Id = types.StringPointerValue(readBucket.Id)
	state.OrgID = types.StringPointerValue(readBucket.OrgID)
	state.Name = types.StringValue(readBucket.Name)
	state.Type = types.StringValue(string(*readBucket.Type))
	state.Description = types.StringPointerValue(readBucket.Description)
	state.CreatedAt = types.StringValue(readBucket.CreatedAt.String())
	state.UpdatedAt = types.StringValue(readBucket.UpdatedAt.String())
	state.RetentionPeriod = types.Int64Value(readBucket.RetentionRules[0].EverySeconds)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *BucketResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan BucketModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	updateBucket := domain.Bucket{
		OrgID:       plan.OrgID.ValueStringPointer(),
		Id:          plan.Id.ValueStringPointer(),
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueStringPointer(),
		RetentionRules: []domain.RetentionRule{{
			EverySeconds: plan.RetentionPeriod.ValueInt64(),
		}},
	}

	// Update existing bucket
	apiResponse, err := r.client.BucketsAPI().UpdateBucket(ctx, &updateBucket)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating bucket",
			"Could not update bucket, unexpected error: "+err.Error(),
		)

		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.Id = types.StringPointerValue(apiResponse.Id)
	plan.OrgID = types.StringPointerValue(apiResponse.OrgID)
	plan.Name = types.StringValue(apiResponse.Name)
	plan.Type = types.StringValue(string(*apiResponse.Type))
	plan.Description = types.StringPointerValue(apiResponse.Description)
	plan.CreatedAt = types.StringValue(apiResponse.CreatedAt.String())
	plan.UpdatedAt = types.StringValue(apiResponse.UpdatedAt.String())
	plan.RetentionPeriod = types.Int64Value(apiResponse.RetentionRules[0].EverySeconds)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *BucketResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state BucketModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing bucket
	err := r.client.BucketsAPI().DeleteBucketWithID(ctx, *state.Id.ValueStringPointer())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting bucket",
			"Could not delete bucket, unexpected error: "+err.Error(),
		)

		return
	}
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
