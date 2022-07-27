package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"

	"github.com/liamcervante/terraform-provider-mock/internal/client"
	"github.com/liamcervante/terraform-provider-mock/internal/complex"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ tfsdk.DataSourceType = complexDataSourceType{}

type complexDataSourceType struct{}

func (t complexDataSourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	tflog.Trace(ctx, "complexDataSourceType.GetSchema")

	return tfsdk.Schema{
		Attributes: complex.Attributes(5),
		Blocks:     complex.Blocks(3),
	}, nil
}

func (t complexDataSourceType) NewDataSource(ctx context.Context, in tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	tflog.Trace(ctx, "complexDataSourceType.NewDataSource")

	provider, diags := convertProviderType(in)

	return client.DataSource{
		Client: provider.client,
	}, diags
}
