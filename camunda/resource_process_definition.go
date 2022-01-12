package camunda

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"terraform-provider-camunda/camunda/client"
)

type resourceDeploymentType struct{}

// Deployment Resource schema
func (r resourceDeploymentType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Type:        types.StringType,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: "The id to identify the deployment by",
			},
			"key": {
				Type:        types.StringType,
				Required:    true,
				Description: "The key to identify the deployment by",
				PlanModifiers: []tfsdk.AttributePlanModifier{
					tfsdk.RequiresReplace(),
				},
			},
			"tenant": {
				Type:        types.StringType,
				Optional:    true,
				Description: "The tenant that the deployment belongs to",
				PlanModifiers: []tfsdk.AttributePlanModifier{
					tfsdk.RequiresReplace(),
				},
			},
			"resources": {
				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
					"name": {
						Type:     types.StringType,
						Required: true,
					},
					"content": {
						Type:     types.StringType,
						Required: true,
					},
				}, tfsdk.ListNestedAttributesOptions{}),
				Required:    true,
				Description: "The resources to deploy",
			},
		},
	}, nil
}

// New resource instance
func (r resourceDeploymentType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return resourceDeployment{
		p: *(p.(*provider)),
	}, nil
}

type resourceDeployment struct {
	p provider
}

// Create a new resource
func (r resourceDeployment) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, likely because it depends on an unknown value from another resource. This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
		)
		return
	}

	var plan Deployment
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resources := map[string]interface{}{}
	for _, item := range plan.Resources {
		resources[item.Name] = client.Named{
			Name:    item.Name,
			Content: item.Content,
		}
	}

	source := "Terraform Provider Camunda"
	changedOnly := true

	created, err := r.p.client.Deployment.Create(client.ReqDeploymentCreate{
		DeploymentName:    plan.Key.Value,
		DeploymentSource:  &source,
		DeployChangedOnly: &changedOnly,
		Resources:         resources,
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating deployment",
			"Could not create key "+plan.Key.Value+": "+err.Error(),
		)
		return
	}

	var result = Deployment{
		Id:        types.String{Value: created.Id},
		Key:       plan.Key,
		Tenant:    plan.Tenant,
		Resources: plan.Resources,
	}

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r resourceDeployment) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, likely because it depends on an unknown value from another resource. This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
		)
		return
	}

	// Get current state
	var state Deployment
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get order from API and then update what is in state from what the API returns
	id := state.Id.Value

	// Get order current value
	deployment, err := r.p.client.Deployment.Get(id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading deployment",
			"Could not read id "+id+": "+err.Error(),
		)
		return
	}

	state.Id = types.String{Value: deployment.Id}
	state.Key = types.String{Value: deployment.Name}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update resource
func (r resourceDeployment) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, likely because it depends on an unknown value from another resource. This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
		)
		return
	}

	var plan Deployment
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resources := map[string]interface{}{}
	for _, item := range plan.Resources {
		resources[item.Name] = client.Named{
			Name:    item.Name,
			Content: item.Content,
		}
	}

	source := "terraform"
	changedOnly := true

	created, err := r.p.client.Deployment.Create(client.ReqDeploymentCreate{
		DeploymentName:    plan.Key.Value,
		DeploymentSource:  &source,
		DeployChangedOnly: &changedOnly,
		Resources:         resources,
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating deployment",
			"Could not create key "+plan.Key.Value+": "+err.Error(),
		)
		return
	}

	var result = Deployment{
		Id:        types.String{Value: created.Id},
		Key:       plan.Key,
		Tenant:    plan.Tenant,
		Resources: plan.Resources,
	}

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete resource
func (r resourceDeployment) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, likely because it depends on an unknown value from another resource. This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
		)
		return
	}

	// Get current state
	var state Deployment
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.Value

	err := r.p.client.Deployment.Delete(
		id,
		map[string]string{},
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting deployment",
			"Could not delete id "+id+": "+err.Error(),
		)
		return
	}

	// Remove resource from state
	resp.State.RemoveResource(ctx)
}

// Import resource
func (r resourceDeployment) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	// Save the import identifier in the key attribute
	tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("key"), req, resp)
}
