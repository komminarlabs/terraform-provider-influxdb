package provider

import "github.com/hashicorp/terraform-plugin-framework/types"

// AuthorizationModel maps InfluxDB authorization schema data.
type AuthorizationModel struct {
	Id          types.String                   `tfsdk:"id"`
	Token       types.String                   `tfsdk:"token"`
	Status      types.String                   `tfsdk:"status"`
	Description types.String                   `tfsdk:"description"`
	OrgID       types.String                   `tfsdk:"org_id"`
	Org         types.String                   `tfsdk:"org"`
	UserId      types.String                   `tfsdk:"user_id"`
	User        types.String                   `tfsdk:"user"`
	CreatedAt   types.String                   `tfsdk:"created_at"`
	UpdatedAt   types.String                   `tfsdk:"updated_at"`
	Permissions []AuthorizationPermissionModel `tfsdk:"permissions"`
}

// AuthorizationPermissionModel maps InfluxDB authorization permission schema data.
type AuthorizationPermissionModel struct {
	Action   types.String                          `tfsdk:"action"`
	Resource AuthorizationPermissionrResourceModel `tfsdk:"resource"`
}

// AuthorizationPermissionrResourceModel maps InfluxDB authorization permission resource schema data.
type AuthorizationPermissionrResourceModel struct {
	Id    types.String `tfsdk:"id"`
	Type  types.String `tfsdk:"type"`
	OrgID types.String `tfsdk:"org_id"`
	Org   types.String `tfsdk:"org"`
}

// BucketModel maps InfluxDB bucket schema data.
type BucketModel struct {
	Id            types.String `tfsdk:"id"`
	OrgID         types.String `tfsdk:"org_id"`
	Type          types.String `tfsdk:"type"`
	Description   types.String `tfsdk:"description"`
	Name          types.String `tfsdk:"name"`
	CreatedAt     types.String `tfsdk:"created_at"`
	UpdatedAt     types.String `tfsdk:"updated_at"`
	RetentionDays types.Int64  `tfsdk:"retention_days"`
}

// OrganizationModel maps InfluxDB organization schema data.
type OrganizationModel struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	CreatedAt   types.String `tfsdk:"created_at"`
	UpdatedAt   types.String `tfsdk:"updated_at"`
}
