package provider

import "github.com/hashicorp/terraform-plugin-framework/types"

// OrganizationModel maps InfluxDB organization schema data.
type OrganizationModel struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	CreatedAt   types.String `tfsdk:"created_at"`
	UpdatedAt   types.String `tfsdk:"updated_at"`
}
