package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/liamcervante/terraform-provider-fakelocal/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/liamcervante/terraform-provider-fakelocal/internal/dynamic"
)

var _ tfsdk.ResourceType = dynamicResourceType{}

type dynamicResourceType struct {
	Object dynamic.Object
}

func (t dynamicResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	tflog.Trace(ctx, "dynamicResourceType.GetSchema")

	var diags diag.Diagnostics

	if _, ok := t.Object.Attributes["id"]; ok {
		diags.AddError(
			"Found `id` value in top level object",
			"Top level dynamic objects cannot define a value called `id` as the provider will generate an ID for them.",
		)

		return tfsdk.Schema{}, diags
	}

	attributes, err := t.Object.ToTerraformAttribute()
	if err != nil {
		diags.AddError("Failed to parse dynamic attributes", err.Error())
	}

	attributes["id"] = tfsdk.Attribute{
		Computed: true,
		PlanModifiers: tfsdk.AttributePlanModifiers{
			tfsdk.UseStateForUnknown(),
		},
		Type: types.StringType,
	}

	return tfsdk.Schema{
		Attributes: attributes,
	}, diags
}

func (t dynamicResourceType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	tflog.Trace(ctx, "dynamicResourceType.GetSchema")

	provider, diags := convertProviderType(in)

	return client.Resource{
		Client:         provider.client,
		CreateResource: client.NewDynamic,
	}, diags
}
