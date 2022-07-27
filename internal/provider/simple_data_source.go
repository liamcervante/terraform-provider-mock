package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"

	"github.com/liamcervante/terraform-provider-mock/internal/client"
	"github.com/liamcervante/terraform-provider-mock/internal/simple"
)

var _ tfsdk.DataSourceType = simpleDataSourceType{}

type simpleDataSourceType struct{}

func (t simpleDataSourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	tflog.Trace(ctx, "simpleDataSourceType.GetSchema")

	return tfsdk.Schema{
		Attributes: simple.Attributes(),
	}, nil
}

func (t simpleDataSourceType) NewDataSource(ctx context.Context, in tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	tflog.Trace(ctx, "simpleDataSourceType.NewDataSource")

	provider, diags := convertProviderType(in)

	return client.DataSource{
		Client: provider.client,
	}, diags
}
