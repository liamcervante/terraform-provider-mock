package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/liamcervante/terraform-provider-mock/internal/client"
	"github.com/liamcervante/terraform-provider-mock/internal/complex"
)

var _ tfsdk.ResourceType = complexResourceType{}

type complexResourceType struct{}

func (t complexResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	tflog.Trace(ctx, "complexResourceType.GetSchema")

	return tfsdk.Schema{
		Attributes: complex.Attributes(3),
		Blocks:     complex.Blocks(3),
	}, nil
}

func (t complexResourceType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	tflog.Trace(ctx, "complexResourceType.NewResource")

	provider, diags := convertProviderType(in)

	return client.Resource{
		Client:       provider.client,
		UseOnlyState: provider.useOnlyState,
	}, diags
}
