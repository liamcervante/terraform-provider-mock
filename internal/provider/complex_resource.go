package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/liamcervante/terraform-provider-fakelocal/internal/resource"
)

var _ tfsdk.ResourceType = complexResourceType{}

type complexResourceType struct{}

func (t complexResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	tflog.Trace(ctx, "complexResourceType.GetSchema")
	return tfsdk.Schema{
		Attributes: resource.ComplexAttributes,
	}, nil
}

func (t complexResourceType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return resource.Resource{
		Client:         provider.client,
		CreateResource: resource.NewComplex,
	}, diags
}