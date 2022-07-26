package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"

	"github.com/liamcervante/terraform-provider-fakelocal/internal/client"
	"github.com/liamcervante/terraform-provider-fakelocal/internal/simple"
)

var _ tfsdk.ResourceType = simpleResourceType{}

type simpleResourceType struct{}

func (t simpleResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	tflog.Trace(ctx, "simpleResourceType.GetSchema")

	return tfsdk.Schema{
		Attributes: simple.Attributes(),
	}, nil
}

func (t simpleResourceType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	tflog.Trace(ctx, "simpleResourceType.NewResource")

	provider, diags := convertProviderType(in)

	return client.Resource{
		Client: provider.client,
	}, diags
}
