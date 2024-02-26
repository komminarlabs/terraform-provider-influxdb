package provider

import "github.com/hashicorp/terraform-plugin-framework/types"

// BucketModel maps InfluxDB bucket schema data.
type BucketModel struct {
	Id              types.String `tfsdk:"id"`
	OrgID           types.String `tfsdk:"org_id"`
	Type            types.String `tfsdk:"type"`
	Description     types.String `tfsdk:"description"`
	Name            types.String `tfsdk:"name"`
	CreatedAt       types.String `tfsdk:"created_at"`
	UpdatedAt       types.String `tfsdk:"updated_at"`
	RetentionPeriod types.Int64  `tfsdk:"retention_period"`
}
