package provider

import (
	"context"

	"github.com/liamcervante/terraform-provider-fakelocal/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/liamcervante/terraform-provider-fakelocal/internal/dynamic"
)

var _ tfsdk.ResourceType = dynamicResourceType{}

type dynamicResourceType struct {
	Object dynamic.Attribute
}

func (t dynamicResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	tflog.Trace(ctx, "dynamicResourceType.GetSchema")
	return t.Object.ToTerraformSchema(true)
}

func (t dynamicResourceType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	tflog.Trace(ctx, "dynamicResourceType.NewResource")

	provider, diags := convertProviderType(in)

	return client.Resource{
		Client: provider.client,
	}, diags
}
