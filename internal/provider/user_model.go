package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// UserModel maps InfluxDB User schema data.
type UserModel struct {
	Id       types.String `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	Password types.String `tfsdk:"password"`
	OrgId    types.String `tfsdk:"org_id"`
	OrgRole  types.String `tfsdk:"org_role"`
	Status   types.String `tfsdk:"status"`
}
