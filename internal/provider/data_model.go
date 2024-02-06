package provider

import "github.com/hashicorp/terraform-plugin-framework/types"

// BucketDataSourceModel maps Buckets schema data.
type BucketDataSourceModel struct {
	ID            types.String `tfsdk:"id"`
	OrgID         types.String `tfsdk:"org_id"`
	Type          types.String `tfsdk:"type"`
	SchemaType    types.String `tfsdk:"schema_type"`
	Description   types.String `tfsdk:"description"`
	Name          types.String `tfsdk:"name"`
	CreatedAt     types.String `tfsdk:"created_at"`
	UpdatedAt     types.String `tfsdk:"updated_at"`
	RetentionDays types.Int64  `tfsdk:"retention_days"`
}
