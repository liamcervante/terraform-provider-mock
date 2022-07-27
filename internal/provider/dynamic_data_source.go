package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/liamcervante/terraform-provider-fakelocal/internal/client"
	"github.com/liamcervante/terraform-provider-fakelocal/internal/dynamic"
)

type dynamicDataSourceType struct {
	Resource dynamic.Resource
}

func (d dynamicDataSourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	tflog.Trace(ctx, "dynamicDataSourceType.GetSchema")

	return d.Resource.ToTerraformSchema(false)
}

func (d dynamicDataSourceType) NewDataSource(ctx context.Context, in tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	tflog.Trace(ctx, "dynamicDataSourceType.NewDataSource")

	provider, diags := convertProviderType(in)

	return client.DataSource{
		Client: provider.client,
	}, diags
}
