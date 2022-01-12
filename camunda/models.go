package camunda

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Deployment -
type Deployment struct {
	Id        types.String `tfsdk:"id"`
	Key       types.String `tfsdk:"key"`
	Tenant    types.String `tfsdk:"tenant"`
	Resources []Resources  `tfsdk:"resources"`
}

type Resources struct {
	Name    string `tfsdk:"name"`
	Content string `tfsdk:"content"`
}
