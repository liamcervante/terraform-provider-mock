package provider

import (
	"context"

	"github.com/liamcervante/terraform-provider-mock/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/liamcervante/terraform-provider-mock/internal/dynamic"
)

var _ tfsdk.ResourceType = dynamicResourceType{}

type dynamicResourceType struct {
	Resource dynamic.Resource
}

func (t dynamicResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	tflog.Trace(ctx, "dynamicResourceType.GetSchema")
	return t.Resource.ToTerraformSchema(true)
}

func (t dynamicResourceType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	tflog.Trace(ctx, "dynamicResourceType.NewResource")

	provider, diags := convertProviderType(in)

	return client.Resource{
		Client:       provider.client,
		UseOnlyState: provider.useOnlyState,
	}, diags
}
